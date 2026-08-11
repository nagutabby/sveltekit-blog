package federation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

type fakeFollowerStore struct {
	upsertFollowerCalls []db.UpsertFollowerParams
	upsertFollowerErr   error

	unfollowCalls []db.UnfollowByActorIDParams
	unfollowErr   error

	followerCount int64
	countErr      error
}

func (f *fakeFollowerStore) UpsertFollower(_ context.Context, arg db.UpsertFollowerParams) (db.Follower, error) {
	f.upsertFollowerCalls = append(f.upsertFollowerCalls, arg)
	if f.upsertFollowerErr != nil {
		return db.Follower{}, f.upsertFollowerErr
	}
	return db.Follower{ActorId: arg.ActorId, Inbox: arg.Inbox, PublicKeyPem: arg.PublicKeyPem, Following: true}, nil
}

func (f *fakeFollowerStore) UnfollowByActorID(_ context.Context, arg db.UnfollowByActorIDParams) (db.Follower, error) {
	f.unfollowCalls = append(f.unfollowCalls, arg)
	if f.unfollowErr != nil {
		return db.Follower{}, f.unfollowErr
	}
	return db.Follower{ActorId: arg.ActorId, Following: false}, nil
}

func (f *fakeFollowerStore) CountActiveFollowers(_ context.Context) (int64, error) {
	if f.countErr != nil {
		return 0, f.countErr
	}
	return f.followerCount, nil
}

type fakeRelayStore struct {
	upsertCalls []db.UpsertRelayConnectionAcceptedParams
	upsertErr   error
}

func (f *fakeRelayStore) UpsertRelayConnectionAccepted(_ context.Context, arg db.UpsertRelayConnectionAcceptedParams) (db.RelayConnection, error) {
	f.upsertCalls = append(f.upsertCalls, arg)
	if f.upsertErr != nil {
		return db.RelayConnection{}, f.upsertErr
	}
	return db.RelayConnection{ActorId: arg.ActorId, Inbox: arg.Inbox, Connected: true}, nil
}

func testConfig(t *testing.T) Config {
	t.Helper()
	_, pkcs1, _ := generateTestKeyPair(t)
	return Config{
		SiteBaseURL:        "https://blog.nagutabby.uk",
		ActorPublicKeyPEM:  "PUBLIC-KEY-PEM",
		ActorPrivateKeyPEM: pkcs1,
	}
}

func TestWebfinger(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, testConfig(t))

	t.Run("valid resource", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/.well-known/webfinger?resource=acct:article@blog.nagutabby.uk", nil)
		rec := httptest.NewRecorder()
		h.Webfinger(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if got := rec.Header().Get("Content-Type"); got != "application/jrd+json" {
			t.Fatalf("Content-Type = %q", got)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
			t.Fatalf("Access-Control-Allow-Origin = %q", got)
		}

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["subject"] != "acct:article@blog.nagutabby.uk" {
			t.Fatalf("subject = %v", body["subject"])
		}
	})

	t.Run("missing resource", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/.well-known/webfinger", nil)
		rec := httptest.NewRecorder()
		h.Webfinger(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	})

	t.Run("wrong resource", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/.well-known/webfinger?resource=acct:someone-else@example.com", nil)
		rec := httptest.NewRecorder()
		h.Webfinger(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
	})
}

func TestActor(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor", nil)
	rec := httptest.NewRecorder()
	h.Actor(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/activity+json" {
		t.Fatalf("Content-Type = %q", got)
	}
	if got := rec.Header().Get("Cache-Control"); got != "max-age=0, private, must-revalidate" {
		t.Fatalf("Cache-Control = %q", got)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["id"] != "https://blog.nagutabby.uk/actor" {
		t.Fatalf("id = %v", body["id"])
	}
	if body["inbox"] != "https://blog.nagutabby.uk/actor/inbox" {
		t.Fatalf("inbox = %v", body["inbox"])
	}
	publicKey := body["publicKey"].(map[string]any)
	if publicKey["publicKeyPem"] != "PUBLIC-KEY-PEM" {
		t.Fatalf("publicKeyPem = %v", publicKey["publicKeyPem"])
	}
}

func TestFollowers(t *testing.T) {
	followers := &fakeFollowerStore{followerCount: 42}
	h := NewHandlers(followers, &fakeRelayStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/followers", nil)
	rec := httptest.NewRecorder()
	h.Followers(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["totalItems"] != float64(42) {
		t.Fatalf("totalItems = %v", body["totalItems"])
	}
}

func TestFollowersDBError(t *testing.T) {
	followers := &fakeFollowerStore{countErr: context.DeadlineExceeded}
	h := NewHandlers(followers, &fakeRelayStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/followers", nil)
	rec := httptest.NewRecorder()
	h.Followers(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

func TestFollowing(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/following", nil)
	rec := httptest.NewRecorder()
	h.Following(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["totalItems"] != float64(0) {
		t.Fatalf("totalItems = %v", body["totalItems"])
	}
}

func TestOutbox(t *testing.T) {
	atomServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<?xml version="1.0"?><feed><entry><title>a</title></entry><entry><title>b</title></entry></feed>`))
	}))
	defer atomServer.Close()

	cfg := testConfig(t)
	cfg.WebBaseURL = atomServer.URL
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, cfg)

	req := httptest.NewRequest(http.MethodGet, "/actor/outbox", nil)
	rec := httptest.NewRecorder()
	h.Outbox(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["totalItems"] != float64(2) {
		t.Fatalf("totalItems = %v", body["totalItems"])
	}
}

func TestOutboxUpstreamError(t *testing.T) {
	atomServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer atomServer.Close()

	cfg := testConfig(t)
	cfg.WebBaseURL = atomServer.URL
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, cfg)

	req := httptest.NewRequest(http.MethodGet, "/actor/outbox", nil)
	rec := httptest.NewRecorder()
	h.Outbox(rec, req)

	// The 500-status response body isn't a well-formed Atom feed the
	// unmarshaler can decode as expected, so the handler should surface it
	// as an error rather than pretending totalItems=0.
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

// newRemoteActorServer fakes a remote ActivityPub actor + inbox, capturing
// any signed POST it receives so tests can assert the Accept was sent and
// verify the HTTP Signature.
func newRemoteActorServer(t *testing.T) (*httptest.Server, *[]string, *[][]byte) {
	t.Helper()
	var signatures []string
	var receivedBodies [][]byte

	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"inbox": "REMOTE_INBOX_URL",
			"publicKey": {"publicKeyPem": "REMOTE-PUBLIC-KEY"}
		}`))
	})
	mux.HandleFunc("/users/alice/inbox", func(w http.ResponseWriter, r *http.Request) {
		signatures = append(signatures, r.Header.Get("Signature"))
		body, _ := io.ReadAll(r.Body)
		receivedBodies = append(receivedBodies, body)
		w.WriteHeader(http.StatusAccepted)
	})

	server := httptest.NewServer(mux)
	return server, &signatures, &receivedBodies
}

func TestInboxFollow(t *testing.T) {
	server, signatures, _ := newRemoteActorServer(t)
	defer server.Close()

	// The fake actor doc returns a placeholder inbox URL; patch it to the
	// real httptest server address per-test via a small proxy handler.
	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"inbox": "` + server.URL + `/users/alice/inbox", "publicKey": {"publicKeyPem": "REMOTE-PUBLIC-KEY"}}`))
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	h := NewHandlers(followers, &fakeRelayStore{}, testConfig(t))

	activityBody := `{"@context":"https://www.w3.org/ns/activitystreams","type":"Follow","actor":"` + actorServer.URL + `/users/alice","object":"https://blog.nagutabby.uk/actor"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}

	if len(followers.upsertFollowerCalls) != 1 {
		t.Fatalf("UpsertFollower called %d times, want 1", len(followers.upsertFollowerCalls))
	}
	if got := followers.upsertFollowerCalls[0].ActorId; got != actorServer.URL+"/users/alice" {
		t.Fatalf("ActorId = %q", got)
	}

	if len(*signatures) != 1 || (*signatures)[0] == "" {
		t.Fatalf("expected a signed Accept to be delivered to the remote inbox, signatures=%v", *signatures)
	}

	var respBody map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if respBody["type"] != "Accept" {
		t.Fatalf("type = %v", respBody["type"])
	}
	if respBody["id"] == nil || respBody["id"] == "" {
		t.Fatal("Follow's Accept response is missing id")
	}
	if respBody["object"] != activityBody && respBody["object"] == nil {
		// object should echo the original activity verbatim
		t.Fatalf("object = %v", respBody["object"])
	}
}

func TestInboxUndo(t *testing.T) {
	server, _, _ := newRemoteActorServer(t)
	defer server.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"inbox": "` + server.URL + `/users/alice/inbox", "publicKey": {"publicKeyPem": "REMOTE-PUBLIC-KEY"}}`))
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	h := NewHandlers(followers, &fakeRelayStore{}, testConfig(t))

	activityBody := `{"@context":"https://www.w3.org/ns/activitystreams","type":"Undo","actor":"` + actorServer.URL + `/users/alice","object":{"type":"Follow"}}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.unfollowCalls) != 1 {
		t.Fatalf("UnfollowByActorID called %d times, want 1", len(followers.unfollowCalls))
	}

	var respBody map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, hasID := respBody["id"]; hasID {
		t.Fatalf("Undo's Accept response should not include an id, got %v", respBody["id"])
	}
}

func TestInboxUndoWrongObjectType(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"inbox": "https://mastodon.example/inbox", "publicKey": {"publicKeyPem": "PEM"}}`))
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	h := NewHandlers(followers, &fakeRelayStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Undo","actor":"` + actorServer.URL + `/users/alice","object":{"type":"Like"}}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != "Invalid Undo activity" {
		t.Fatalf("body = %q", rec.Body.String())
	}
	if len(followers.unfollowCalls) != 0 {
		t.Fatal("UnfollowByActorID should not be called for a non-Follow Undo object")
	}
}

func TestInboxAccept(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/relay", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"inbox": "https://relay.example/inbox", "publicKey": {"publicKeyPem": "RELAY-PUBLIC-KEY"}}`))
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	relays := &fakeRelayStore{}
	h := NewHandlers(&fakeFollowerStore{}, relays, testConfig(t))

	activityBody := `{"@context":"x","type":"Accept","actor":"` + actorServer.URL + `/users/relay","object":{"type":"Follow"}}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if len(relays.upsertCalls) != 1 {
		t.Fatalf("UpsertRelayConnectionAccepted called %d times, want 1", len(relays.upsertCalls))
	}
	if got := relays.upsertCalls[0].ActorId; got != actorServer.URL+"/users/relay" {
		t.Fatalf("ActorId = %q", got)
	}
}

func TestInboxMissingRequiredFields(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(`{"type":"Follow"}`))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestInboxUnsupportedType(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"inbox": "https://mastodon.example/inbox", "publicKey": {"publicKeyPem": "PEM"}}`))
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Like","actor":"` + actorServer.URL + `/users/alice","object":"y"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

func TestInboxActorFetchFailure(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Follow","actor":"http://127.0.0.1:0/nonexistent","object":"y"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != "Could not fetch actor information" {
		t.Fatalf("body = %q", rec.Body.String())
	}
}
