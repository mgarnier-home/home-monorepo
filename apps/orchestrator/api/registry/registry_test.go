package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type rewriteToServerTransport struct {
	target *url.URL
	base   http.RoundTripper
}

func (t *rewriteToServerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.URL.Scheme = t.target.Scheme
	clone.URL.Host = t.target.Host
	return t.base.RoundTrip(clone)
}

func withMockRegistryServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}

	originalClient := httpClient
	t.Cleanup(func() {
		httpClient = originalClient
	})

	transport := &rewriteToServerTransport{
		target: targetURL,
		base:   server.Client().Transport,
	}

	httpClient = &http.Client{Transport: transport}
}

func TestParseSemver(t *testing.T) {
	tests := []struct {
		in      string
		ok      bool
		expects semVer
	}{
		{in: "1.2.3", ok: true, expects: semVer{Major: 1, Minor: 2, Patch: 3, Raw: "1.2.3"}},
		{in: "v10.20.30", ok: true, expects: semVer{Major: 10, Minor: 20, Patch: 30, Raw: "v10.20.30"}},
		{in: "1.2", ok: false},
		{in: "latest", ok: false},
	}

	for _, tc := range tests {
		got, ok := parseSemver(tc.in)
		if ok != tc.ok {
			t.Fatalf("parseSemver(%q) ok=%v, want %v", tc.in, ok, tc.ok)
		}
		if !tc.ok {
			continue
		}
		if got != tc.expects {
			t.Fatalf("parseSemver(%q)=%+v, want %+v", tc.in, got, tc.expects)
		}
	}
}

func TestSemVerNewerThan(t *testing.T) {
	if !((semVer{Minor: 3, Patch: 0}).newerThan(semVer{Minor: 2, Patch: 99})) {
		t.Fatalf("expected higher minor to be newer")
	}
	if !((semVer{Minor: 2, Patch: 5}).newerThan(semVer{Minor: 2, Patch: 4})) {
		t.Fatalf("expected higher patch to be newer")
	}
	if (semVer{Minor: 2, Patch: 4}).newerThan(semVer{Minor: 2, Patch: 4}) {
		t.Fatalf("expected equal versions to not be newer")
	}
}

func TestGetDockerHubTags_Paginates(t *testing.T) {
	withMockRegistryServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v2/repositories/library/nginx/tags") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.Query().Get("page") == "2" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"next":    "",
				"results": []map[string]string{{"name": "1.25.1"}},
			})
			return
		}

		next := "https://hub.docker.com/v2/repositories/library/nginx/tags?page=2"
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"next":    next,
			"results": []map[string]string{{"name": "1.24.0"}},
		})
	})

	tags, err := getDockerHubTags("nginx")
	if err != nil {
		t.Fatalf("getDockerHubTags failed: %v", err)
	}
	if len(tags) != 2 || tags[0] != "1.24.0" || tags[1] != "1.25.1" {
		t.Fatalf("unexpected tags: %#v", tags)
	}
}

func TestGetGHCRTags_UsesTokenAndAuth(t *testing.T) {
	withMockRegistryServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/token":
			if r.URL.Query().Get("scope") != "repository:owner/app:pull" {
				t.Fatalf("unexpected token scope: %s", r.URL.Query().Get("scope"))
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"token": "abc123"})
		case r.URL.Path == "/v2/owner/app/tags/list":
			if got := r.Header.Get("Authorization"); got != "Bearer abc123" {
				t.Fatalf("expected bearer auth, got %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"tags": []string{"1.0.0", "1.1.0"}})
		default:
			t.Fatalf("unexpected GHCR path: %s", r.URL.Path)
		}
	})

	tags, err := getGHCRTags("owner/app")
	if err != nil {
		t.Fatalf("getGHCRTags failed: %v", err)
	}
	if len(tags) != 2 || tags[0] != "1.0.0" || tags[1] != "1.1.0" {
		t.Fatalf("unexpected tags: %#v", tags)
	}
}

func TestGetTagsForImage_UnsupportedRegistry(t *testing.T) {
	_, err := getTagsForImage("quay.io/example/app")
	if err == nil {
		t.Fatalf("expected unsupported registry error")
	}
}

func TestGetLatestCompatibleVersion_PrefixAndMajorRules(t *testing.T) {
	withMockRegistryServer(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v2/repositories/library/nginx/tags") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"next": "",
			"results": []map[string]string{
				{"name": "1.2.4"},
				{"name": "v1.3.0"},
				{"name": "2.0.0"},
				{"name": "latest"},
			},
		})
	})

	got, err := GetLatestCompatibleVersion("nginx", "1.2.3")
	if err != nil {
		t.Fatalf("GetLatestCompatibleVersion failed: %v", err)
	}
	if got != "1.3.0" {
		t.Fatalf("expected 1.3.0, got %q", got)
	}

	got, err = GetLatestCompatibleVersion("nginx", "v1.2.3")
	if err != nil {
		t.Fatalf("GetLatestCompatibleVersion failed: %v", err)
	}
	if got != "v1.3.0" {
		t.Fatalf("expected v1.3.0, got %q", got)
	}

	got, err = GetLatestCompatibleVersion("nginx", "latest")
	if err != nil {
		t.Fatalf("GetLatestCompatibleVersion for non-semver should not fail: %v", err)
	}
	if got != "" {
		t.Fatalf("expected no update for non-semver current version, got %q", got)
	}
}

func TestGetTagsForImage_RoutesGHCR(t *testing.T) {
	withMockRegistryServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/token":
			_ = json.NewEncoder(w).Encode(map[string]string{"token": "token"})
		case r.URL.Path == "/v2/org/service/tags/list":
			if got := r.Header.Get("Authorization"); got != "Bearer token" {
				t.Fatalf("expected auth header, got %q", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"tags": []string{"0.1.0"}})
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, "not found")
		}
	})

	tags, err := getTagsForImage("ghcr.io/org/service")
	if err != nil {
		t.Fatalf("getTagsForImage failed: %v", err)
	}
	if len(tags) != 1 || tags[0] != "0.1.0" {
		t.Fatalf("unexpected tags: %#v", tags)
	}
}
