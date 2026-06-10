package config

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/osutils"
	"mgarnier11.fr/go/libs/s3"
	"mgarnier11.fr/go/orchestrator-common/types"
	"mgarnier11.fr/go/orchestrator-common/utils"
)

type Config struct {
	composeDir string
	s3Config   *s3.Config
	logger     *logger.Logger
}

func NewConfig(composeDir string, s3Config *s3.Config) *Config {
	return &Config{
		composeDir: composeDir,
		s3Config:   s3Config,
		logger:     logger.NewLogger("[COMPOSE-CONFIG]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
	}
}

func (c *Config) GetComposeConfigs(commands []*types.Command) ([]*types.ComposeConfig, error) {
	composeConfigs := []*types.ComposeConfig{}

	for _, command := range commands {
		if command.ComposeFile == nil {
			c.logger.Errorf("Command %s has no compose file", command.Command)
			continue
		}

		osCommand := &osutils.OsCommand{
			OsCommand:     "docker",
			OsCommandArgs: c.getConfigArgs(command.ComposeFile),
			Dir:           c.composeDir,
		}

		configOutput, err := osutils.ExecOsCommandOutput(osCommand, command.Command)
		if err != nil {
			c.logger.Errorf("Error executing command %s %s %s: %v", command.ComposeFile.Stack, command.ComposeFile.Host, command.Action, err)
			continue
		}

		var composeConfig types.ComposeFileSource
		if err := yaml.Unmarshal([]byte(configOutput), &composeConfig); err != nil {
			return nil, fmt.Errorf("error parsing compose config: %w", err)
		}

		osCommand.OsCommandArgs = append(osCommand.OsCommandArgs, "--no-interpolate")
		configOutputNoInterpolation, err := osutils.ExecOsCommandOutput(osCommand, command.Command)
		if err != nil {
			c.logger.Errorf("Error executing command with no interpolation %s %s %s: %v", command.ComposeFile.Stack, command.ComposeFile.Host, command.Action, err)
			continue
		}

		var composeConfigNoInterpolation types.ComposeFileSource
		if err := yaml.Unmarshal([]byte(configOutputNoInterpolation), &composeConfigNoInterpolation); err != nil {
			return nil, fmt.Errorf("error parsing compose config with no interpolation: %w", err)
		}

		for serviceName, service := range composeConfig.Services {
			service.Image = composeConfigNoInterpolation.Services[serviceName].Image

			c.logger.Debugf("Found service %s in config %s %s %s", serviceName, command.ComposeFile.Stack, command.ComposeFile.Host, command.Action)
			c.logger.Debugf("Service %s config: container_name=%s image=%s", serviceName, service.ContainerName, service.Image)
		}

		composeConfigs = append(composeConfigs, &types.ComposeConfig{
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

func (c *Config) getConfigArgs(command *types.ComposeFile) []string {
	args := []string{
		"compose",
	}

	// err := c.DownloadEnvFiles(context.Background(), c.composeDir, command.Stack)
	// if err != nil {
	// 	c.logger.Errorf("Error downloading env files: %v", err)
	// }

	envFiles := utils.GetEnvFiles(c.composeDir, command.Stack)

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

func (c *Config) DownloadEnvFiles(
	ctx context.Context,
	targetDir string,
	stack string,
) error {
	client, err := s3.NewClient(ctx, *c.s3Config)

	if err != nil {
		return fmt.Errorf("error creating S3 client: %w", err)
	}
	objects, err := client.ListObjects(ctx, "")
	if err != nil {
		return fmt.Errorf("error listing objects in bucket: %w", err)
	}

	for _, object := range objects {
		c.logger.Debugf("Object : %s", object.Key)
		if strings.HasPrefix(object.Key, "") && strings.HasSuffix(object.Key, ".env") {
			localPath := path.Join(targetDir, object.Key)
			err := client.DownloadToFile(ctx, object.Key, localPath)
			if err != nil {
				c.logger.Errorf("Error downloading file %s: %v", object.Key, err)
				continue
			}
			c.logger.Infof("Downloaded file %s to %s", object.Key, localPath)
		}
	}

	return nil
}
