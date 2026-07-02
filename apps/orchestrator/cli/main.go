package main

import (
	"os"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/s3"
	"mgarnier11.fr/go/orchestrator-cli/commands"
	"mgarnier11.fr/go/orchestrator-cli/config"
	common "mgarnier11.fr/go/orchestrator-common"
)

func main() {
	logger.InitAppLogger("orchestrator")

	commonLib := common.NewCommonLib(
		config.Env.ComposeDirPath,
		&s3.Config{
			Endpoint:        config.Env.S3Endpoint,
			AccessKeyID:     config.Env.S3AccessKey,
			SecretAccessKey: config.Env.S3SecretKey,
			Bucket:          config.Env.S3Bucket,
		},
	)

	rootCommand := &commands.Command{
		Command: "orchestrator",

		SubCommands: make(map[string]*commands.Command),
	}

	commandsStrings, err := commands.GetCommands(commonLib)

	if err != nil {
		logger.Errorf("Error getting commands: %v", err)
		os.Exit(1)
	}

	for _, commandString := range commandsStrings {
		commands.SetSubCommands(commandString, rootCommand)
	}

	rootCobraCommand := commands.GetCobraCommand(commonLib, rootCommand, nil)
	rootCobraCommand.PersistentFlags().String("mode", "", "Choose execution mode (local, hybrid, remote)")
	rootCobraCommand.PersistentFlags().String("service", "", "Execute command for a specific service")

	rootCobraCommand.AddCommand(commands.CompletionCommand())
	rootCobraCommand.AddCommand(commands.UpdateCliCommand())

	err = rootCobraCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
