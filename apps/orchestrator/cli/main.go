package main

import (
	"os"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
	"mgarnier11.fr/go/orchestrator-cli/commands"
)

func main() {
	logger.InitAppLogger("orchestrator-cli")

	commandsString, err := api.GetCommands()

	if err != nil {
		logger.Errorf("Error getting commands from orchestrator API: %v", err)
		return
	}

	rootCommand := &commands.Command{
		Command:     "orchestrator-cli",
		SubCommands: make(map[string]*commands.Command),
	}

	for _, commandString := range commandsString {
		commands.SetSubCommands(commandString, rootCommand)
	}

	rootCobraCommand := commands.GetCobraCommand(rootCommand, nil)

	rootCobraCommand.AddCommand(commands.CompletionCommand())
	rootCobraCommand.AddCommand(commands.UpdateCommand())

	err = rootCobraCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
