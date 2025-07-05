package server

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
	"github.com/zishang520/socket.io/v2/socket"
	"mgarnier11.fr/go/dashboard/config"
	"mgarnier11.fr/go/dashboard/server/state"
	"mgarnier11.fr/go/libs/httputils"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/version"
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

	s.logger.Infof("Starting server on port %d", s.port)

	fs := http.FileServer(http.Dir(config.Config.AppDistPath))
	io := socket.NewServer(nil, nil)

	router.Handle("/socket.io/", io.ServeHandler(nil))
	router.PathPrefix("/").Handler(fs)

	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		// defer client.Conn().Close(true)
		logger.Infof("Client connected: %s", client.Id())
		// client.On("event", func(datas ...any) {
		// 	logger.Infof("Received event from client %s: %v", client.Id(), datas)
		// 	// Here you can handle the event and send a response back to the client if needed
		// 	client.Emit("response", "Event received successfully")
		// })
		client.On("disconnect", func(...any) {
			logger.Infof("Client disconnected: %s", client.Id())

			state.ClientDisconnect()
		})

		dashboardConfig, err := config.Config.GetDashboardConfig()

		if err != nil {
			logger.Errorf("Error loading dashboardConfig: %v", err)
			return
		}

		client.Emit("dashboardConfig", dashboardConfig)

		state.ClientConnect(s.logger, io, client, dashboardConfig)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)
}
