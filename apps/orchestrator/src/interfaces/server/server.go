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
	"mgarnier11.fr/go/orchestrator/implementation/command"
	"mgarnier11.fr/go/orchestrator/implementation/compose"
	"mgarnier11.fr/go/orchestrator/implementation/composeConfig.go"
	"mgarnier11.fr/go/orchestrator/implementation/execution"
	"mgarnier11.fr/go/orchestrator/models"
)

type Server struct {
	logger *logger.Logger
	port   int

	config *models.OrchestratorConfig

	composeService       *compose.ComposeService
	executionService     *execution.ExecutionService
	commandService       *command.CommandService
	composeConfigService *composeConfig.ComposeConfigService
}

func NewServer(
	orchestratorConfig *models.OrchestratorConfig,
) *Server {
	return &Server{
		logger:               logger.NewLogger("[SERVER]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
		port:                 orchestratorConfig.ServerPort,
		config:               orchestratorConfig,
		composeService:       compose.GetComposeService(),
		executionService:     execution.GetExecutionService(),
		commandService:       command.GetCommandService(),
		composeConfigService: composeConfig.GetComposeConfigService(),
	}
}

func (server *Server) Start() error {

	router := mux.NewRouter()

	router.Use(httputils.LogRequestMiddleware(server.logger))
	router.Use(httputils.CorsMiddleware)

	version.SetupVersionRoute(router)

	server.logger.Infof("Starting server on port %d", server.port)

	router.HandleFunc("/", server.get).Methods("GET")
	router.HandleFunc("/cli", server.getCli).Methods("GET")
	router.HandleFunc("/compose", server.getComposeFiles).Methods("GET")
	router.HandleFunc("/commands", server.getCommands).Methods("GET")
	router.HandleFunc("/{command}/exec", server.streamExecCommand).Methods("GET")
	router.HandleFunc("/{command}/configs", server.getCommandsConfigs).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%d", server.port), router)
}

func (server *Server) getComposeFiles(w http.ResponseWriter, r *http.Request) {
	if err := server.composeService.RefreshComposeFiles(); err != nil {
		server.logger.Errorf("Error refreshing compose files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	composeFiles := server.composeService.GetComposeFiles()

	server.logger.Infof("Found %d compose files", len(composeFiles))

	httputils.WriteYamlResponse(w, composeFiles)

	server.logger.Infof("Successfully served %d compose files", len(composeFiles))

}

func (server *Server) streamExecCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	command := vars["command"]

	if command == "" {
		server.logger.Errorf("No command provided in request")
		http.Error(w, "Bad Request: No command provided", http.StatusBadRequest)
		return
	}
	service := r.URL.Query().Get("service")

	server.logger.Debugf("Received command: %s %s", command, service)

	w.Header().Set("Content-Type", "text/plain")

	err := server.executionService.ExecCommand(command, service, w)
	if err != nil {
		server.logger.Errorf("Error executing command %s: %v", command, err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (server *Server) getCommands(w http.ResponseWriter, r *http.Request) {
	if err := server.composeService.RefreshComposeFiles(); err != nil {
		server.logger.Errorf("Error refreshing commands: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commands := server.commandService.GetCommands()

	commandsStr := make([]string, len(commands))
	for i, command := range commands {
		commandsStr[i] = command.Command
	}

	server.logger.Infof("Found %d commands", len(commandsStr))

	httputils.WriteYamlResponse(w, commandsStr)
}

func (server *Server) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Orchestrator API is running"))
	server.logger.Infof("Handled GET request to /")
}

func (server *Server) getCli(w http.ResponseWriter, r *http.Request) {
	arch := r.URL.Query().Get("arch")
	osName := r.URL.Query().Get("os")

	fileName := "orchestrator"
	if arch != "" && osName != "" {
		fileName += fmt.Sprintf("-%s-%s", osName, arch)

		if osName == "windows" {
			fileName += ".exe"
		}
	}

	binPath := path.Join(server.config.BinariesPath, fileName)

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		server.logger.Errorf("CLI binary not found at %s", binPath)
		http.Error(w, fmt.Sprintf("CLI binary not found: %s", binPath), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, binPath)
	server.logger.Infof("Served CLI binary from %s", binPath)
}

func (server *Server) getCommandsConfigs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	command := vars["command"]

	if command == "" {
		server.logger.Errorf("No command provided in request")
		http.Error(w, "Bad Request: No command provided", http.StatusBadRequest)
		return
	}

	server.logger.Debugf("Received command: %s", command)

	commandsToExecute, err := server.commandService.GetCommandsToExecute(command)
	if err != nil {
		server.logger.Errorf("Error getting commands to execute: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	server.logger.Debugf("Found %d commands to execute", len(commandsToExecute))

	composeConfigs, err := server.composeConfigService.GetComposeConfigs(commandsToExecute)
	if err != nil {
		server.logger.Errorf("Error getting compose configs: %v", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	server.logger.Debugf("Found %d compose configs", len(composeConfigs))

	httputils.WriteYamlResponse(w, composeConfigs)
}
