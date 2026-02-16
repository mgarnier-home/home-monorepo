package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"mgarnier11.fr/go/orchestrator-cli/config"
	compose "mgarnier11.fr/go/orchestrator-common"
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
	url := fmt.Sprintf("%s/%s/exec", config.Env.OrchestratorApiUrl, command)
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

func DownloadCliBinary(arch, osName string) (string, error) {
	url := fmt.Sprintf("%s/cli?arch=%s&os=%s", config.Env.OrchestratorApiUrl, arch, osName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading CLI binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error downloading CLI binary: %s", resp.Status)
	}

	fileName := fmt.Sprintf("orchestrator-cli-%s-%s", osName, arch)
	if osName == "windows" {
		fileName += ".exe"
	}

	binPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	filePath := path.Join(filepath.Dir(binPath), fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	return filePath, nil
}

func GetComposeConfigs(command string) ([]*compose.ComposeConfig, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/configs", config.Env.OrchestratorApiUrl, command))
	if err != nil {
		return nil, fmt.Errorf("error getting compose configs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting compose configs: %s", resp.Status)
	}

	var configs []*compose.ComposeConfig
	if err := yaml.NewDecoder(resp.Body).Decode(&configs); err != nil {
		return nil, fmt.Errorf("error decoding compose configs response: %w", err)
	}

	return configs, nil
}
