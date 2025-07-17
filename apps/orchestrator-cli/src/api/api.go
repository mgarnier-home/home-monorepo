package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
	"mgarnier11.fr/go/orchestrator-cli/config"
)

func GetCommands() ([]string, error) {
	// Make a request to the orchestrator API to get the commands
	resp, err := http.Get(config.Env.OrchestratorApiUrl + "/commands")

	if err != nil {
		return nil, fmt.Errorf("error getting commands from orchestrator API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting commands from orchestrator API: %s", resp.Status)
	}
	defer resp.Body.Close()

	var commands []string
	if err := yaml.NewDecoder(resp.Body).Decode(&commands); err != nil {
		return nil, fmt.Errorf("error decoding commands response: %w", err)
	}

	return commands, nil
}

func ExecCommandStream(command string) error {
	url := fmt.Sprintf("%s/exec-command/%s", config.Env.OrchestratorApiUrl, command)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error executing command: %s", resp.Status)
	}

	// Stream response to stdout
	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}
