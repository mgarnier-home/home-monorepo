package watcher

import (
	"bufio"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"mgarnier11.fr/go/libs/logger"
	ntfy "mgarnier11.fr/go/libs/ntfy"
	common "mgarnier11.fr/go/orchestrator-common"
	compose_files "mgarnier11.fr/go/orchestrator-common/files"
)

var Logger = logger.NewLogger("[WATCHER]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF88")), nil)


// matches: image: some/image:${VAR_NAME}
var imageRegexp = regexp.MustCompile(`(?m)^\s*image:\s+([^:\s]+):\$\{([A-Z0-9_]+)\}`)

type imageInfo struct {
	Image   string
	VarName string
}

// Start launches the version watcher in a background goroutine.
// It runs immediately on startup, then every hour.
func Start(composeDirPath string) {
	go func() {
		run(composeDirPath)

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			run(composeDirPath)
		}
	}()
}

func run(composeDirPath string) {
	Logger.Infof("Running version check")

	composeFiles, err := compose_files.GetComposeFiles(composeDirPath)
	if err != nil {
		Logger.Errorf("Error getting compose files: %v", err)
		return
	}

	images := parseComposeImages(composeFiles)
	Logger.Infof("Found %d unique image references to check", len(images))

	versionsEnvPath := path.Join(composeDirPath, "versions.env")
	versions, err := readVersionsEnv(versionsEnvPath)
	if err != nil {
		Logger.Errorf("Error reading versions.env: %v", err)
		return
	}

	updates := map[string]string{}

	for _, img := range images {
		currentVersion, ok := versions[img.VarName]
		if !ok {
			Logger.Debugf("No version entry for %s", img.VarName)
			continue
		}

		Logger.Debugf("Checking %s (%s = %s)", img.Image, img.VarName, currentVersion)

		newVersion, err := GetLatestCompatibleVersion(img.Image, currentVersion)
		if err != nil {
			Logger.Debugf("Skipping %s: %v", img.Image, err)
			continue
		}

		if newVersion != "" {
			Logger.Infof("Update available for %s: %s → %s", img.Image, currentVersion, newVersion)
			updates[img.VarName] = newVersion
		}
	}

	if len(updates) == 0 {
		Logger.Infof("No updates found")
		return
	}

	for varName, newVersion := range updates {
		versions[varName] = newVersion
	}

	if err := writeVersionsEnv(versionsEnvPath, versions); err != nil {
		Logger.Errorf("Error writing versions.env: %v", err)
		return
	}

	var sb strings.Builder
	for varName, newVersion := range updates {
		sb.WriteString(varName + "=" + newVersion + "\n")
	}

	if err := sendNotification("Orchestrator - Version Updates", sb.String(), "package,arrow_up"); err != nil {
		Logger.Errorf("Error sending ntfy notification: %v", err)
	}

	Logger.Infof("Updated %d version(s) in versions.env", len(updates))
}

func parseComposeImages(files []*common.ComposeFile) []imageInfo {
	seen := map[string]bool{}
	images := []imageInfo{}

	for _, f := range files {
		content, err := os.ReadFile(f.Path)
		if err != nil {
			Logger.Debugf("Error reading compose file %s: %v", f.Path, err)
			continue
		}

		matches := imageRegexp.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			image := m[1]
			varName := m[2]
			key := image + ":" + varName
			if seen[key] {
				continue
			}
			seen[key] = true
			images = append(images, imageInfo{Image: image, VarName: varName})
		}
	}

	return images
}

func readVersionsEnv(filePath string) (map[string]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	versions := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			versions[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return versions, scanner.Err()
}

func writeVersionsEnv(filePath string, versions map[string]string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if newVal, ok := versions[key]; ok {
				lines[i] = key + "=" + newVal
			}
		}
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}
