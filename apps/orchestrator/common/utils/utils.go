package utils

import (
	"fmt"
	"os"
	"strings"

	"mgarnier11.fr/go/libs/utils"
)

func GetHostConfig(host string) string {
	return utils.GetEnv(strings.ToUpper(host)+"_HOST", "")
}

func GetEnvFiles(composeDir string, stack string) []string {
	globalEnvFiles := getEnvFilesPaths(composeDir, "")
	stackEnvFiles := getEnvFilesPaths(composeDir, stack)

	return append(globalEnvFiles, stackEnvFiles...)
}

func getEnvFilesPaths(composeDir string, stack string) []string {
	dir := ""
	if stack != "" {
		dir = fmt.Sprintf("%s/%s", composeDir, stack)
	} else {
		dir = composeDir
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	envFiles := []string{}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".env") {
			if stack != "" {
				envFiles = append(envFiles, fmt.Sprintf("%s/%s", stack, file.Name()))
			} else {
				envFiles = append(envFiles, file.Name())
			}
		}
	}

	return envFiles
}
