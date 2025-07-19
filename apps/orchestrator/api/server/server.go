package server

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
	"mgarnier11.fr/go/libs/httputils"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/version"
	"mgarnier11.fr/go/orchestrator-api/compose"
	"mgarnier11.fr/go/orchestrator-api/config"
)

type Server struct {
	port int
}

var Logger = logger.NewLogger("[SERVER]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {

	router := mux.NewRouter()

	router.Use(httputils.LogRequestMiddleware(Logger))
	router.Use(httputils.CorsMiddleware)

	version.SetupVersionRoute(router)

	Logger.Infof("Starting server on port %d", s.port)

	router.HandleFunc("/cli", s.getCli).Methods("GET")

	router.HandleFunc("/compose", s.getComposeFiles).Methods("GET")
	router.HandleFunc("/commands", s.getCommands).Methods("GET")
	router.HandleFunc("/exec-command/{command}", s.streamExecCommand).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)
}

func (s *Server) getComposeFiles(w http.ResponseWriter, r *http.Request) {
	composeFiles, err := compose.GetComposeFiles()

	if err != nil {
		Logger.Errorf("Error getting compose files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Logger.Infof("Found %d compose files", len(composeFiles))

	httputils.WriteYamlResponse(w, composeFiles)

	Logger.Infof("Successfully served %d compose files", len(composeFiles))

}

func (s *Server) streamExecCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	command := vars["command"]

	if command == "" {
		Logger.Errorf("No command provided in request")
		http.Error(w, "Bad Request: No command provided", http.StatusBadRequest)
		return
	}

	Logger.Debugf("Received command: %s", command)

	commandsToExecute, err := compose.GetCommandsToExecute(command)
	if err != nil {
		Logger.Errorf("Error getting commands to execute: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	Logger.Debugf("Found %d commands to execute", len(commandsToExecute))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	compose.ExecCommandsStream(commandsToExecute, w)
}

func (s *Server) getCommands(w http.ResponseWriter, r *http.Request) {
	composeFiles, err := compose.GetComposeFiles()

	if err != nil {
		Logger.Errorf("Error getting compose files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commands, err := compose.GetCommands(composeFiles)

	if err != nil {
		Logger.Errorf("Error getting commands: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commandsStr := make([]string, len(commands))
	for i, command := range commands {
		commandsStr[i] = command.Command
	}

	Logger.Infof("Found %d commands", len(commandsStr))

	httputils.WriteYamlResponse(w, commandsStr)
}

func (s *Server) getCli(w http.ResponseWriter, r *http.Request) {
	arch := r.URL.Query().Get("arch")
	osName := r.URL.Query().Get("os")

	fileName := "orchestrator-cli"
	if arch != "" && osName != "" {
		fileName += fmt.Sprintf("-%s-%s", osName, arch)

		if osName == "windows" {
			fileName += ".exe"
		}
	}

	binPath := path.Join(config.Env.BinariesPath, fileName)

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		Logger.Errorf("CLI binary not found at %s", binPath)
		http.Error(w, fmt.Sprintf("CLI binary not found: %s", binPath), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, binPath)
	Logger.Infof("Served CLI binary from %s", binPath)
}
