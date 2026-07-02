package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/utils"
)

var _logger = logger.NewLogger("[COMMON-UTILS]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func GetHostConfig(host string) string {
	hostConfig := utils.GetEnv(strings.ToUpper(host)+"_HOST", "")

	if hostConfig == "" {
		panic(fmt.Sprintf("Host config for %s not found in environment variables", host))
	}
	return hostConfig
}

func GetEnvFiles(dir string, stack string) []string {
	globalEnvFiles := getEnvFilesPaths(dir)
	stackEnvFiles := getEnvFilesPaths(fmt.Sprintf("%s/%s", dir, stack))

	return append(globalEnvFiles, stackEnvFiles...)
}

func getEnvFilesPaths(dir string) []string {
	_logger.Debugf("Getting env files from directory: %s", dir)

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
