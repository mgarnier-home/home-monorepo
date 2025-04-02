package validation

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"mgarnier11.fr/go/mineager/server/objects/dto"
)

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

func ValidateServerPostRequest(r *http.Request) (*dto.CreateServerRequestDto, error) {
	var requestData dto.CreateServerRequestDto

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return nil, errors.New("failed to parse request body")
	}

	requestData.Name = strings.ToLower(requestData.Name)
	requestData.MapName = strings.ToLower(requestData.MapName)

	if err := validateName(requestData.Name, "name"); err != nil {
		return nil, err
	}

	if err := validateVersion(requestData.Version, "version", true); err != nil {
		return nil, err
	}

	if err := validateMemory(requestData.Memory); err != nil {
		return nil, err
	}

	if err := validateName(requestData.MapName, "mapName"); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateServerDeleteRequest(r *http.Request) (*dto.DeleteServerRequestDto, error) {
	var requestData dto.DeleteServerRequestDto

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return nil, errors.New("failed to parse request body")
	}

	return &requestData, nil
}
