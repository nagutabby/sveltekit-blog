package d1_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db/d1"
)

// newTestClient points a d1.Client at a fake D1 HTTP query API that
// dispatches on the SQL statement's leading keyword, so each test can
// stub just the responses it needs.
func newTestClient(t *testing.T, handler func(w http.ResponseWriter, sql string, params []any)) *d1.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("Authorization header = %q, want %q", got, "Bearer test-token")
		}
		var body struct {
			SQL    string `json:"sql"`
			Params []any  `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		handler(w, body.SQL, body.Params)
	}))
	t.Cleanup(server.Close)

	return d1.New(d1.Config{
		AccountID:  "test-account",
		DatabaseID: "test-db",
		APIToken:   "test-token",
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	})
}

func writeD1Result(w http.ResponseWriter, results []map[string]any) {
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"errors":  []any{},
		"result": []map[string]any{
			{"success": true, "results": results},
		},
	})
}

func TestGetFollowerByActorID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, sql string, params []any) {
		if !strings.Contains(sql, `FROM "Follower"`) {
			t.Fatalf("unexpected SQL: %s", sql)
		}
		writeD1Result(w, []map[string]any{
			{
				"id":           float64(1),
				"actorId":      "https://mastodon.example/users/alice",
				"inbox":        "https://mastodon.example/users/alice/inbox",
				"publicKeyPem": "PEM",
				"following":    float64(1),
				"createdAt":    "2026-01-01T00:00:00Z",
				"updatedAt":    "2026-01-01T00:00:00Z",
			},
		})
	})

	got, err := client.GetFollowerByActorID(context.Background(), "https://mastodon.example/users/alice")
	if err != nil {
		t.Fatalf("GetFollowerByActorID returned error: %v", err)
	}
	if got.ID != 1 || !got.Following || got.PublicKeyPem != "PEM" {
		t.Fatalf("unexpected follower: %+v", got)
	}
}

func TestGetFollowerByActorIDNotFound(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, sql string, params []any) {
		writeD1Result(w, nil)
	})

	_, err := client.GetFollowerByActorID(context.Background(), "https://mastodon.example/users/nobody")
	if err != sql.ErrNoRows {
		t.Fatalf("err = %v, want sql.ErrNoRows", err)
	}
}

func TestUpsertRelayConnectionAcceptedNullLastAcceptedAt(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, sql string, params []any) {
		if len(params) < 3 || params[2] != nil {
			t.Fatalf("expected nil lastAcceptedAt param, got %#v", params)
		}
		writeD1Result(w, []map[string]any{
			{
				"id":             float64(1),
				"actorId":        "https://relay.example/actor",
				"inbox":          "https://relay.example/inbox",
				"connected":      float64(1),
				"lastAcceptedAt": nil,
				"createdAt":      "2026-01-01T00:00:00Z",
				"updatedAt":      "2026-01-01T00:00:00Z",
			},
		})
	})

	got, err := client.UpsertRelayConnectionAccepted(context.Background(), db.UpsertRelayConnectionAcceptedParams{
		ActorId:   "https://relay.example/actor",
		Inbox:     "https://relay.example/inbox",
		CreatedAt: "2026-01-01T00:00:00Z",
		UpdatedAt: "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("UpsertRelayConnectionAccepted returned error: %v", err)
	}
	if got.LastAcceptedAt.Valid {
		t.Fatalf("LastAcceptedAt.Valid = true, want false")
	}
}

func TestQueryErrorResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, sql string, params []any) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"errors": []map[string]any{
				{"code": 7500, "message": "syntax error"},
			},
			"result": []any{},
		})
	})

	_, err := client.GetFollowerByActorID(context.Background(), "https://mastodon.example/users/alice")
	if err == nil || !strings.Contains(err.Error(), "syntax error") {
		t.Fatalf("err = %v, want an error mentioning %q", err, "syntax error")
	}
}
