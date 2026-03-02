package config

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/osutils"
	common "mgarnier11.fr/go/orchestrator-common"
	"mgarnier11.fr/go/orchestrator-common/utils"
)

var Logger = logger.NewLogger("[COMPOSE-CONFIG]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func GetComposeConfigs(composeDir string, commands []*common.Command) ([]*common.ComposeConfig, error) {
	composeConfigs := []*common.ComposeConfig{}

	for _, command := range commands {
		if command.ComposeFile == nil {
			Logger.Errorf("Command %s has no compose file", command.Command)
			continue
		}

		osCommand := &osutils.OsCommand{
			OsCommand:     "docker",
			OsCommandArgs: getConfigArgs(composeDir, command.ComposeFile),
			Dir:           composeDir,
		}

		configOutput, err := osutils.ExecOsCommandOutput(osCommand, command.Command)
		if err != nil {
			Logger.Errorf("Error executing command %s %s %s: %v", command.ComposeFile.Stack, command.ComposeFile.Host, command.Action, err)
			continue
		}

		var composeConfig common.ComposeFileSource
		if err := yaml.Unmarshal([]byte(configOutput), &composeConfig); err != nil {
			return nil, fmt.Errorf("error parsing compose config: %w", err)
		}

		composeConfigs = append(composeConfigs, &common.ComposeConfig{
			Host:       command.ComposeFile.Host,
			Stack:      command.ComposeFile.Stack,
			Action:     command.Action,
			Config:     configOutput,
			HostConfig: utils.GetHostConfig(command.ComposeFile.Host),
			Services:   composeConfig.Services,
		})
	}

	return composeConfigs, nil

}

func getConfigArgs(composeDir string, command *common.ComposeFile) []string {
	args := []string{
		"compose",
	}

	envFiles := utils.GetEnvFiles(composeDir, command.Stack)

	for _, envFile := range envFiles {
		args = append(args, "--env-file", envFile)
	}

	args = append(args,
		"-f",
		fmt.Sprintf("%s/%s.%s.yml", command.Stack, command.Host, command.Stack),
		"config",
	)

	return args
}
