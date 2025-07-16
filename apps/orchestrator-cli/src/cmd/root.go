package cmd

import (
	"os"

	"mgarnier11.fr/go/orchestrator-cli/command"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "home-cli",
}

func Execute() {

	rootCmd.AddCommand(command.GetCobraCommands()...)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.CompletionOptions.DisableDescriptions = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
