package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federationadmin"
)

type noopArticleStore struct{}

func (noopArticleStore) GetArticle(id string) (content.Article, error) {
	return content.Article{}, nil
}

type noopRelayStore struct{}

func (noopRelayStore) ListRelayConnections(_ context.Context) ([]db.RelayConnection, error) {
	return nil, nil
}

func federationAdminPath(t *testing.T) string {
	t.Helper()
	return "/blog.federationadmin.v1.FederationAdminService/PublishArticleActivity"
}

func newTestFederationAdminHandler(token string) http.Handler {
	svc := federationadmin.NewService(noopArticleStore{}, noopRelayStore{}, federation.Config{})
	return NewHandler(Config{FederationAdmin: svc, FederationAdminToken: token})
}

func TestFederationAdminRejectsMissingAuthorization(t *testing.T) {
	handler := newTestFederationAdminHandler("secret-token")

	req := httptest.NewRequest(http.MethodPost, federationAdminPath(t), strings.NewReader(`{}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestFederationAdminRejectsWrongToken(t *testing.T) {
	handler := newTestFederationAdminHandler("secret-token")

	req := httptest.NewRequest(http.MethodPost, federationAdminPath(t), strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer wrong-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestFederationAdminRejectsEverythingWhenTokenUnset(t *testing.T) {
	handler := newTestFederationAdminHandler("")

	req := httptest.NewRequest(http.MethodPost, federationAdminPath(t), strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer ")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestFederationAdminAllowsCorrectToken(t *testing.T) {
	handler := newTestFederationAdminHandler("secret-token")

	req := httptest.NewRequest(http.MethodPost, federationAdminPath(t), strings.NewReader(`{"articleId":"x","changeType":"CHANGE_TYPE_DELETE"}`))
	req.Header.Set("Authorization", "Bearer secret-token")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusUnauthorized {
		t.Fatalf("status = %d, want anything but 401 (request should reach the handler), body=%s", rec.Code, rec.Body.String())
	}
}
