package compose

import (
	"fmt"
	"os"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"mgarnier11.fr/go/libs/logger"
	common "mgarnier11.fr/go/orchestrator-common"
)

var Logger = logger.NewLogger("[COMPOSE-FILES]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func composeFile(composeDir, stackName, hostName string) *common.ComposeFile {
	stackPath := path.Join(composeDir, stackName)

	return &common.ComposeFile{
		Name:  hostName + "-" + stackName,
		Stack: stackName,
		Host:  hostName,
		Path:  path.Join(stackPath, hostName+"."+stackName+".yml"),
	}

}

func GetComposeFiles(composeDir string) ([]*common.ComposeFile, error) {
	Logger.Debugf("Getting compose files from directory: %s", composeDir)

	stacks, err := os.ReadDir(composeDir)
	if err != nil {
		return nil, err
	}

	composeFiles := []*common.ComposeFile{}

	for _, stack := range stacks {

		if !stack.IsDir() {
			continue
		}

		stackName := stack.Name()

		Logger.Verbosef("Found stack: %s", stackName)

		stackPath := path.Join(composeDir, stackName)
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

			composeFiles = append(composeFiles, composeFile(composeDir, stackName, hostName))
		}
	}

	return composeFiles, nil
}

func GetCommands(composeFiles []*common.ComposeFile) ([]*common.Command, error) {
	hosts := []string{}
	for _, composeFile := range composeFiles {
		if slices.Contains(hosts, composeFile.Host) {
			continue
		}
		hosts = append(hosts, composeFile.Host)
	}

	// Foreach stack, generate commands ${stack} ${host} ${action} and ${stack} ${action} ${host} and ${stack} ${action}
	commands := []*common.Command{}
	for _, composeFile := range composeFiles {
		for _, action := range common.ActionList {
			// Command for specific host
			commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s %s", composeFile.Stack, composeFile.Host, action), ComposeFile: composeFile, Action: action})
			commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s %s", composeFile.Stack, action, composeFile.Host), ComposeFile: composeFile, Action: action})
			// Command for all hosts in stack
			commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s", composeFile.Stack, action), ComposeFile: &common.ComposeFile{Stack: composeFile.Stack, Host: "all", Name: "all-" + composeFile.Stack}, Action: action})
		}
	}

	// Foreach host, generate commands ${host} ${stack} ${action} and ${host} ${action} ${stack} and ${host} ${action}
	for _, host := range hosts {
		for _, composeFile := range composeFiles {
			if composeFile.Host != host {
				continue
			}
			for _, action := range common.ActionList {
				// Command for specific stack
				commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s %s", host, action, composeFile.Stack), ComposeFile: composeFile, Action: action})
				commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s %s", host, composeFile.Stack, action), ComposeFile: composeFile, Action: action})
				// Command for all stacks
				commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s", host, action), ComposeFile: &common.ComposeFile{Host: host, Name: host + "-all", Stack: "all"}, Action: action})
			}
		}
	}

	// Foreach action, generate commands ${action} ${stack} ${host} and ${action} ${host} ${stack} and ${action} ${stack} and ${action} ${host}
	for _, action := range common.ActionList {
		for _, composeFile := range composeFiles {
			// Command for specific stack and host
			commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s %s", action, composeFile.Stack, composeFile.Host), ComposeFile: composeFile, Action: action})
			commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s %s", action, composeFile.Host, composeFile.Stack), ComposeFile: composeFile, Action: action})
			// Command for all stacks and hosts
			commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s", action, composeFile.Stack), ComposeFile: &common.ComposeFile{Stack: composeFile.Stack, Name: "all-" + composeFile.Stack, Host: "all"}, Action: action})
		}

		for _, host := range hosts {
			// Command for specific host
			commands = append(commands, &common.Command{Command: fmt.Sprintf("%s %s", action, host), ComposeFile: &common.ComposeFile{Host: host, Name: host + "-all", Stack: "all"}, Action: action})
		}

		commands = append(commands, &common.Command{Command: action, ComposeFile: &common.ComposeFile{Host: "all", Name: "all-all", Stack: "all"}, Action: action}) // Command for all actions
	}

	return commands, nil
}

func GetCommandsToExecute(composeDir, commandString string) ([]*common.Command, error) {
	composeFiles, err := GetComposeFiles(composeDir)

	if err != nil {
		return nil, fmt.Errorf("error getting compose files: %w", err)
	}

	allCommands, err := GetCommands(composeFiles)

	if err != nil {
		return nil, fmt.Errorf("error getting commands: %w", err)
	}

	commandIndex := slices.IndexFunc(allCommands, func(c *common.Command) bool {
		return c.Command == commandString
	})

	if commandIndex == -1 {
		return nil, fmt.Errorf("command %s not found", commandString)
	}

	command := allCommands[commandIndex]

	commands := []*common.Command{}

	for _, composeFile := range composeFiles {
		if (command.ComposeFile.Host == "all" && command.ComposeFile.Stack == "all") ||
			(command.ComposeFile.Host == "all" && command.ComposeFile.Stack == composeFile.Stack) ||
			(command.ComposeFile.Host == composeFile.Host && command.ComposeFile.Stack == "all") ||
			(command.ComposeFile.Host == composeFile.Host && command.ComposeFile.Stack == composeFile.Stack) {
			commands = append(commands, &common.Command{
				Command:     fmt.Sprintf("%s %s %s", composeFile.Stack, composeFile.Host, command.Action),
				ComposeFile: composeFile,
				Action:      command.Action,
			})
		}
	}

	sort.SliceStable(commands, func(i, j int) bool {
		if commands[i].ComposeFile.Host == commands[j].ComposeFile.Host {
			return commands[i].ComposeFile.Stack < commands[j].ComposeFile.Stack
		}
		return commands[i].ComposeFile.Host < commands[j].ComposeFile.Host
	})

	return commands, nil
}
