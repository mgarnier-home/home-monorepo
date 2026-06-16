package exec

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
	"mgarnier11.fr/go/orchestrator-cli/config"
	common "mgarnier11.fr/go/orchestrator-common"
	"mgarnier11.fr/go/orchestrator-common/types"
)

var Logger = logger.NewLogger("[CLI-EXEC]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func formatErrors(errors map[*types.ComposeConfig]error) string {
	if len(errors) == 0 {
		return ""
	}

	errorString := "Errors:\n"

	for config, err := range errors {
		errorString += fmt.Sprintf("%s %s %s - Error : %v\n", config.Action, config.Host, config.Stack, err)
	}

	return errorString
}

func ExecCommand(commonLib *common.CommonLib, command string, service string) error {
	var configs []*types.ComposeConfig = make([]*types.ComposeConfig, 0)
	var err error

	switch config.Env.Mode {
	case config.ModeFullLocal:
		Logger.Infof("Getting commands to execute from local... %s", config.Env.ComposeDirPath)

		commands, err := commonLib.Files.GetCommandsToExecute(command)

		if err != nil {
			return fmt.Errorf("error getting commands to execute from local: %w", err)
		}

		configs, err = commonLib.Config.GetComposeConfigs(commands)
		if err != nil {
			return fmt.Errorf("error getting compose configs from local: %w", err)
		}

		Logger.Infof("Executing command on local... %s", config.Env.ComposeDirPath)
		errors := commonLib.Exec.ExecCommandsStream(configs, service, nil)
		if len(errors) > 0 {
			return fmt.Errorf("error executing command on local: %v", formatErrors(errors))
		}
	case config.ModeHybrid:
		Logger.Infof("Getting commands to execute from api... %s", config.Env.ApiUrl)

		configs, err = api.GetComposeConfigs(command)

		if err != nil {
			return fmt.Errorf("error getting compose configs from api: %w", err)
		}

		Logger.Infof("Executing command on local... %s", config.Env.ComposeDirPath)
		errors := commonLib.Exec.ExecCommandsStream(configs, service, nil)
		if len(errors) > 0 {
			return fmt.Errorf("error executing command on local: %v", formatErrors(errors))
		}
	case config.ModeFullApi:
		Logger.Infof("Executing command on api... %s", config.Env.ApiUrl)
		err = api.ExecCommandStream(command, service)
		if err != nil {
			return fmt.Errorf("error executing command on api: %v", err)
		}
	}

	return nil
}
