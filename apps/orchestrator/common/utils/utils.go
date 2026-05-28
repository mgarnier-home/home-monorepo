package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"mgarnier11.fr/go/libs/osutils"
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

func UpdateRepo(writer io.Writer, folder string, token string) error {
	gitEnv := []string{"GIT_TERMINAL_PROMPT=0"}

	gitToken := strings.TrimSpace(token)
	if gitToken != "" {
		basicAuth := base64.StdEncoding.EncodeToString([]byte("x-access-token:" + gitToken))
		gitEnv = append(gitEnv,
			"GIT_CONFIG_COUNT=1",
			"GIT_CONFIG_KEY_0=http.extraheader",
			fmt.Sprintf("GIT_CONFIG_VALUE_0=AUTHORIZATION: basic %s", basicAuth),
		)
	}

	err := osutils.ExecOsCommandStream(&osutils.OsCommand{
		OsCommand:     "git",
		OsCommandArgs: []string{"pull"},
		Dir:           folder,
		Env:           gitEnv,
	}, writer, "git pull")
	if err != nil {
		return fmt.Errorf("git pull failed in %s: %w", folder, err)
	}

	return nil
}
