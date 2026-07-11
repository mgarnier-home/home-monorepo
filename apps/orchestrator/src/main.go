package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mgarnier11.fr/go/libs/config"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/osutils"
	"mgarnier11.fr/go/libs/s3"
	"mgarnier11.fr/go/orchestrator/implementation/command"
	"mgarnier11.fr/go/orchestrator/implementation/compose"
	"mgarnier11.fr/go/orchestrator/implementation/composeConfig.go"
	"mgarnier11.fr/go/orchestrator/implementation/env"
	"mgarnier11.fr/go/orchestrator/implementation/execution"
	"mgarnier11.fr/go/orchestrator/interfaces/cli"
	"mgarnier11.fr/go/orchestrator/models"
)

func main() {
	logger.InitAppLogger("orchestrator")

	orchestratorConfig := &models.OrchestratorConfig{}

	errors := config.GetConfig(orchestratorConfig)

	if len(errors) > 0 {
		for _, err := range errors {
			logger.Errorf("Error loading config: %v", err)
		}
		os.Exit(1)
	}

	if orchestratorConfig.SSHPrivateKey != "" {
		err := createSSHKeyFile(orchestratorConfig.SSHPrivateKey)
		if err != nil {
			logger.Errorf("Error creating SSH key file: %v", err)
			os.Exit(1)
		}
	}

	composeService := compose.InitComposeService(orchestratorConfig)
	defer composeService.Destroy()
	commandService := command.InitCommandService(orchestratorConfig, composeService)
	envService := env.InitEnvService(&s3.Config{
		Endpoint:        orchestratorConfig.S3Endpoint,
		AccessKeyID:     orchestratorConfig.S3AccessKey,
		SecretAccessKey: orchestratorConfig.S3SecretKey,
		Bucket:          orchestratorConfig.S3Bucket,
		Region:          "nantes",
	})
	defer envService.Destroy()
	composeConfigService := composeConfig.InitComposeConfigService(envService, composeService)
	execution.InitExecutionService(orchestratorConfig, composeService, composeConfigService, commandService)

	rootCommand := &cobra.Command{
		Use: "orchestrator",
	}

	cli.ActionCommands(orchestratorConfig, rootCommand)

	rootCommand.AddCommand(cli.CompletionCommand())
	rootCommand.AddCommand(cli.UpdateCliCommand(orchestratorConfig))
	rootCommand.AddCommand(cli.StartServerCommand(orchestratorConfig))

	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func createSSHKeyFile(sshPrivateKey string) error {

	// Create /ssh directory if it doesn't exist
	err := os.MkdirAll("ssh", 0700)
	if err != nil {
		return fmt.Errorf("error creating ssh directory: %w", err)
	}

	// If the id_rsa file already exists, do nothing
	exists, err := osutils.FileExists("ssh/ssh_private_key")
	if err != nil {
		return fmt.Errorf("error checking if SSH private key file exists: %w", err)
	}
	if !exists {
		// Create the id_rsa file with the content of SSH_PRIVATE_KEY
		err = os.WriteFile("ssh/ssh_private_key", []byte(sshPrivateKey), 0600)
		if err != nil {
			return fmt.Errorf("error writing SSH private key to file: %w", err)
		}
	}

	return nil
}
