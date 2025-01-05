package validation

import (
	"encoding/json"
	"errors"
	"mgarnier11/mineager/config"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type postServerRequest struct {
	HostName string `json:"hostName"`
	Name     string `json:"name"`
	Version  string `json:"version",omitempty`
	MapName  string `json:"mapName"`
	Memory   string `json:"memory"`
	Url      string `json:"url"`
}

const minMemory = 1
const maxMemory = 16

func validateMemory(memoryStr string) error {
	if memoryStr == "" {
		return errors.New("memory is required")
	}

	memoryRegex := regexp.MustCompile(`^\d+G$`)
	if !memoryRegex.MatchString(memoryStr) {
		return errors.New("memory must follow the format xG (e.g., 1G)")
	}

	memoryNum := memoryStr[:len(memoryStr)-1]
	memory, err := strconv.Atoi(memoryNum)
	if err != nil {
		return errors.New("failed to parse memory")
	}

	if memory < minMemory || memory > maxMemory {
		return errors.New("memory must be between 1G and 16G")
	}

	return nil
}

func validateUrl(url string) error {
	if url == "" {
		return errors.New("url is required")
	}

	// Regex to validate the URL format
	urlRegex := regexp.MustCompile(`^[a-zA-Z0-9-]+\.mgarnier11\.fr$`)
	if !urlRegex.MatchString(url) {
		return errors.New("url must be in the format WHATEVER.mgarnier11.fr and WHATEVER cannot be empty")
	}

	return nil
}

func validateHostName(hostName string) error {
	if hostName == "" {
		return errors.New("hostName is required")
	}

	_, err := config.GetHost(hostName)

	if err != nil {
		return errors.New("hostName does not exist")
	}

	return nil
}

func ValidateServerPostRequest(r *http.Request) (*postServerRequest, error) {
	var requestData postServerRequest

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return nil, errors.New("failed to parse request body")
	}

	requestData.HostName = strings.ToLower(requestData.HostName)
	requestData.Name = strings.ToLower(requestData.Name)
	requestData.MapName = strings.ToLower(requestData.MapName)

	if err := validateName(requestData.Name, "name"); err != nil {
		return nil, err
	}

	if err := validateVersion(requestData.Version, "version", false); err != nil {
		return nil, err
	}

	if err := validateMemory(requestData.Memory); err != nil {
		return nil, err
	}

	if err := validateName(requestData.MapName, "mapName"); err != nil {
		return nil, err
	}

	if err := validateUrl(requestData.Url); err != nil {
		return nil, err
	}

	if err := validateHostName(requestData.HostName); err != nil {
		return nil, err
	}

	return &requestData, nil
}
