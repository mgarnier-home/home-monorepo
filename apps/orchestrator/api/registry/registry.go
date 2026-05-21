package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"mgarnier11.fr/go/libs/version"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// GetLatestCompatibleVersion returns the newest tag for the given image that
// has the same major version as currentVersionStr but a higher minor or patch.
// Returns "" if no update is found or if currentVersionStr is not semver.
func GetLatestCompatibleVersion(image string) (string, error) {
	current, ok := version.ParseSemver(currentVersionStr)
	if !ok {
		return "", nil // not semver, skip silently
	}

	tags, err := getTagsForImage(image)
	if err != nil {
		return "", err
	}

	hasVPrefix := strings.HasPrefix(currentVersionStr, "v")

	var best *version.SemVer
	for _, tag := range tags {
		sv, ok := version.ParseSemver(tag)
		if !ok {
			continue
		}
		if sv.Major != current.Major {
			continue
		}
		if !sv.NewerThan(current) {
			continue
		}
		if best == nil || sv.newerThan(*best) {
			copy := sv
			// Normalise the v-prefix to match the convention used in versions.env
			rawStripped := strings.TrimPrefix(copy.Raw, "v")
			if hasVPrefix {
				copy.Raw = "v" + rawStripped
			} else {
				copy.Raw = rawStripped
			}
			best = &copy
		}
	}

	if best == nil {
		return "", nil
	}
	return best.Raw, nil
}

func getTagsForImage(image string) ([]string, error) {
	if strings.HasPrefix(image, "ghcr.io/") {
		namespacedImage := strings.TrimPrefix(image, "ghcr.io/")
		return getGHCRTags(namespacedImage)
	}

	// Any other registry with a hostname (contains a dot before first slash) — skip
	slashIdx := strings.Index(image, "/")
	if slashIdx > 0 && strings.Contains(image[:slashIdx], ".") {
		return nil, fmt.Errorf("unsupported registry")
	}

	return getDockerHubTags(image)
}

// --- Docker Hub ---

type dockerHubTagsResponse struct {
	Next    *string `json:"next"`
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func getDockerHubTags(image string) ([]string, error) {
	namespace := "library"
	name := image
	if strings.Contains(image, "/") {
		parts := strings.SplitN(image, "/", 2)
		namespace = parts[0]
		name = parts[1]
	}

	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags?page_size=100&ordering=-name", namespace, name)
	var tags []string

	for url != "" {
		resp, err := httpClient.Get(url)
		if err != nil {
			return nil, fmt.Errorf("fetching Docker Hub tags: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Docker Hub API status %d for %s/%s", resp.StatusCode, namespace, name)
		}

		var result dockerHubTagsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decoding Docker Hub response: %w", err)
		}

		for _, r := range result.Results {
			tags = append(tags, r.Name)
		}

		if result.Next != nil && *result.Next != "" {
			url = *result.Next
		} else {
			url = ""
		}
	}

	return tags, nil
}

// --- GitHub Container Registry ---

func getGHCRTags(namespacedImage string) ([]string, error) {
	tokenURL := fmt.Sprintf("https://ghcr.io/token?service=ghcr.io&scope=repository:%s:pull", namespacedImage)
	resp, err := httpClient.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("getting GHCR token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decoding GHCR token: %w", err)
	}

	tagsURL := fmt.Sprintf("https://ghcr.io/v2/%s/tags/list", namespacedImage)
	req, err := http.NewRequest("GET", tagsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)

	resp2, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching GHCR tags: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GHCR tags API status %d for %s", resp2.StatusCode, namespacedImage)
	}

	var tagsResp struct {
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("decoding GHCR tags: %w", err)
	}

	return tagsResp.Tags, nil
}
