package compose

import (
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/osutils"
	"mgarnier11.fr/go/libs/utils"
	"mgarnier11.fr/go/orchestrator-api/config"
)

var Logger = logger.NewLogger("[COMPOSE]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

type Command struct {
	Command     string       `yaml:"command"`
	ComposeFile *ComposeFile `yaml:"compose_file"`
	Action      string       `yaml:"action"`
}

type ComposeConfig struct {
	Host   string `yaml:"host"`
	Action string `yaml:"action"`
	Config string `yaml:"config"`
}

var ActionList = []string{"up", "down", "restart"}

func GetCommandsToExecute(commandString string) ([]*Command, error) {

	composeFiles, err := GetComposeFiles()

	if err != nil {
		return nil, fmt.Errorf("error getting compose files: %w", err)
	}

	allCommands, err := GetCommands(composeFiles)

	if err != nil {
		return nil, fmt.Errorf("error getting commands: %w", err)
	}

	commandIndex := slices.IndexFunc(allCommands, func(c *Command) bool {
		return c.Command == commandString
	})

	if commandIndex == -1 {
		return nil, fmt.Errorf("command %s not found", commandString)
	}

	command := allCommands[commandIndex]

	commands := []*Command{}

	for _, composeFile := range composeFiles {
		if (command.ComposeFile.Host == "all" && command.ComposeFile.Stack == "all") ||
			(command.ComposeFile.Host == "all" && command.ComposeFile.Stack == composeFile.Stack) ||
			(command.ComposeFile.Host == composeFile.Host && command.ComposeFile.Stack == "all") ||
			(command.ComposeFile.Host == composeFile.Host && command.ComposeFile.Stack == composeFile.Stack) {
			commands = append(commands, &Command{
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

func GetComposeConfigs(commands []*Command) ([]*ComposeConfig, error) {
	composeConfigs := []*ComposeConfig{}

	for _, command := range commands {
		if command.ComposeFile == nil {
			Logger.Errorf("Command %s has no compose file", command.Command)
			continue
		}

		osCommand := &osutils.OsCommand{
			OsCommand:     "docker",
			OsCommandArgs: getComposeCommandArgs(command.ComposeFile, "config"),
			Dir:           config.Env.ComposeDir,
		}

		configOutput, err := osutils.ExecOsCommandOutput(osCommand, command.Command)
		if err != nil {
			Logger.Errorf("Error executing command %s %s %s: %v", command.ComposeFile.Stack, command.ComposeFile.Host, command.Action, err)
			continue
		}

		composeConfigs = append(composeConfigs, &ComposeConfig{
			Host:   command.ComposeFile.Host,
			Action: command.Action,
			Config: configOutput,
		})

		Logger.Debugf("Compose config for command %s: %s", command.Command, configOutput)
	}

	return composeConfigs, nil

}

func ExecCommandsStream(commands []*Command, writer io.Writer) {
	exec := func(command *Command) error {
		Logger.Infof("Executing command: %s %s %s", command.ComposeFile.Stack, command.ComposeFile.Host, command.Action)

		err := setContext(command.ComposeFile.Host, writer)

		if err != nil {
			return fmt.Errorf("error setting context for host %s: %w", command.ComposeFile.Host, err)
		}

		osCommand := &osutils.OsCommand{
			OsCommand:     "docker",
			OsCommandArgs: getComposeCommandArgs(command.ComposeFile, command.Action),
			Dir:           config.Env.ComposeDir,
		}

		err = osutils.ExecOsCommandStream(osCommand, writer, command.Command)

		if err != nil {
			return fmt.Errorf("error executing command %s %s %s: %w", command.ComposeFile.Stack, command.ComposeFile.Host, command.Action, err)
		}

		return nil
	}

	results := make(map[*Command]error)

	for _, command := range commands {
		results[command] = exec(command)
	}

	results[&Command{Command: "docker context use default"}] = osutils.ExecOsCommandStream(&osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "use", "default"},
		Dir:           config.Env.ComposeDir,
	}, writer, "docker context use default")

	for cmd, err := range results {
		if err != nil {
			log := color.RedString("%s - Error : %v", cmd.Command, err)

			Logger.Errorf("%s", log)
			writer.Write([]byte(fmt.Sprintf("%s\n", log)))
		} else {
			log := color.GreenString("%s - Success", cmd.Command)

			Logger.Infof("%s", log)
			writer.Write([]byte(fmt.Sprintf("%s\n", log)))
		}
	}
}

func setContext(host string, writer io.Writer) error {
	dockerContextCreateCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "create", host, "--docker", "host=" + utils.GetEnv(strings.ToUpper(host)+"_HOST", "")},
		Dir:           config.Env.ComposeDir,
	}

	err := osutils.ExecOsCommandStream(dockerContextCreateCommand, writer, "docker context create "+host)
	if err != nil {
		Logger.Debugf("Context %s already exists, skipping creation", host)
	}

	dockerContextUseCommand := &osutils.OsCommand{
		OsCommand:     "docker",
		OsCommandArgs: []string{"context", "use", host},
		Dir:           config.Env.ComposeDir,
	}

	err = osutils.ExecOsCommandStream(dockerContextUseCommand, writer, "docker context use "+host)

	if err != nil {
		return err
	}

	return nil
}

func getComposeCommandArgs(command *ComposeFile, action string) []string {
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

	switch action {
	case "config":
		args = append(args, "config")
	case "up":
		args = append(args, "up", "-d", "--pull", "always")
	case "down":
		args = append(args, "down", "-v")
	case "restart":
		args = append(args, "up", "-d", "--pull", "always", "--force-recreate")
	}

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
