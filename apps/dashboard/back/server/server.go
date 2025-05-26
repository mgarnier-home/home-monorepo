package server

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"mgarnier11.fr/go/dashboard/config"
	"mgarnier11.fr/go/libs/httputils"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/version"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for demo purposes
		return true
	},
}

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

	router.HandleFunc("/ws", handleWebSocket).Methods("GET")

	s.logger.Infof("Starting server on port %d", s.port)

	fs := http.FileServer(http.Dir(config.Config.AppDistPath))

	router.PathPrefix("/").Handler(fs)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	logger.Infof("Client connected")

	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Read error: %v", err)
			break
		}
		logger.Debugf("Received: %s", message)

		// Echo message back to client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			logger.Errorf("Write error: %v", err)
			break
		}
	}
}
