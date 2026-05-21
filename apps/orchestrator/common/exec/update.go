package exec

import (
	"context"
	"fmt"
	"os"

	"regexp"
	"strings"

	"github.com/joho/godotenv"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/ref"

	"mgarnier11.fr/go/libs/version"
	common "mgarnier11.fr/go/orchestrator-common"
)

var envVarRegexp = regexp.MustCompile(`\$\{([^}]+)\}`)

type ImageTag struct {
	ImageName string
	Tag       string
	Version   version.SemVer
}

func updateVersion(versionFilePath string, composeConfig *common.ComposeConfig, service string) error {
	/*
		The image tag can be either : mgarnier11/autoscaler:${AUTOSCALER_VERSION}
		or mgarnier11/autoscaler:1.3.0
		to retrieve the current version, we need to :
			parse the env variable name from the image tag
			retrieve the env variable value using the versions.env file
	*/

	if service == "" {
		for serviceName, serviceConfig := range composeConfig.Services {
			updateServiceVersion(versionFilePath, serviceName, serviceConfig)
		}
	} else {
		if composeConfig.Services[service] == nil {
			return fmt.Errorf("service %s not found in compose config %s %s %s", service, composeConfig.Host, composeConfig.Stack, composeConfig.Action)
		}

		return updateServiceVersion(versionFilePath, service, composeConfig.Services[service])
	}

	return nil

}

func updateServiceVersion(versionFilePath string, service string, serviceConfig *common.ComposeService) error {
	Logger.Debugf("Updating image %s version", serviceConfig.Image)

	parts := strings.Split(serviceConfig.Image, ":")
	if len(parts) < 2 {
		return fmt.Errorf("image tag is missing for service %s", service)
	}

	imageName := parts[0]
	imageTag := parts[1]

	if envVarMatch := envVarRegexp.FindStringSubmatch(imageTag); len(envVarMatch) == 2 {
		// If the image tag is an env variable
		envVarName := envVarMatch[1]

		semver, err := getEnvVarValue(envVarName, versionFilePath)
		if err != nil {
			return fmt.Errorf("error getting current version for env variable %s: %v", envVarName, err)
		}

		currentVersion, ok := version.ParseSemver(semver)
		if !ok {
			return fmt.Errorf("current version %s for env variable %s is not a valid semver", currentVersion.Raw, envVarName)
		}

		latestCompatibleVersion, mostRecentVersion, err := getNewImageVersions(imageName, currentVersion)
		if err != nil {
			return fmt.Errorf("error getting new image versions for %s: %v", imageName, err)
		}

		Logger.Infof("Service %s is at version %s, latest compatible version is %s, most recent version is %s", service, currentVersion.Raw, latestCompatibleVersion.Raw, mostRecentVersion.Raw)

		if latestCompatibleVersion.NewerThan(currentVersion) {
			comment := ""
			if latestCompatibleVersion.Raw != mostRecentVersion.Raw {
				comment = fmt.Sprintf("note: the most recent version is %s", mostRecentVersion.Raw)
			}
			err = updateEnvVarValue(envVarName, latestCompatibleVersion.Raw, comment, versionFilePath)
			if err != nil {
				return fmt.Errorf("error updating env variable %s: %v", envVarName, err)
			}

			Logger.Infof("Updated service %s to version %s", service, latestCompatibleVersion.Raw)
		} else {
			Logger.Infof("Service %s is already at the latest compatible version %s", service, currentVersion.Raw)
		}
	} else {
		return fmt.Errorf("image tag for service %s is not an env variable (%s), skipping update", service, imageTag)
	}

	return nil
}

func getEnvVarValue(envVarName, versionFilePath string) (string, error) {
	versions, err := godotenv.Read(versionFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading versions.env: %v", err)
	}

	value, ok := versions[envVarName]
	if !ok {
		return "", fmt.Errorf("env variable %s not found in versions.env", envVarName)
	}

	return value, nil
}

func getNewImageVersions(imageName string, currentVersion version.SemVer) (latestCompatibleVersion version.SemVer, mostRecentVersion version.SemVer, err error) {
	ctx := context.Background()

	rc := regclient.New()

	r, err := ref.New(imageName)
	if err != nil {
		Logger.Errorf("%v", err)
	}

	tags, err := rc.TagList(ctx, r)
	if err != nil {
		Logger.Errorf("%v", err)
	}

	latestCompatibleVersion = version.SemVer{Major: currentVersion.Major, Minor: currentVersion.Minor, Patch: currentVersion.Patch, Raw: currentVersion.Raw}
	mostRecentVersion = version.SemVer{Major: currentVersion.Major, Minor: currentVersion.Minor, Patch: currentVersion.Patch, Raw: currentVersion.Raw}

	for _, tag := range tags.Tags {
		semver, ok := version.ParseSemver(tag)
		if !ok {
			continue
		}

		if semver.NewerThan(latestCompatibleVersion) && semver.Major == currentVersion.Major {
			latestCompatibleVersion = semver
		}
		if semver.NewerThan(mostRecentVersion) {
			mostRecentVersion = semver
		}
	}

	return latestCompatibleVersion, mostRecentVersion, nil
}

func updateEnvVarValue(envVarName, newValue, commentValue, versionFilePath string) error {
	// search the env variable in the file and updates it in place
	// if comment value is defined, add it as a comment at the end of the line

	input, err := os.ReadFile(versionFilePath)
	if err != nil {
		return fmt.Errorf("error reading versions.env: %v", err)
	}

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, envVarName+"=") {
			newLine := fmt.Sprintf("%s=%s", envVarName, newValue)
			if commentValue != "" {
				newLine += fmt.Sprintf(" # %s", commentValue)
			}
			lines[i] = newLine
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(versionFilePath, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("error writing versions.env: %v", err)
	}

	return nil
}
