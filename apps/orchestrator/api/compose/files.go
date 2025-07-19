package compose

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"mgarnier11.fr/go/orchestrator-api/config"
)

type ComposeFile struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Host  string `yaml:"host"`
	Stack string `yaml:"stack"`
}

func composeFile(stackName, hostName string) *ComposeFile {
	stackPath := path.Join(config.Env.ComposeDir, stackName)

	return &ComposeFile{
		Name:  hostName + "-" + stackName,
		Stack: stackName,
		Host:  hostName,
		Path:  path.Join(stackPath, hostName+"."+stackName+".yml"),
	}

}

func GetComposeFiles() ([]*ComposeFile, error) {
	Logger.Infof("Getting compose files from directory: %s", config.Env.ComposeDir)

	stacks, err := os.ReadDir(config.Env.ComposeDir)
	if err != nil {
		return nil, err
	}

	composeFiles := []*ComposeFile{}

	for _, stack := range stacks {

		if !stack.IsDir() {
			continue
		}

		stackName := stack.Name()

		Logger.Verbosef("Found stack: %s", stackName)

		stackPath := path.Join(config.Env.ComposeDir, stackName)
		hostsFiles, err := os.ReadDir(stackPath)

		if err != nil {
			Logger.Errorf("Error reading directory %s: %v", stackPath, err)
			continue
		}

		for _, hostFile := range hostsFiles {
			hostFileName := hostFile.Name()

			parts := strings.Split(hostFileName, ".")

			if len(parts) != 3 || parts[1] != stackName || parts[2] != "yml" {
				continue
			}

			hostName := parts[0]

			Logger.Verbosef("Found compose file: %s for host: %s for stack: %s", hostFileName, hostName, stackName)

			composeFiles = append(composeFiles, composeFile(stackName, hostName))
		}
	}

	return composeFiles, nil
}

func GetCommands(composeFiles []*ComposeFile) ([]*Command, error) {
	hosts := []string{}
	for _, composeFile := range composeFiles {
		if slices.Contains(hosts, composeFile.Host) {
			continue
		}
		hosts = append(hosts, composeFile.Host)
	}

	// Foreach stack, generate commands ${stack} ${host} ${action} and ${stack} ${action} ${host} and ${stack} ${action}
	commands := []*Command{}
	for _, composeFile := range composeFiles {
		for _, action := range ActionList {
			// Command for specific host
			commands = append(commands, &Command{Command: fmt.Sprintf("%s %s %s", composeFile.Stack, composeFile.Host, action), ComposeFile: composeFile, Action: action})
			commands = append(commands, &Command{Command: fmt.Sprintf("%s %s %s", composeFile.Stack, action, composeFile.Host), ComposeFile: composeFile, Action: action})
			// Command for all hosts in stack
			commands = append(commands, &Command{Command: fmt.Sprintf("%s %s", composeFile.Stack, action), ComposeFile: &ComposeFile{Stack: composeFile.Stack, Host: "all", Name: "all-" + composeFile.Stack}, Action: action})
		}
	}

	// Foreach host, generate commands ${host} ${stack} ${action} and ${host} ${action} ${stack} and ${host} ${action}
	for _, host := range hosts {
		for _, composeFile := range composeFiles {
			if composeFile.Host != host {
				continue
			}
			for _, action := range ActionList {
				// Command for specific stack
				commands = append(commands, &Command{Command: fmt.Sprintf("%s %s %s", host, action, composeFile.Stack), ComposeFile: composeFile, Action: action})
				commands = append(commands, &Command{Command: fmt.Sprintf("%s %s %s", host, composeFile.Stack, action), ComposeFile: composeFile, Action: action})
				// Command for all stacks
				commands = append(commands, &Command{Command: fmt.Sprintf("%s %s", host, action), ComposeFile: &ComposeFile{Host: host, Name: host + "-all", Stack: "all"}, Action: action})
			}
		}
	}

	// Foreach action, generate commands ${action} ${stack} ${host} and ${action} ${host} ${stack} and ${action} ${stack} and ${action} ${host}
	for _, action := range ActionList {
		for _, composeFile := range composeFiles {
			// Command for specific stack and host
			commands = append(commands, &Command{Command: fmt.Sprintf("%s %s %s", action, composeFile.Stack, composeFile.Host), ComposeFile: composeFile, Action: action})
			commands = append(commands, &Command{Command: fmt.Sprintf("%s %s %s", action, composeFile.Host, composeFile.Stack), ComposeFile: composeFile, Action: action})
			// Command for all stacks and hosts
			commands = append(commands, &Command{Command: fmt.Sprintf("%s %s", action, composeFile.Stack), ComposeFile: &ComposeFile{Stack: composeFile.Stack, Name: "all-" + composeFile.Stack, Host: "all"}, Action: action})
		}

		for _, host := range hosts {
			// Command for specific host
			commands = append(commands, &Command{Command: fmt.Sprintf("%s %s", action, host), ComposeFile: &ComposeFile{Host: host, Name: host + "-all", Stack: "all"}, Action: action})
		}

		commands = append(commands, &Command{Command: action, ComposeFile: &ComposeFile{Host: "all", Name: "all-all", Stack: "all"}, Action: action}) // Command for all actions
	}

	return commands, nil
}
