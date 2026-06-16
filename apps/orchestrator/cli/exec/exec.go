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

func formatResults(results map[*types.ComposeConfig]error) string {
	if len(results) == 0 {
		return ""
	}

	resultString := "Results:\n"

	for config, err := range results {
		resultString += fmt.Sprintf("%s %s %s - Error : %v\n", config.Action, config.Host, config.Stack, err)
	}

	return resultString
}

func ExecCommand(commonLib *common.CommonLib, command string, service string) error {
	var configs []*types.ComposeConfig = make([]*types.ComposeConfig, 0)
	var err error
	var results map[*types.ComposeConfig]error

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
		results = commonLib.Exec.ExecCommandsStream(configs, service, nil)
	case config.ModeHybrid:
		Logger.Infof("Getting commands to execute from api... %s", config.Env.ApiUrl)

		configs, err = api.GetComposeConfigs(command)

		if err != nil {
			return fmt.Errorf("error getting compose configs from api: %w", err)
		}

		Logger.Infof("Executing command on local... %s", config.Env.ComposeDirPath)
		results = commonLib.Exec.ExecCommandsStream(configs, service, nil)
	case config.ModeFullApi:
		Logger.Infof("Executing command on api... %s", config.Env.ApiUrl)
		err = api.ExecCommandStream(command, service)
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
