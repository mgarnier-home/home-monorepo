package command

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator/implementation/compose"
	"mgarnier11.fr/go/orchestrator/models"
)

type CommandService struct {
	config         *models.OrchestratorConfig
	logger         *logger.Logger
	composeService *compose.ComposeService

	commands []*models.Command
}

var (
	instance *CommandService
)

func InitCommandService(config *models.OrchestratorConfig, composeService *compose.ComposeService) *CommandService {

	instance = &CommandService{
		config:         config,
		logger:         logger.NewLogger("[SERVICE:COMMAND]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
		composeService: composeService,
		commands:       []*models.Command{},
	}

	err := instance.RefreshCommands()

	if err != nil {
		instance.logger.Errorf("Error initializing command service: %v", err)
		os.Exit(1)
	}
	return instance
}

func GetCommandService() *CommandService {
	if instance == nil {
		fmt.Println("CommandService is not initialized. Please call InitCommandService first.")
		os.Exit(1)
	}
	return instance
}

func (service *CommandService) GetCommands() []*models.Command {
	return service.commands
}

func (service *CommandService) RefreshCommands() error {
	switch service.config.Mode {
	case models.ModeFullLocal:
		service.commands = getCommandsFromComposeFiles(service.composeService.GetComposeFiles())
		return nil
	case models.ModeHybrid, models.ModeFullApi:
		commandsString, err := getCommandsFromApi()

		if err != nil {
			return fmt.Errorf("error getting commands from api: %w", err)
		}

		service.commands = commandsString

		return nil
	}

	return fmt.Errorf("invalid mode: %s", service.config.Mode)
}

func getCommandsFromComposeFiles(composeFiles []*models.ComposeFile) []*models.Command {
	hosts := []string{}
	for _, composeFile := range composeFiles {
		if slices.Contains(hosts, composeFile.Host) {
			continue
		}
		hosts = append(hosts, composeFile.Host)
	}

	// Foreach stack, generate commands ${stack} ${host} ${action} and ${stack} ${action} ${host} and ${stack} ${action}
	commands := []*models.Command{}
	for _, composeFile := range composeFiles {
		for _, action := range models.ActionList {
			// Command for specific host
			commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s %s", composeFile.Stack, composeFile.Host, action), ComposeFile: composeFile, Action: action})
			commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s %s", composeFile.Stack, action, composeFile.Host), ComposeFile: composeFile, Action: action})
			// Command for all hosts in stack
			commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s", composeFile.Stack, action), ComposeFile: &models.ComposeFile{Stack: composeFile.Stack, Host: "all", Name: "all-" + composeFile.Stack}, Action: action})
		}
	}

	// Foreach host, generate commands ${host} ${stack} ${action} and ${host} ${action} ${stack} and ${host} ${action}
	for _, host := range hosts {
		for _, composeFile := range composeFiles {
			if composeFile.Host != host {
				continue
			}
			for _, action := range models.ActionList {
				// Command for specific stack
				commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s %s", host, action, composeFile.Stack), ComposeFile: composeFile, Action: action})
				commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s %s", host, composeFile.Stack, action), ComposeFile: composeFile, Action: action})
				// Command for all stacks
				commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s", host, action), ComposeFile: &models.ComposeFile{Host: host, Name: host + "-all", Stack: "all"}, Action: action})
			}
		}
	}

	// Foreach action, generate commands ${action} ${stack} ${host} and ${action} ${host} ${stack} and ${action} ${stack} and ${action} ${host}
	for _, action := range models.ActionList {
		for _, composeFile := range composeFiles {
			// Command for specific stack and host
			commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s %s", action, composeFile.Stack, composeFile.Host), ComposeFile: composeFile, Action: action})
			commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s %s", action, composeFile.Host, composeFile.Stack), ComposeFile: composeFile, Action: action})
			// Command for all stacks and hosts
			commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s", action, composeFile.Stack), ComposeFile: &models.ComposeFile{Stack: composeFile.Stack, Name: "all-" + composeFile.Stack, Host: "all"}, Action: action})
		}

		for _, host := range hosts {
			// Command for specific host
			commands = append(commands, &models.Command{Command: fmt.Sprintf("%s %s", action, host), ComposeFile: &models.ComposeFile{Host: host, Name: host + "-all", Stack: "all"}, Action: action})
		}

		commands = append(commands, &models.Command{Command: action, ComposeFile: &models.ComposeFile{Host: "all", Name: "all-all", Stack: "all"}, Action: action}) // Command for all actions
	}

	return commands
}

func getCommandsFromApi() ([]*models.Command, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (service *CommandService) GetCommandsToExecute(commandString string) ([]*models.Command, error) {
	allCommands := service.GetCommands()
	composeFiles := service.composeService.GetComposeFiles()

	commandIndex := slices.IndexFunc(allCommands, func(c *models.Command) bool {
		return c.Command == commandString
	})

	if commandIndex == -1 {
		return nil, fmt.Errorf("command %s not found", commandString)
	}

	command := allCommands[commandIndex]

	commands := []*models.Command{}

	for _, composeFile := range composeFiles {
		if (command.ComposeFile.Host == "all" && command.ComposeFile.Stack == "all") ||
			(command.ComposeFile.Host == "all" && command.ComposeFile.Stack == composeFile.Stack) ||
			(command.ComposeFile.Host == composeFile.Host && command.ComposeFile.Stack == "all") ||
			(command.ComposeFile.Host == composeFile.Host && command.ComposeFile.Stack == composeFile.Stack) {
			commands = append(commands, &models.Command{
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
