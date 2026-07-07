package cli

import (
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"mgarnier11.fr/go/libs/logger"
	apiclient "mgarnier11.fr/go/orchestrator/implementation/apiClient"
	"mgarnier11.fr/go/orchestrator/implementation/execution"
	"mgarnier11.fr/go/orchestrator/interfaces/server"
	"mgarnier11.fr/go/orchestrator/models"
)

var Logger = logger.NewLogger("[CLI]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

type actionCommand struct {
	Name     string
	Commands map[string]*actionCommand
}

// Execution
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

func getActionCommandFunc( /* orchestratorConfig *models.OrchestratorConfig */ ) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		service := cmd.Flag("service").Value.String()

		executionService := execution.GetExecutionService()

		err := executionService.ExecCommand(getCliCommand(cmd), service, os.Stdout)

		if err != nil {
			Logger.Errorf("Error executing command %s: %v", getCliCommand(cmd), err)
			return err
		}

		return nil
	}
}

// Subcommands
func setSubCommands(commandString string, command *actionCommand) {
	parts := strings.Split(commandString, " ")
	mainCmd := parts[0]

	if command.Commands[mainCmd] == nil {
		command.Commands[mainCmd] = &actionCommand{
			Name:     mainCmd,
			Commands: make(map[string]*actionCommand),
		}
	}

	if len(parts) == 1 {
		return
	}

	setSubCommands(strings.Join(parts[1:], " "), command.Commands[mainCmd])
}

func addActionCommands(orchestratorConfig *models.OrchestratorConfig, actionCommand *actionCommand, parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:  actionCommand.Name,
		RunE: getActionCommandFunc( /* orchestratorConfig */ ),
	}

	if parent != nil {
		parent.AddCommand(cmd)
	}

	for _, subCommand := range actionCommand.Commands {
		addActionCommands(orchestratorConfig, subCommand, cmd)
	}

	return cmd
}

// Main cobra commands

func ActionCommands(orchestratorConfig *models.OrchestratorConfig, rootCmd *cobra.Command, commands []*models.Command) error {
	rootActionCommand := &actionCommand{
		Name:     "orchestrator",
		Commands: make(map[string]*actionCommand),
	}

	for _, command := range commands {
		setSubCommands(command.Command, rootActionCommand)
	}

	for _, subCommand := range rootActionCommand.Commands {
		rootCmd.AddCommand(addActionCommands(orchestratorConfig, subCommand, rootCmd))
	}

	rootCmd.PersistentFlags().String("mode", "", "Choose execution mode (local, hybrid, remote)")
	rootCmd.PersistentFlags().String("service", "", "Execute command for a specific service")

	return nil
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

func UpdateCliCommand(orchestratorConfig *models.OrchestratorConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "update-cli",
		Short: "Update the orchestrator",
		RunE: func(cmd *cobra.Command, args []string) error {
			Logger.Infof("Updating orchestrator binary...")

			oldFilePath, err := os.Executable()
			if err != nil {
				Logger.Errorf("Error getting current executable path: %v", err)
				return err
			}

			filePath, err := apiclient.DownloadCliBinary(orchestratorConfig.ApiUrl, runtime.GOARCH, runtime.GOOS)
			if err != nil {
				Logger.Errorf("Error downloading CLI binary: %v", err)
				return err
			}

			err = os.Rename(oldFilePath, oldFilePath+".old")
			if err != nil {
				Logger.Errorf("Error renaming old binary file: %v", err)
				return err
			}

			err = os.Rename(filePath, oldFilePath)
			if err != nil {
				Logger.Errorf("Error renaming new binary file to old binary file path: %v", err)
				return err
			}

			err = os.Chmod(oldFilePath, 0755)
			if err != nil {
				Logger.Errorf("Error setting permissions on new binary file: %v", err)
				return err
			}

			Logger.Infof("Update completed successfully.")

			return nil
		},
	}
}

func StartServerCommand(orchestratorConfig *models.OrchestratorConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "start-server",
		Short: "Start the orchestrator server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.NewServer(orchestratorConfig).Start()
		},
	}
}
