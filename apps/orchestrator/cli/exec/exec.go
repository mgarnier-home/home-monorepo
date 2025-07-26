package exec

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"mgarnier11.fr/go/libs/osutils"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
)

var Logger = logger.NewLogger("[EXEC]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func ExecCommand(command string) error {
	Logger.Infof("Running command: %s", command)

	configs, err := api.GetComposeConfigs(command)

	if err != nil {
		return fmt.Errorf("error getting compose configs: %w", err)
	}

	results := make(map[*api.ComposeConfig]error)

	for _, config := range configs {
		results[config] = execComposeConfig(config)
	}

	err = osutils.ExecOsCommand(&osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "use", "default"},
		Dir:           os.TempDir(),
	}, "docker context use default")

	if err != nil {
		Logger.Errorf("error resetting docker context to default: %v", err)
	}

	for config, err := range results {
		if err != nil {
			Logger.Errorf("%s", color.RedString("%s %s %s - Error : %v", config.Host, config.Stack, config.Action, err))
		} else {
			Logger.Infof("%s", color.GreenString("%s %s %s - Success", config.Host, config.Stack, config.Action))
		}
	}

	return nil

}

func execComposeConfig(config *api.ComposeConfig) error {
	Logger.Infof("Executing %s %s %s", config.Action, config.Host, config.Stack)

	// Write the config to a file
	filePath, err := writeComposeConfigToFile(config.Config)
	if err != nil {
		return fmt.Errorf("error writing compose config to file for host %s: %w", config.Host, err)
	}

	Logger.Debugf("Compose config written to file: %s", filePath)

	// Delete the file after execution
	defer os.Remove(filePath)

	// Create a context for the host
	if err := setConfigContext(config, os.Stdout); err != nil {
		return fmt.Errorf("error setting context for host %s: %w", config.Host, err)
	}

	// Execute the compose command using the file and context
	if err := execComposeCommand(config, filePath); err != nil {
		return fmt.Errorf("error executing compose command for host %s: %w", config.Host, err)
	}

	return nil
}

func writeComposeConfigToFile(config string) (string, error) {
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

func setConfigContext(config *api.ComposeConfig, writer io.Writer) error {
	dockerContextCreateCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "create", config.Host, "--docker", "host=" + config.HostConfig},
		Dir:           os.TempDir(),
	}

	err := osutils.ExecOsCommand(dockerContextCreateCommand, "docker context create "+config.Host)
	if err != nil {
		Logger.Debugf("Context %s already exists, skipping creation", config.Host)
	}

	dockerContextUseCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "use", config.Host},
		Dir:           os.TempDir(),
	}

	err = osutils.ExecOsCommand(dockerContextUseCommand, "docker context use "+config.Host)

	if err != nil {
		return err
	}

	return nil
}

func execComposeCommand(config *api.ComposeConfig, composeFileName string) error {
	args := []string{
		"compose",
		"-f", composeFileName,
	}

	switch config.Action {
	case "up":
		args = append(args, "up", "-d", "--pull", "always")
	case "down":
		args = append(args, "down", "-v")
	case "restart":
		args = append(args, "up", "-d", "--pull", "always", "--force-recreate")
	default:
		return fmt.Errorf("unknown action: %s", config.Action)
	}

	osCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: args,
		Dir:           os.TempDir(),
	}

	return osutils.ExecOsCommand(osCommand, fmt.Sprintf("docker compose %s on host %s", config.Action, config.Host))
}
