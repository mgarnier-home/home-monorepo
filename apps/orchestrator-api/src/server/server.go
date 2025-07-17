package server

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
	"mgarnier11.fr/go/libs/httputils"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/version"
	"mgarnier11.fr/go/orchestrator-api/compose"
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

	router.HandleFunc("/compose", s.getComposeFiles).Methods("GET")
	router.HandleFunc("/commands", s.getCommands).Methods("GET")
	router.HandleFunc("/exec-command/{command}", s.streamExecCommand).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)
}

func (s *Server) getComposeFiles(w http.ResponseWriter, r *http.Request) {
	composeFiles, err := compose.GetComposeFiles()

	if err != nil {
		s.logger.Errorf("Error getting compose files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.logger.Infof("Found %d compose files", len(composeFiles))

	httputils.WriteYamlResponse(w, composeFiles)

	s.logger.Infof("Successfully served %d compose files", len(composeFiles))

}

func (s *Server) streamExecCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	command := vars["command"]

	if command == "" {
		s.logger.Errorf("No command provided in request")
		http.Error(w, "Bad Request: No command provided", http.StatusBadRequest)
		return
	}

	s.logger.Infof("Received command: %s", command)

	commandsToExecute, err := compose.GetCommandsToExecute(command)
	if err != nil {
		s.logger.Errorf("Error getting commands to execute: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	logger.Debugf("Found %d commands to execute", len(commandsToExecute))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	compose.ExecCommandsStream(commandsToExecute, w)
}

func (s *Server) getCommands(w http.ResponseWriter, r *http.Request) {
	composeFiles, err := compose.GetComposeFiles()

	if err != nil {
		s.logger.Errorf("Error getting compose files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commands, err := compose.GetCommands(composeFiles)

	if err != nil {
		s.logger.Errorf("Error getting commands: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commandsStr := make([]string, len(commands))
	for i, command := range commands {
		commandsStr[i] = command.Command
	}

	s.logger.Infof("Found %d commands", len(commandsStr))

	httputils.WriteYamlResponse(w, commandsStr)
}

// func (s *Server) getEnv(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	stack := vars["stack"]
// 	host := vars["host"]

// 	composeFile, err := compose.GetComposeFile(stack, host)

// 	if err != nil {
// 		s.logger.Errorf("Error getting compose file for stack %s and host %s: %v", stack, host, err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	config, err := compose.GetComposeFileConfig(composeFile)
// 	if err != nil {
// 		s.logger.Errorf("Error getting config from compose file for stack %s and host %s: %v", stack, host, err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	httputils.WriteTextResponse(w, config)

// 	s.logger.Infof("Successfully served config for stack %s and host %s", stack, host)
// }
