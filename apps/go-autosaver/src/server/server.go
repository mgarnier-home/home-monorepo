package server

import (
	"fmt"
	"net/http"

	"mgarnier11.fr/go/libs/httputils"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/version"

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

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.Use(httputils.LogRequestMiddleware(s.logger))
	router.Use(httputils.CorsMiddleware)

	version.SetupVersionRoute(router)

	// routers.NewHostsRouter(router, s.logger)
	// routers.NewMapsRouter(router, s.logger)
	// routers.NewServersRouter(router, s.logger)

	s.logger.Infof("Starting server on port %d", s.port)

	fs := http.FileServer(http.Dir("frontend"))

	router.PathPrefix("/").Handler(fs)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)

}
