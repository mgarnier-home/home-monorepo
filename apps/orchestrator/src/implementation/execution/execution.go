package execution

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/osutils"
	apiclient "mgarnier11.fr/go/orchestrator/implementation/apiClient"
	"mgarnier11.fr/go/orchestrator/implementation/command"
	"mgarnier11.fr/go/orchestrator/implementation/compose"
	"mgarnier11.fr/go/orchestrator/implementation/composeConfig.go"
	"mgarnier11.fr/go/orchestrator/models"
)

type ExecutionService struct {
	config *models.OrchestratorConfig
	logger *logger.Logger

	composeService       *compose.ComposeService
	composeConfigService *composeConfig.ComposeConfigService
	commandService       *command.CommandService
}

var (
	instance *ExecutionService
	once     sync.Once
)

func InitExecutionService(config *models.OrchestratorConfig, composeService *compose.ComposeService, composeConfigService *composeConfig.ComposeConfigService, commandService *command.CommandService) *ExecutionService {
	once.Do(func() {
		instance = &ExecutionService{
			config:               config,
			logger:               logger.NewLogger("[SERVICE:EXECUTION]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
			composeService:       composeService,
			composeConfigService: composeConfigService,
			commandService:       commandService,
		}
	})

	return instance
}

func GetExecutionService() *ExecutionService {
	if instance == nil {
		fmt.Println("ExecutionService is not initialized. Please call InitExecutionService first.")
		os.Exit(1)
	}

	return instance
}

func (service *ExecutionService) ExecCommand(command string, targetService string, writer io.Writer) error {
	var configs []*models.ComposeConfig = make([]*models.ComposeConfig, 0)
	var err error
	var results map[*models.ComposeConfig]error

	switch service.config.Mode {
	case models.ModeFullLocal:
		service.logger.Infof("Getting commands to execute from local... %s", service.config.ComposeDirPath)

		commands, err := service.commandService.GetCommandsToExecute(command)
		if err != nil {
			return fmt.Errorf("error getting commands to execute from local: %w", err)
		}

		configs, err = service.composeConfigService.GetComposeConfigs(commands)
		if err != nil {
			return fmt.Errorf("error getting compose configs from local: %w", err)
		}

		service.logger.Infof("Executing command on local... %s", service.config.ComposeDirPath)
		results = service.execCommandsStream(configs, targetService, writer)
	case models.ModeHybrid:
		service.logger.Infof("Getting commands to execute from api... %s", service.config.ApiUrl)

		configs, err = apiclient.GetComposeConfigs(service.config.ApiUrl, command)

		if err != nil {
			return fmt.Errorf("error getting compose configs from api: %w", err)
		}

		service.logger.Infof("Executing command on local... %s", service.config.ComposeDirPath)
		results = service.execCommandsStream(configs, targetService, writer)
	case models.ModeFullApi:
		service.logger.Infof("Executing command on api... %s", service.config.ApiUrl)
		err = apiclient.ExecCommandStream(service.config.ApiUrl, command, targetService)
		if err != nil {
			return fmt.Errorf("error executing command on api: %v", err)
		}
	}

	hasError := false

	for _, err := range results {
		if err != nil {
			hasError = true
		}
	}

	if hasError {
		return fmt.Errorf("error executing command: %s", formatResults(results))
	}

	return nil
}

func (service *ExecutionService) execCommandsStream(composeConfigs []*models.ComposeConfig, targetService string, writer io.Writer) map[*models.ComposeConfig]error {

	results := make(map[*models.ComposeConfig]error)

	// Execute commands for each compose config
	for _, composeConfig := range composeConfigs {
		if targetService != "" && composeConfig.Services[targetService] == nil {
			service.logger.Infof("Skipping config %s %s %s as it does not contain service %s", composeConfig.Host, composeConfig.Stack, composeConfig.Action, targetService)
			continue
		}
		results[composeConfig] = service.execComposeConfigStream(composeConfig, targetService, writer)
	}

	// Reset to default context
	results[&models.ComposeConfig{Host: "default", Stack: "", Action: "context reset"}] = osutils.ExecOsCommandStream(&osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "use", "default"},
		Dir:           os.TempDir(),
	}, writer, "docker context use default")

	// Log results
	for config, err := range results {
		if err != nil {
			log := color.RedString("%s %s %s - Error : %v", config.Action, config.Host, config.Stack, err)
			service.logger.Errorf("%s", log)
			if writer != nil {
				writer.Write([]byte(fmt.Sprintf("%s\n", log)))
			}
		} else {
			log := color.GreenString("%s %s %s - Success", config.Action, config.Host, config.Stack)
			service.logger.Infof("%s", log)
			if writer != nil {
				writer.Write([]byte(fmt.Sprintf("%s\n", log)))
			}
		}
	}

	return results
}

func (service *ExecutionService) execComposeConfigStream(config *models.ComposeConfig, targetService string, writer io.Writer) error {
	service.logger.Infof("Executing %s %s %s %s", config.Action, config.Host, config.Stack, targetService)

	// Write the config to a file
	filePath, err := service.writeComposeConfigToTempFile(config.Config)
	if err != nil {
		return fmt.Errorf("error writing compose config to file for host %s: %w", config.Host, err)
	}

	service.logger.Debugf("Compose config written to file: %s", filePath)

	// Delete the file after execution
	defer os.Remove(filePath)

	// Create a context for the host
	if err := service.setContextStream(config, writer); err != nil {
		return fmt.Errorf("error setting context for host %s: %w", config.Host, err)
	}

	// Execute the compose command using the file and context
	if err := service.execComposeCommandStream(config, filePath, targetService, writer); err != nil {
		return fmt.Errorf("error executing compose command for host %s: %w", config.Host, err)
	}

	return nil
}

func (service *ExecutionService) writeComposeConfigToTempFile(config string) (string, error) {
	file, err := os.CreateTemp("", "compose-*.yml")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(config)
	if err != nil {
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}

	return file.Name(), nil
}

func (service *ExecutionService) setContextStream(config *models.ComposeConfig, writer io.Writer) error {
	dockerContextCreateCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "create", config.Host, "--docker", "host=" + config.HostConfig},
	}

	err := osutils.ExecOsCommandStream(dockerContextCreateCommand, writer, "docker context create "+config.Host)
	if err != nil {
		service.logger.Debugf("Context %s already exists, skipping creation", config.Host)
	}

	dockerContextUseCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "use", config.Host},
	}

	err = osutils.ExecOsCommandStream(dockerContextUseCommand, writer, "docker context use "+config.Host)

	if err != nil {
		return err
	}

	return nil
}

func (service *ExecutionService) execComposeCommandStream(
	config *models.ComposeConfig,
	composeFileName string,
	targetService string,
	writer io.Writer,
) error {
	args := []string{
		"compose",
		"-f", composeFileName,
	}

	switch config.Action {
	case "up":
		args = append(args, "up", "--remove-orphans", "-d", "--pull", "always")
	case "down":
		args = append(args, "down", "--remove-orphans", "-v")
	case "restart":
		args = append(args, "up", "--remove-orphans", "-d", "--pull", "always", "--force-recreate")
	default:
		return fmt.Errorf("unknown action: %s", config.Action)
	}

	if targetService != "" {
		args = append(args, targetService)
	}

	osCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: args,
		Dir:           os.TempDir(),
	}
	return osutils.ExecOsCommandStream(osCommand, writer, fmt.Sprintf("docker %s %s %s", config.Action, config.Host, config.Stack))
}

func formatResults(results map[*models.ComposeConfig]error) string {
	if len(results) == 0 {
		return ""
	}

	resultString := "Results:\n"

	for config, err := range results {
		resultString += fmt.Sprintf("%s %s %s - Error : %v\n", config.Action, config.Host, config.Stack, err)
	}

	return resultString
}
