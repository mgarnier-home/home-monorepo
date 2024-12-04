package server

import (
	"context"
	"fmt"
	"mgarnier11/go-proxy/host"
	"mgarnier11/go-proxy/hostManager"
	"mgarnier11/go-proxy/hostState"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"

	"mgarnier11/go/logger"
)

type contextKey string

const hostContextKey contextKey = "host"

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

func (s *Server) getHostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hostName := vars["host"]
		host := hostManager.GetHost(hostName)

		if host == nil {
			log.Errorf("Host %s not found", hostName)
			http.Error(w, "Host not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), hostContextKey, host)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) Start() error {
	log = logger.NewLogger("[SERVER]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

	router := mux.NewRouter()

	router.Use(s.logRequestMiddleware)

	controlRouter := router.PathPrefix("/control/{host}").Subrouter()
	controlRouter.Use(s.getHostMiddleware)

	controlRouter.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		host.StartHost()

		if host.State == hostState.Started {
			w.Write([]byte(fmt.Sprintf("Host %s has successfully started", host.Config.Name)))
		} else {
			w.Write([]byte(fmt.Sprintf("Host %s failed to start, check logs", host.Config.Name)))
		}
	})

	controlRouter.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		host.StopHost()

		if host.State == hostState.Stopped {
			w.Write([]byte(fmt.Sprintf("Host %s has successfully stopped", host.Config.Name)))
		} else {
			w.Write([]byte(fmt.Sprintf("Host %s failed to stop, check logs", host.Config.Name)))
		}
	})

	controlRouter.HandleFunc("/start-stop", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		if host.State == hostState.Started {
			host.StopHost()
		} else if host.State == hostState.Stopped {
			host.StartHost()
		}

		if host.State == hostState.Started {
			w.Write([]byte(fmt.Sprintf("Host %s has successfully started", host.Config.Name)))
		} else if host.State == hostState.Stopped {
			w.Write([]byte(fmt.Sprintf("Host %s has successfully stopped", host.Config.Name)))
		} else {
			w.Write([]byte(fmt.Sprintf("Host %s failed to start/stop, check logs", host.Config.Name)))
		}

	})

	controlRouter.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		w.WriteHeader(210 + int(host.State))
		w.Write([]byte(fmt.Sprintf("Host %s is %s", host.Config.Name, host.State.String())))
	})

	controlRouter.HandleFunc("/autostop-toggle", func(w http.ResponseWriter, r *http.Request) {
		host := r.Context().Value(hostContextKey).(*host.Host)

		host.Config.Autostop = !host.Config.Autostop
		host.Config.Save()

		if host.Config.Autostop {
			w.Write([]byte(fmt.Sprintf("Host %s autostop enabled", host.Config.Name)))
		} else {
			w.Write([]byte(fmt.Sprintf("Host %s autostop disabled", host.Config.Name)))
		}
	})

	log.Infof("Starting server on port %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)

}
