package compose

import (
	"fmt"
	"os"
	"path"
	"strings"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-api/config"
	osUtils "mgarnier11.fr/go/orchestrator-api/os-utils"
)

type ComposeFile struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Host  string `yaml:"host"`
	Stack string `yaml:"stack"`
}

func GetComposeFile(stackName, hostName string) (*ComposeFile, error) {
	logger.Debugf("Getting compose file for stack: %s and host: %s", stackName, hostName)

	filePath := path.Join(config.Env.ComposeDir, stackName, fmt.Sprintf("%s.%s.yml", hostName, stackName))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("compose file %s does not exist for stack %s and host %s", filePath, stackName, hostName)
	}

	return &ComposeFile{
		Name:  hostName + "-" + stackName,
		Path:  filePath,
		Host:  hostName,
		Stack: stackName,
	}, nil
}

func GetStackComposeFiles(stackName string) ([]*ComposeFile, error) {
	composeFiles := []*ComposeFile{}

	stackPath := path.Join(config.Env.ComposeDir, stackName)

	hostsFiles, err := os.ReadDir(stackPath)

	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", stackPath, err)
	}

	for _, hostFile := range hostsFiles {
		hostFileName := hostFile.Name()

		if !strings.HasSuffix(hostFileName, stackName+".yml") && !strings.HasSuffix(hostFileName, stackName+".yaml") {
			continue
		}
		hostName := strings.Split(hostFileName, ".")[0]

		logger.Debugf("Found compose file: %s for host: %s for stack: %s", hostFileName, hostName, stackName)

		composeFile, err := GetComposeFile(stackName, hostName)

		if err != nil {
			logger.Errorf("Error getting compose file for stack %s and host %s: %v", stackName, hostName, err)
			continue
		}

		composeFiles = append(composeFiles, composeFile)
	}

	return composeFiles, nil
}

func GetComposeFiles() ([]*ComposeFile, error) {
	logger.Infof("Getting compose files from directory: %s", config.Env.ComposeDir)

	stacks, err := os.ReadDir(config.Env.ComposeDir)
	if err != nil {
		return nil, err
	}

	composeFiles := []*ComposeFile{}

	for _, stack := range stacks {

		if !stack.IsDir() {
			continue
		}

		logger.Debugf("Found stack: %s", stack.Name())

		stackComposeFiles, err := GetStackComposeFiles(stack.Name())

		if err != nil {
			logger.Errorf("Error getting compose files for stack %s: %v", stack.Name(), err)
			continue
		}

		composeFiles = append(composeFiles, stackComposeFiles...)
	}

	return composeFiles, nil
}

func GetComposeFileConfig(file *ComposeFile) (string, error) {
	args := []string{
		"compose",
	}

	envFiles := getEnvFiles(file.Stack)

	for _, envFile := range envFiles {
		args = append(args, "--env-file", envFile)
	}

	args = append(args,
		"-f",
		fmt.Sprintf("%s/%s.%s.yml", file.Stack, file.Host, file.Stack),
	)

	args = append(args, "config", "--format", "yaml")

	dockerComposeCommand := &osUtils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: args,
		Dir:           config.Env.ComposeDir,
	}

	return osUtils.ExecOsCommandOutput(dockerComposeCommand, fmt.Sprintf("%s %s", file.Stack, file.Host))
}

func getComposeCommandArgs(command *ComposeFile) []string {
	args := []string{
		"compose",
	}

	envFiles := getEnvFiles(command.Stack)

	for _, envFile := range envFiles {
		args = append(args, "--env-file", envFile)
	}

	args = append(args,
		"-f",
		fmt.Sprintf("%s/%s.%s.yml", command.Stack, command.Host, command.Stack),
	)

	return args
}

func getEnvFiles(stack string) []string {
	composeDir := config.Env.ComposeDir
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
