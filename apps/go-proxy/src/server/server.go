package server

import (
	"context"
	"fmt"
	"mgarnier11/go-proxy/host"
	"mgarnier11/go-proxy/hostmanager"
	"mgarnier11/go-proxy/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type contextKey string

const hostContextKey contextKey = "host"

type Server struct {
	port int
}

var logger = utils.NewLogger("API", "[%s] ", nil)

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("Request: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) getHostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hostName := vars["host"]
		host := hostmanager.GetHost(hostName)

		if host == nil {
			logger.Errorf("Host %s not found", hostName)
			http.Error(w, "Host not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), hostContextKey, host)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.Use(s.logRequestMiddleware)

	controlRouter := router.PathPrefix("/control/{host}").Subrouter()
	controlRouter.Use(s.getHostMiddleware)

	controlRouter.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		w.Write([]byte(fmt.Sprintf("Starting host %s", host.Config.Name)))

		host.StartHost(nil)
	})

	controlRouter.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		w.Write([]byte(fmt.Sprintf("Stopping host %s", host.Config.Name)))

		host.StopHost()
	})

	controlRouter.HandleFunc("/start-stop", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		started, err := host.HostStarted()

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Failed to check host status: %v", err)))
			return
		}

		if started {
			host.StopHost()

			w.Write([]byte(fmt.Sprintf("Stopping host %s", host.Config.Name)))
		} else {
			host.StartHost(nil)

			w.Write([]byte(fmt.Sprintf("Starting host %s", host.Config.Name)))
		}
	})

	controlRouter.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		started, err := host.HostStarted()

		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Failed to check host status: %v", err)))
			return
		}

		if started {
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("Host %s is started", host.Config.Name)))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Host %s is stopped", host.Config.Name)))
		}
	})

	logger.Infof("Starting server on port %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)

}
