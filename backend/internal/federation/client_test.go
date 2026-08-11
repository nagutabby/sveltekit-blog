package federation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchActorSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/activity+json" {
			t.Errorf("Accept header = %q", got)
		}
		w.Header().Set("Content-Type", "application/activity+json")
		_, _ = w.Write([]byte(`{
			"inbox": "https://mastodon.example/users/alice/inbox",
			"publicKey": {
				"id": "https://mastodon.example/users/alice#main-key",
				"owner": "https://mastodon.example/users/alice",
				"publicKeyPem": "PEM-DATA"
			}
		}`))
	}))
	defer server.Close()

	actor, err := fetchActor(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatalf("fetchActor returned error: %v", err)
	}
	if actor.Inbox != "https://mastodon.example/users/alice/inbox" {
		t.Fatalf("Inbox = %q", actor.Inbox)
	}
	if actor.PublicKey.PublicKeyPem != "PEM-DATA" {
		t.Fatalf("PublicKeyPem = %q", actor.PublicKey.PublicKeyPem)
	}
}

func TestFetchActorMissingInbox(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"publicKey": {"publicKeyPem": "PEM-DATA"}}`))
	}))
	defer server.Close()

	_, err := fetchActor(context.Background(), server.Client(), server.URL)
	if err == nil {
		t.Fatal("expected an error for a missing inbox, got nil")
	}
}

func TestFetchActorNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchActor(context.Background(), server.Client(), server.URL)
	if err == nil {
		t.Fatal("expected an error for a 404 response, got nil")
	}
}

func TestFetchActorMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	_, err := fetchActor(context.Background(), server.Client(), server.URL)
	if err == nil {
		t.Fatal("expected an error for malformed JSON, got nil")
	}
}
