package watcher

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	common "mgarnier11.fr/go/orchestrator-common"
)

func TestParseComposeImages_UniqueAndFiltered(t *testing.T) {
	tempDir := t.TempDir()

	f1 := filepath.Join(tempDir, "a.yml")
	f2 := filepath.Join(tempDir, "b.yml")

	if err := os.WriteFile(f1, []byte(`services:
  one:
    image: nginx:${NGINX_VERSION}
  two:
    image: ghcr.io/test/app:${APP_VERSION}
`), 0644); err != nil {
		t.Fatalf("write f1: %v", err)
	}

	if err := os.WriteFile(f2, []byte(`services:
  one:
    image: nginx:${NGINX_VERSION}
  two:
    image: busybox:latest
`), 0644); err != nil {
		t.Fatalf("write f2: %v", err)
	}

	images := parseComposeImages([]*common.ComposeFile{
		{Path: f1},
		{Path: f2},
	})

	if len(images) != 2 {
		t.Fatalf("expected 2 unique variable-based images, got %d", len(images))
	}

	got := map[string]string{}
	for _, img := range images {
		got[img.Image] = img.VarName
	}

	if got["nginx"] != "NGINX_VERSION" {
		t.Fatalf("expected nginx -> NGINX_VERSION, got %#v", got)
	}
	if got["ghcr.io/test/app"] != "APP_VERSION" {
		t.Fatalf("expected ghcr.io/test/app -> APP_VERSION, got %#v", got)
	}
}

func TestReadVersionsEnv_ParsesExpectedEntries(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.env")

	content := "# comment\nINVALID\nA = 1.0.0\nB= v2.3.4\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("write versions.env: %v", err)
	}

	versions, err := readVersionsEnv(filePath)
	if err != nil {
		t.Fatalf("readVersionsEnv failed: %v", err)
	}

	if versions["A"] != "1.0.0" {
		t.Fatalf("expected A=1.0.0, got %q", versions["A"])
	}
	if versions["B"] != "v2.3.4" {
		t.Fatalf("expected B=v2.3.4, got %q", versions["B"])
	}
	if _, ok := versions["INVALID"]; ok {
		t.Fatalf("did not expect INVALID key")
	}
}

func TestWriteVersionsEnv_UpdatesOnlyKnownKeys(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "versions.env")

	initial := "# keep\nA=1.0.0\nB=2.0.0\n"
	if err := os.WriteFile(filePath, []byte(initial), 0644); err != nil {
		t.Fatalf("write versions.env: %v", err)
	}

	err := writeVersionsEnv(filePath, map[string]string{
		"A": "1.1.0",
		"C": "3.0.0",
	})
	if err != nil {
		t.Fatalf("writeVersionsEnv failed: %v", err)
	}

	updated, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read updated versions.env: %v", err)
	}

	got := string(updated)
	if !strings.Contains(got, "A=1.1.0") {
		t.Fatalf("expected A to be updated, got:\n%s", got)
	}
	if !strings.Contains(got, "B=2.0.0") {
		t.Fatalf("expected B to remain unchanged, got:\n%s", got)
	}
	if strings.Contains(got, "C=3.0.0") {
		t.Fatalf("did not expect C to be added, got:\n%s", got)
	}
}

func TestRun_UpdatesVersionsAndSendsNotification(t *testing.T) {
	tempDir := t.TempDir()
	composePath := filepath.Join(tempDir, "stack.yml")
	versionsPath := filepath.Join(tempDir, "versions.env")

	composeContent := "services:\n  s:\n    image: nginx:${NGINX_VERSION}\n"
	if err := os.WriteFile(composePath, []byte(composeContent), 0644); err != nil {
		t.Fatalf("write compose file: %v", err)
	}
	if err := os.WriteFile(versionsPath, []byte("NGINX_VERSION=1.24.0\n"), 0644); err != nil {
		t.Fatalf("write versions.env: %v", err)
	}

	oldGetComposeFiles := getComposeFiles
	oldGetLatest := getLatestCompatibleVersion
	oldSendNotification := sendNotification
	t.Cleanup(func() {
		getComposeFiles = oldGetComposeFiles
		getLatestCompatibleVersion = oldGetLatest
		sendNotification = oldSendNotification
	})

	getComposeFiles = func(composeDir string) ([]*common.ComposeFile, error) {
		if composeDir != tempDir {
			t.Fatalf("unexpected compose dir: %s", composeDir)
		}
		return []*common.ComposeFile{{Path: composePath}}, nil
	}
	getLatestCompatibleVersion = func(image, currentVersion string) (string, error) {
		if image == "nginx" && currentVersion == "1.24.0" {
			return "1.25.1", nil
		}
		return "", nil
	}

	notified := false
	sendNotification = func(title, body, tags string) error {
		notified = true
		if title != "Orchestrator - Version Updates" {
			t.Fatalf("unexpected title: %s", title)
		}
		if !strings.Contains(body, "NGINX_VERSION=1.25.1") {
			t.Fatalf("unexpected body: %s", body)
		}
		if tags != "package,arrow_up" {
			t.Fatalf("unexpected tags: %s", tags)
		}
		return nil
	}

	run(tempDir)

	updated, err := os.ReadFile(versionsPath)
	if err != nil {
		t.Fatalf("read updated versions.env: %v", err)
	}
	if !strings.Contains(string(updated), "NGINX_VERSION=1.25.1") {
		t.Fatalf("expected updated version, got:\n%s", string(updated))
	}
	if !notified {
		t.Fatalf("expected notification to be sent")
	}
}

func TestRun_NoUpdatesNoNotification(t *testing.T) {
	tempDir := t.TempDir()
	composePath := filepath.Join(tempDir, "stack.yml")
	versionsPath := filepath.Join(tempDir, "versions.env")

	if err := os.WriteFile(composePath, []byte("services:\n  s:\n    image: nginx:${NGINX_VERSION}\n"), 0644); err != nil {
		t.Fatalf("write compose file: %v", err)
	}
	if err := os.WriteFile(versionsPath, []byte("NGINX_VERSION=1.24.0\n"), 0644); err != nil {
		t.Fatalf("write versions.env: %v", err)
	}

	oldGetComposeFiles := getComposeFiles
	oldGetLatest := getLatestCompatibleVersion
	oldSendNotification := sendNotification
	t.Cleanup(func() {
		getComposeFiles = oldGetComposeFiles
		getLatestCompatibleVersion = oldGetLatest
		sendNotification = oldSendNotification
	})

	getComposeFiles = func(string) ([]*common.ComposeFile, error) {
		return []*common.ComposeFile{{Path: composePath}}, nil
	}
	getLatestCompatibleVersion = func(string, string) (string, error) {
		return "", nil
	}
	notified := false
	sendNotification = func(string, string, string) error {
		notified = true
		return errors.New("should not be called")
	}

	run(tempDir)

	updated, err := os.ReadFile(versionsPath)
	if err != nil {
		t.Fatalf("read versions.env: %v", err)
	}
	if strings.TrimSpace(string(updated)) != "NGINX_VERSION=1.24.0" {
		t.Fatalf("expected unchanged versions.env, got:\n%s", string(updated))
	}
	if notified {
		t.Fatalf("did not expect notification")
	}
}
