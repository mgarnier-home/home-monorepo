package composeConfig

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/osutils"
	"mgarnier11.fr/go/libs/utils"
	"mgarnier11.fr/go/orchestrator/implementation/compose"
	"mgarnier11.fr/go/orchestrator/implementation/env"
	"mgarnier11.fr/go/orchestrator/models"
)

type ComposeConfigService struct {
	logger *logger.Logger

	envService     *env.EnvService
	composeService *compose.ComposeService
}

var (
	instance *ComposeConfigService
)

func InitComposeConfigService(envService *env.EnvService, composeService *compose.ComposeService) *ComposeConfigService {
	instance = &ComposeConfigService{
		logger:         logger.NewLogger("[SERVICE:COMPOSE-CONFIG]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
		envService:     envService,
		composeService: composeService,
	}
	return instance
}

func GetComposeConfigService() *ComposeConfigService {
	if instance == nil {
		fmt.Println("ComposeConfigService is not initialized. Please call InitComposeConfigService first.")
		os.Exit(1)
	}

	return instance
}

func (service *ComposeConfigService) GetComposeConfigs(commands []*models.Command) ([]*models.ComposeConfig, error) {
	composeConfigs := []*models.ComposeConfig{}

	for _, command := range commands {
		if command.ComposeFile == nil {
			service.logger.Errorf("Command %s has no compose file", command.Command)
			continue
		}

		args := []string{
			"compose",
		}
		envFiles := getEnvFiles(service.envService.GetEnvFilesDir(), command.ComposeFile.Stack)
		for _, envFile := range envFiles {
			args = append(args, "--env-file", envFile)
		}
		args = append(args,
			"-f",
			fmt.Sprintf("%s/%s.%s.yml", command.ComposeFile.Stack, command.ComposeFile.Host, command.ComposeFile.Stack),
			"config",
		)

		osCommand := &osutils.OsCommand{
			OsCommand:     "docker",
			OsCommandArgs: args,
			Dir:           service.composeService.GetComposeFilesDir(),
		}

		configOutput, err := osutils.ExecOsCommandOutput(osCommand)
		if err != nil {
			service.logger.Errorf("Error executing command %s %s %s: %v", command.ComposeFile.Stack, command.ComposeFile.Host, command.Action, err)
			continue
		}

		var composeConfig models.ComposeFileSource
		if err := yaml.Unmarshal([]byte(configOutput), &composeConfig); err != nil {
			return nil, fmt.Errorf("error parsing compose config: %w", err)
		}

		composeConfigs = append(composeConfigs, &models.ComposeConfig{
			Host:       command.ComposeFile.Host,
			Stack:      command.ComposeFile.Stack,
			Action:     command.Action,
			Config:     configOutput,
			HostConfig: getHostConfig(command.ComposeFile.Host),
			Services:   composeConfig.Services,
		})
	}

	return composeConfigs, nil

}

func getEnvFiles(dir string, stack string) []string {
	globalEnvFiles := getEnvFilesPaths(dir)
	stackEnvFiles := getEnvFilesPaths(fmt.Sprintf("%s/%s", dir, stack))

	return append(globalEnvFiles, stackEnvFiles...)
}

func getEnvFilesPaths(dir string) []string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []string{}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	envFiles := []string{}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".env") {
			envFiles = append(envFiles, fmt.Sprintf("%s/%s", dir, file.Name()))
		}
	}

	return envFiles
}

func getHostConfig(host string) string {
	hostConfig := utils.GetEnv(strings.ToUpper(host)+"_HOST", "")

	if hostConfig == "" {
		panic(fmt.Sprintf("Host config for %s not found in environment variables", host))
	}
	return hostConfig
}
