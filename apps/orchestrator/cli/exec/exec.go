package exec

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
	"mgarnier11.fr/go/orchestrator-cli/config"
	compose_common "mgarnier11.fr/go/orchestrator-common"
	compose_config "mgarnier11.fr/go/orchestrator-common/config"
	compose_exec "mgarnier11.fr/go/orchestrator-common/exec"
	compose_files "mgarnier11.fr/go/orchestrator-common/files"
)

var Logger = logger.NewLogger("[CLI-EXEC]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func ExecCommand(command string, service string) error {
	var configs []*compose_common.ComposeConfig = make([]*compose_common.ComposeConfig, 0)
	var err error

	switch config.Env.Mode {
	case config.ModeFullLocal:
		Logger.Infof("Getting commands to execute from local... %s", config.Env.ComposeDir)

		commands, err := compose_files.GetCommandsToExecute(config.Env.ComposeDir, command)

		if err != nil {
			return fmt.Errorf("error getting commands to execute from local: %w", err)
		}

		configs, err = compose_config.GetComposeConfigs(config.Env.ComposeDir, commands)
		if err != nil {
			return fmt.Errorf("error getting compose configs from local: %w", err)
		}
	case config.ModeHybrid:
		Logger.Infof("Getting commands to execute from api... %s", config.Env.ApiUrl)

		configs, err = api.GetComposeConfigs(command)

		if err != nil {
			return fmt.Errorf("error getting compose configs from api: %w", err)
		}
	}

	switch config.Env.Mode {
	case config.ModeFullLocal, config.ModeHybrid:
		Logger.Infof("Executing command on local... %s", config.Env.ComposeDir)
		compose_exec.ExecCommandsStream(configs, service, nil)
	case config.ModeFullApi:
		Logger.Infof("Executing command on api... %s", config.Env.ApiUrl)
		api.ExecCommandStream(command, service)

	}

	return nil
}
