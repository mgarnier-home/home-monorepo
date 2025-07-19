package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
)

type Command struct {
	Command     string
	SubCommands map[string]*Command
}

func setSubCommands(commandString string, command *Command) {
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

	setSubCommands(strings.Join(parts[1:], " "), command.SubCommands[mainCmd])
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

func getCobraCommand(command *Command, parentCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: command.Command,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Running command: %s with args: %v", getCliCommand(cmd), args)

			err := api.ExecCommandStream(getCliCommand(cmd))

			if err != nil {
				logger.Errorf("Error executing command %s: %v", getCliCommand(cmd), err)
				return
			}
		},
	}

	if parentCmd != nil {
		parentCmd.AddCommand(cmd)
	}

	for _, subCommand := range command.SubCommands {
		cmd.AddCommand(getCobraCommand(subCommand, cmd))
	}

	return cmd
}

func completionCommand() *cobra.Command {
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

func updateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the orchestrator-cli",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("Updating orchestrator-cli...")

			oldFilePath, err := os.Executable()
			if err != nil {
				logger.Errorf("Error getting current executable path: %v", err)
				return
			}

			filePath, err := api.DownloadCliBinary(runtime.GOARCH, runtime.GOOS)
			if err != nil {
				logger.Errorf("Error downloading CLI binary: %v", err)
				return
			}

			err = os.Rename(oldFilePath, oldFilePath+".old")
			if err != nil {
				logger.Errorf("Error renaming old file: %v", err)
				return
			}

			err = os.Rename(filePath, oldFilePath)
			if err != nil {
				logger.Errorf("Error renaming new file to old file path: %v", err)
				return
			}

			logger.Infof("Update completed successfully.")
		},
	}
}

func main() {
	logger.InitAppLogger("orchestrator-cli")

	commandsString, err := api.GetCommands()

	if err != nil {
		logger.Errorf("Error getting commands from orchestrator API: %v", err)
		return
	}

	rootCommand := &Command{
		Command:     "orchestrator-cli",
		SubCommands: make(map[string]*Command),
	}

	for _, commandString := range commandsString {
		setSubCommands(commandString, rootCommand)
	}

	rootCobraCommand := getCobraCommand(rootCommand, nil)

	rootCobraCommand.AddCommand(completionCommand())
	rootCobraCommand.AddCommand(updateCommand())

	err = rootCobraCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
