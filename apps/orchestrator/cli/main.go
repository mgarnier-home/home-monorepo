package main

import (
	"os"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/commands"
)

func main() {
	logger.InitAppLogger("orchestrator-cli")

	rootCommand := &commands.Command{
		Command:     "orchestrator-cli",
		SubCommands: make(map[string]*commands.Command),
	}

	commandsStrings, err := commands.GetCommands()

	if err != nil {
		logger.Errorf("Error getting commands: %v", err)
		os.Exit(1)
	}

	for _, commandString := range commandsStrings {
		commands.SetSubCommands(commandString, rootCommand)
	}

	rootCobraCommand := commands.GetCobraCommand(rootCommand, nil)
	rootCobraCommand.PersistentFlags().String("mode", "", "Choose execution mode (local, hybrid, remote)")
	rootCobraCommand.PersistentFlags().String("service", "", "Execute command for a specific service")

	rootCobraCommand.AddCommand(commands.CompletionCommand())
	rootCobraCommand.AddCommand(commands.UpdateCommand())

	err = rootCobraCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
