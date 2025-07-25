package exec

import (
	"fmt"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-cli/api"
)

func ExecCommand(command string) error {
	logger.Infof("Running command: %s", command)

	// err := api.ExecCommandStream(getCliCommand(cmd))

	// if err != nil {
	// 	logger.Errorf("Error executing command %s: %v", getCliCommand(cmd), err)
	// 	return
	// }

	configs, err := api.GetComposeConfigs(command)

	if err != nil {
		logger.Errorf("Error getting compose configs for command %s: %v", command, err)
		return fmt.Errorf("error getting compose configs: %w", err)
	}

	for _, config := range configs {
		logger.Infof("Got compose config for host %s, action %s", config.Host, config.Action)

		// Write the config to a file
		// Create a context for the host
		// Execute the compose command using the file and context
		// Delete the file after execution
		// Reset the context

	}

	return nil
}
