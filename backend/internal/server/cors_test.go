package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORSHandlesPreflight(t *testing.T) {
	called := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
	handler := withCORS("https://blog.nagutabby.uk", inner)

	req := httptest.NewRequest(http.MethodOptions, "/blog.contact.v1.ContactService/SubmitContact", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://blog.nagutabby.uk" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "POST" {
		t.Fatalf("Access-Control-Allow-Methods = %q", got)
	}
	if called {
		t.Fatal("preflight request should not reach the inner handler")
	}
}

func TestWithCORSPassesThroughNonPreflightRequests(t *testing.T) {
	called := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler := withCORS("https://blog.nagutabby.uk", inner)

	req := httptest.NewRequest(http.MethodPost, "/blog.contact.v1.ContactService/SubmitContact", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected the inner handler to be called for a non-preflight request")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://blog.nagutabby.uk" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}
