package cmd

import (
	"fmt"
	"mgarnier11/home-cli/command"
	"mgarnier11/home-cli/utils"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var commandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "Get all commands",
	Run: func(cmd *cobra.Command, args []string) {
		commandsPaths := utils.GetSubCommandsPaths(command.GetCobraCommands())

		preciseStackControlCommands := []string{}
		stacksControlCommands := []string{}
		hostControlCommands := []string{}
		actionCommands := []string{}

		for _, commandPath := range commandsPaths {
			commandParts := strings.Split(commandPath, " ")

			switch len(commandParts) {
			case 1:
				actionCommands = append(actionCommands, commandPath)
			case 2:
				if slices.Contains(utils.StackList, commandParts[0]) {
					stacksControlCommands = append(stacksControlCommands, commandPath)
				} else if slices.Contains(utils.HostList, commandParts[0]) {
					hostControlCommands = append(hostControlCommands, commandPath)
				}
			case 3:
				preciseStackControlCommands = append(preciseStackControlCommands, commandPath)
			}
		}

		fmt.Println("===========Precise stack control commands===========")
		fmt.Println(strings.Join(preciseStackControlCommands, "\n"))

		fmt.Println("===========Stacks control commands===========")
		fmt.Println(strings.Join(stacksControlCommands, "\n"))

		fmt.Println("===========Host control commands===========")
		fmt.Println(strings.Join(hostControlCommands, "\n"))

		fmt.Println("===========Action commands===========")
		fmt.Println(strings.Join(actionCommands, "\n"))

	},
}

func init() {
	commandsCmd.Hidden = true
	rootCmd.AddCommand(commandsCmd)
}
