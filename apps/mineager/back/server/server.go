package server

import (
	"fmt"
	"mgarnier11/go/logger"
	"mgarnier11/go/version"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/server/routes"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
)

type Server struct {
	port int
}

var log *logger.Logger

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Request: %s %s", r.Method, r.URL.Path)

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
	log = logger.NewLogger("[SERVER]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

	router := mux.NewRouter()

	router.Use(s.logRequestMiddleware)
	if config.Config.ApiToken != "" {
		router.Use(s.checkApiTokenMiddleware)
	}
	router.Use(s.corsMiddleware)

	version.SetupVersionRoute(router)

	routes.MapRoutes(router)
	routes.ServerRoutes(router)

	log.Infof("Starting server on port %d", s.port)

	fs := http.FileServer(http.Dir(config.Config.FrontendPath))

	router.PathPrefix("/").Handler(fs)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)

}
