package validation

import (
	"fmt"
	"regexp"
)

func validateName(name string, key string) error {
	if name == "" {
		return fmt.Errorf("%s is required", key)
	}
	if len(name) > 15 {
		return fmt.Errorf("%s must be 15 characters or less", key)
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("%s must only contain letters, numbers, and hyphens", key)
	}

	return nil
}

func validateVersion(version string, key string, required bool) error {
	if version == "" {
		if !required {
			return nil
		}
		return fmt.Errorf("%s is required", key)
	}

	versionRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !versionRegex.MatchString(version) {
		return fmt.Errorf("%s must be in the format x.x.x", key)
	}

	return nil
}
