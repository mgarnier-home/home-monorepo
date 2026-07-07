package apiclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"go.yaml.in/yaml/v3"
	"mgarnier11.fr/go/orchestrator/models"
)

func DownloadCliBinary(apiUrl, arch, osName string) (string, error) {
	url := fmt.Sprintf("%s/cli?arch=%s&os=%s", apiUrl, arch, osName)
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
		return "", fmt.Errorf("error getting current executable path: %w", err)
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

func GetComposeConfigs(apiUrl, command string) ([]*models.ComposeConfig, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/configs", apiUrl, command))
	if err != nil {
		return nil, fmt.Errorf("error getting compose configs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting compose configs: %s", resp.Status)
	}

	var configs []*models.ComposeConfig
	if err := yaml.NewDecoder(resp.Body).Decode(&configs); err != nil {
		return nil, fmt.Errorf("error decoding compose configs response: %w", err)
	}

	return configs, nil
}

func ExecCommandStream(apiUrl, command, service string) error {
	url := fmt.Sprintf("%s/%s/exec", apiUrl, command)
	if service != "" {
		url = fmt.Sprintf("%s?service=%s", url, service)
	}
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
