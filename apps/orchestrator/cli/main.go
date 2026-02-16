package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
	"mgarnier11.fr/go/orchestrator-cli/commands"
	"mgarnier11.fr/go/orchestrator-cli/config"
	compose "mgarnier11.fr/go/orchestrator-common"
)

func main() {
	logger.InitAppLogger("orchestrator-cli")

	rootCommand := &commands.Command{
		Command:     "orchestrator-cli",
		SubCommands: make(map[string]*commands.Command),
	}

	var local bool
	var err error

	bootstrap := pflag.NewFlagSet("bootstrap", pflag.ContinueOnError)
	bootstrap.SetOutput(io.Discard) // silence pre-parse usage/errors
	bootstrap.ParseErrorsWhitelist.UnknownFlags = true
	bootstrap.BoolVar(&local, "local", config.Env.Local, "Execute command locally")
	_ = bootstrap.Parse(os.Args[1:]) // local is now available here

	commandsStrings, err := getCommands(local)

	if err != nil {
		logger.Errorf("Error getting commands: %v", err)
		os.Exit(1)
	}

	for _, commandString := range commandsStrings {
		commands.SetSubCommands(commandString, rootCommand)
	}

	rootCobraCommand := commands.GetCobraCommand(rootCommand, nil)
	rootCobraCommand.PersistentFlags().BoolVar(&local, "local", config.Env.Local, "Execute command locally")

	rootCobraCommand.AddCommand(commands.CompletionCommand())
	rootCobraCommand.AddCommand(commands.UpdateCommand())

	err = rootCobraCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getCommands(local bool) ([]string, error) {

	if local {
		logger.Infof("Executing in local mode, fetching commands from %s", config.Env.ComposeDir)

		composeFiles, err := compose.GetComposeFiles(config.Env.ComposeDir)

		if err != nil {
			return nil, fmt.Errorf("error getting compose files from local: %w", err)
		}

		commands, err := compose.GetCommands(composeFiles)

		if err != nil {
			return nil, fmt.Errorf("error getting commands from local: %w", err)
		}

		commandsString := make([]string, len(commands))

		for i, command := range commands {
			commandsString[i] = command.Command
		}

		return commandsString, nil
	} else {

		commandsString, err := api.GetCommands()

		if err != nil {
			return nil, fmt.Errorf("error getting commands from api: %w", err)
		}

		return commandsString, nil
	}
}
