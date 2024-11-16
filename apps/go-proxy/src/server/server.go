package server

import (
	"fmt"
	"mgarnier11/go-proxy/host"
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/gorilla/mux"
)

type Server struct {
	port int

	hosts *(map[string]*host.Host)
}

func NewServer(port int, hosts *(map[string]*host.Host)) *Server {
	return &Server{
		port:  port,
		hosts: hosts,
	}
}

func (s *Server) getHostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	controlRouter := router.PathPrefix("/control").Subrouter()
	controlRouter.Use(s.getHostMiddleware)

	log.Infof("Starting server on port %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)

}
