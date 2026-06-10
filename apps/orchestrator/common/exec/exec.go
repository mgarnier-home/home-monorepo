package exec

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/osutils"
	"mgarnier11.fr/go/orchestrator-common/types"
)

type Exec struct {
	logger *logger.Logger
}

func NewExec() *Exec {
	return &Exec{
		logger: logger.NewLogger("[COMPOSE-EXEC]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
	}
}

func (e *Exec) ExecCommandsStream(composeConfigs []*types.ComposeConfig, service string, writer io.Writer) map[*types.ComposeConfig]error {

	results := make(map[*types.ComposeConfig]error)

	for _, composeConfig := range composeConfigs {
		if service != "" && composeConfig.Services[service] == nil {
			e.logger.Infof("Skipping config %s %s %s as it does not contain service %s", composeConfig.Host, composeConfig.Stack, composeConfig.Action, service)
			continue
		}
		results[composeConfig] = e.execComposeConfigStream(composeConfig, service, writer)
	}

	// Reset to default context
	results[&types.ComposeConfig{Host: "default", Stack: "", Action: "context reset"}] = osutils.ExecOsCommandStream(&osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "use", "default"},
		Dir:           os.TempDir(),
	}, writer, "docker context use default")

	for config, err := range results {
		if err != nil {
			log := color.RedString("%s %s %s - Error : %v", config.Action, config.Host, config.Stack, err)
			e.logger.Errorf("%s", log)
			if writer != nil {
				writer.Write([]byte(fmt.Sprintf("%s\n", log)))
			}
		} else {
			log := color.GreenString("%s %s %s - Success", config.Action, config.Host, config.Stack)
			e.logger.Infof("%s", log)
			if writer != nil {
				writer.Write([]byte(fmt.Sprintf("%s\n", log)))
			}
		}
	}

	return results
}

func (e *Exec) execComposeConfigStream(config *types.ComposeConfig, service string, writer io.Writer) error {
	e.logger.Infof("Executing %s %s %s %s", config.Action, config.Host, config.Stack, service)

	// Write the config to a file
	filePath, err := e.writeComposeConfigToTempFile(config.Config)
	if err != nil {
		return fmt.Errorf("error writing compose config to file for host %s: %w", config.Host, err)
	}

	e.logger.Debugf("Compose config written to file: %s", filePath)

	// Delete the file after execution
	defer os.Remove(filePath)

	// Create a context for the host
	if err := e.setContextStream(config, writer); err != nil {
		return fmt.Errorf("error setting context for host %s: %w", config.Host, err)
	}

	// Execute the compose command using the file and context
	if err := e.execComposeCommandStream(config, filePath, service, writer); err != nil {
		return fmt.Errorf("error executing compose command for host %s: %w", config.Host, err)
	}

	return nil
}

func (e *Exec) writeComposeConfigToTempFile(config string) (string, error) {
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

func (e *Exec) setContextStream(config *types.ComposeConfig, writer io.Writer) error {
	dockerContextCreateCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "create", config.Host, "--docker", "host=" + config.HostConfig},
	}

	err := osutils.ExecOsCommandStream(dockerContextCreateCommand, writer, "docker context create "+config.Host)
	if err != nil {
		e.logger.Debugf("Context %s already exists, skipping creation", config.Host)
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

func (e *Exec) execComposeCommandStream(
	config *types.ComposeConfig,
	composeFileName string,
	service string,
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

	if service != "" {
		args = append(args, service)
	}

	osCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: args,
		Dir:           os.TempDir(),
	}
	return osutils.ExecOsCommandStream(osCommand, writer, fmt.Sprintf("docker %s %s %s", config.Action, config.Host, config.Stack))
}
