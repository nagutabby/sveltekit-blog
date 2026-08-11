package server_test

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"connectrpc.com/connect"

	contentv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/content/v1"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/content/v1/contentv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/server"
)

// TestContentServiceOverHTTP exercises the full stack for ContentService,
// the same way TestHealthServiceOverHTTP does for HealthService: a real
// HTTP server serving server.NewHandler(), hit by a generated Connect RPC
// client, backed by a fixture content directory on disk.
func TestContentServiceOverHTTP(t *testing.T) {
	dir := t.TempDir()
	articlesDir := filepath.Join(dir, "articles")
	if err := os.MkdirAll(articlesDir, 0o755); err != nil {
		t.Fatalf("failed to create fixture dir: %v", err)
	}
	articleMD := "---\ntitle: 統合テスト記事\nimage: images/test.png\npublishedAt: 2025-01-01\nupdatedAt: 2025-01-01\n---\n本文\n"
	if err := os.WriteFile(filepath.Join(articlesDir, "integration-test.md"), []byte(articleMD), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	cfg := server.Config{Content: content.NewLoader(dir)}
	srv := httptest.NewServer(server.NewHandler(cfg))
	defer srv.Close()

	client := contentv1connect.NewContentServiceClient(srv.Client(), srv.URL)

	resp, err := client.GetArticle(context.Background(), connect.NewRequest(&contentv1.GetArticleRequest{Id: "integration-test"}))
	if err != nil {
		t.Fatalf("GetArticle returned error: %v", err)
	}
	if got := resp.Msg.GetArticle().GetTitle(); got != "統合テスト記事" {
		t.Fatalf("Title = %q, want %q", got, "統合テスト記事")
	}
	if got := resp.Msg.GetArticle().GetImage(); got != "/content/articles/images/test.png" {
		t.Fatalf("Image = %q", got)
	}
}
