package server

import (
	"fmt"
	"mgarnier11/go/logger"
	"mgarnier11/go/version"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/server/routers"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
)

type Server struct {
	port   int
	logger *logger.Logger
}

func NewServer(port int) *Server {
	return &Server{
		port:   port,
		logger: logger.NewLogger("[SERVER]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
	}
}

func (s *Server) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("Request: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) checkApiTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiToken := r.Header.Get("Api-Token")

		if apiToken != config.Config.ApiToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace * with specific origin(s) if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Api-Token, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.Use(s.logRequestMiddleware)
	if config.Config.ApiToken != "" {
		router.Use(s.checkApiTokenMiddleware)
	}
	router.Use(s.corsMiddleware)

	version.SetupVersionRoute(router)

	hostsRouter := routers.NewHostsRouter(router, s.logger)

	routes.MapsRoutes(router)
	hostNameRouter := routes.HostsRoutes(router)
	routes.ServersRoutes(hostNameRouter)

	s.logger.Infof("Starting server on port %d", s.port)

	fs := http.FileServer(http.Dir(config.Config.FrontendPath))

	router.PathPrefix("/").Handler(fs)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)

}
