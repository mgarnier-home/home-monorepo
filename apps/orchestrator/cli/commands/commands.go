package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
	"mgarnier11.fr/go/orchestrator-cli/config"
	"mgarnier11.fr/go/orchestrator-cli/exec"

	composefiles "mgarnier11.fr/go/orchestrator-common/files"
)

type Command struct {
	Command     string
	SubCommands map[string]*Command
}

var Logger = logger.NewLogger("[CLI-COMMANDS]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

func SetSubCommands(commandString string, command *Command) {
	parts := strings.Split(commandString, " ")
	mainCmd := parts[0]

	if command.SubCommands[mainCmd] == nil {
		command.SubCommands[mainCmd] = &Command{
			Command:     mainCmd,
			SubCommands: make(map[string]*Command),
		}
	}

	if len(parts) == 1 {
		return
	}

	SetSubCommands(strings.Join(parts[1:], " "), command.SubCommands[mainCmd])
}

func getFullCommand(cobraCmd *cobra.Command) string {
	if cobraCmd.Parent() == nil {
		return cobraCmd.Use
	}
	return getFullCommand(cobraCmd.Parent()) + " " + cobraCmd.Use
}

func getCliCommand(cobraCmd *cobra.Command) string {
	fullCmd := getFullCommand(cobraCmd)

	parts := strings.Split(fullCmd, " ")

	return strings.Join(parts[1:], " ")
}

func GetCobraCommand(command *Command, parentCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: command.Command,
		Run: func(cmd *cobra.Command, args []string) {
			service := cmd.Flag("service").Value.String()

			err := exec.ExecCommand(getCliCommand(cmd), service)

			if err != nil {
				Logger.Errorf("Error executing command %s: %v", getCliCommand(cmd), err)
				return
			}
		},
	}

	if parentCmd != nil {
		parentCmd.AddCommand(cmd)
	}

	for _, subCommand := range command.SubCommands {
		cmd.AddCommand(GetCobraCommand(subCommand, cmd))
	}

	return cmd
}

func CompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate completion script",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
		Hidden: true,
	}
}

func UpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the orchestrator-cli",
		Run: func(cmd *cobra.Command, args []string) {
			Logger.Infof("Updating orchestrator-cli...")

			oldFilePath, err := os.Executable()
			if err != nil {
				Logger.Errorf("Error getting current executable path: %v", err)
				return
			}

			filePath, err := api.DownloadCliBinary(runtime.GOARCH, runtime.GOOS)
			if err != nil {
				Logger.Errorf("Error downloading CLI binary: %v", err)
				return
			}

			err = os.Rename(oldFilePath, oldFilePath+".old")
			if err != nil {
				Logger.Errorf("Error renaming old file: %v", err)
				return
			}

			err = os.Rename(filePath, oldFilePath)
			if err != nil {
				Logger.Errorf("Error renaming new file to old file path: %v", err)
				return
			}

			Logger.Infof("Update completed successfully.")
		},
	}
}

func GetCommands() ([]string, error) {
	switch config.Env.Mode {
	case config.ModeFullLocal:
		composeFiles, err := composefiles.GetComposeFiles(config.Env.ComposeDir)

		if err != nil {
			return nil, fmt.Errorf("error getting compose files from local: %w", err)
		}

		commands, err := composefiles.GetCommands(composeFiles)

		if err != nil {
			return nil, fmt.Errorf("error getting commands from local: %w", err)
		}

		commandsString := make([]string, len(commands))

		for i, command := range commands {
			commandsString[i] = command.Command
		}

		return commandsString, nil
	case config.ModeHybrid, config.ModeFullApi:
		commandsString, err := api.GetCommands()

		if err != nil {
			return nil, fmt.Errorf("error getting commands from api: %w", err)
		}

		return commandsString, nil
	}

	return nil, fmt.Errorf("invalid mode: %s", config.Env.Mode)
}
