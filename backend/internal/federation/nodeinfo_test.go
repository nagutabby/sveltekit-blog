package federation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
)

func TestWellKnownNodeInfo(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/.well-known/nodeinfo", nil)
	rec := httptest.NewRecorder()
	h.WellKnownNodeInfo(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body struct {
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if len(body.Links) != 2 {
		t.Fatalf("links = %v, want 2 entries", body.Links)
	}
	if body.Links[0].Href != "https://blog.nagutabby.uk/nodeinfo/2.1" {
		t.Fatalf("2.1 href = %q", body.Links[0].Href)
	}
	if body.Links[1].Href != "https://blog.nagutabby.uk/nodeinfo/2.0" {
		t.Fatalf("2.0 href = %q", body.Links[1].Href)
	}
}

func TestNodeInfo21(t *testing.T) {
	articles := &fakeArticleStore{articles: map[string]content.Article{
		"a": {Title: "a"},
		"b": {Title: "b"},
	}}
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, articles, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/nodeinfo/2.1", nil)
	rec := httptest.NewRecorder()
	h.NodeInfo21(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != nodeInfo21ContentType {
		t.Fatalf("Content-Type = %q, want %q", got, nodeInfo21ContentType)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["version"] != "2.1" {
		t.Fatalf("version = %v", body["version"])
	}
	software := body["software"].(map[string]any)
	if software["name"] != softwareName {
		t.Fatalf("software.name = %v", software["name"])
	}
	if software["repository"] != softwareRepository {
		t.Fatalf("2.1 software.repository = %v, want %q", software["repository"], softwareRepository)
	}
	protocols := body["protocols"].([]any)
	if len(protocols) != 1 || protocols[0] != "activitypub" {
		t.Fatalf("protocols = %v", protocols)
	}
	usage := body["usage"].(map[string]any)
	if usage["localPosts"] != float64(2) {
		t.Fatalf("usage.localPosts = %v, want 2", usage["localPosts"])
	}
}

func TestNodeInfo20OmitsRepository(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/nodeinfo/2.0", nil)
	rec := httptest.NewRecorder()
	h.NodeInfo20(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != nodeInfo20ContentType {
		t.Fatalf("Content-Type = %q, want %q", got, nodeInfo20ContentType)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["version"] != "2.0" {
		t.Fatalf("version = %v", body["version"])
	}
	software := body["software"].(map[string]any)
	if _, hasRepository := software["repository"]; hasRepository {
		t.Fatal("2.0 software should not include a repository field")
	}
}
