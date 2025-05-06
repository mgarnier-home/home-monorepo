package command

import (
	"fmt"
	"mgarnier11/home-cli/utils"

	"github.com/spf13/cobra"
)

func createActionsCommands(stacks []string, hosts []string) []*cobra.Command {
	actionsCmd := []*cobra.Command{}

	for _, action := range utils.ActionList {
		actionCmd := &cobra.Command{
			Use:       action,
			ValidArgs: []string{"parallel"},
			Args:      cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Running action", cmd.CommandPath(), "with args ", args)
				fmt.Println("host", hosts)
				fmt.Println("stack", stacks)
				fmt.Println("action", action)

				ExecCommand(stacks, hosts, action, args)
			},
		}

		actionsCmd = append(actionsCmd, actionCmd)
	}

	return actionsCmd
}

func GetCobraCommands() []*cobra.Command {
	compose := utils.GetCompose()
	stacksCmd := []*cobra.Command{}

	for stack, hosts := range compose.Stacks {
		stackCmd := &cobra.Command{
			Use: stack,
		}

		for _, host := range hosts {
			hostCmd := &cobra.Command{
				Use: host,
			}

			hostCmd.AddCommand(createActionsCommands([]string{stack}, []string{host})...)
			stackCmd.AddCommand(hostCmd)
		}

		stackCmd.AddCommand(createActionsCommands([]string{stack}, hosts)...)
		stacksCmd = append(stacksCmd, stackCmd)
	}

	for host, stacks := range compose.Hosts {
		hostCmd := &cobra.Command{
			Use: host,
		}

		for _, stack := range stacks {
			stackCmd := &cobra.Command{
				Use: stack,
			}

			stackCmd.AddCommand(createActionsCommands([]string{stack}, []string{host})...)
			hostCmd.AddCommand(stackCmd)
		}

		hostCmd.AddCommand(createActionsCommands(stacks, []string{host})...)
		stacksCmd = append(stacksCmd, hostCmd)
	}

	stacksCmd = append(stacksCmd, createActionsCommands(utils.StackList, utils.HostList)...)

	return stacksCmd
}
