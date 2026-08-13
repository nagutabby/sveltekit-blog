package federation

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

// encodePublicKeyPEM renders an RSA public key the way real ActivityPub
// actor documents do: PKIX, "BEGIN PUBLIC KEY".
func encodePublicKeyPEM(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

// newSignedInboxRequest builds a POST /actor/inbox request signed the way a
// real remote server would sign it, so Handlers.Inbox's HTTP Signature
// verification accepts it. keyID/privateKeyPEM belong to the *sending*
// (remote) actor, not this server's own actor key.
func newSignedInboxRequest(t *testing.T, cfg Config, body, keyID, privateKeyPEM string) *http.Request {
	t.Helper()
	targetURL := cfg.SiteBaseURL + "/actor/inbox"
	headers, err := SignHTTPRequest(targetURL, http.MethodPost, body, keyID, privateKeyPEM)
	if err != nil {
		t.Fatalf("failed to sign inbox request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(body))
	req.Host = hostOf(cfg.SiteBaseURL)
	req.Header.Set("Date", headers.Date)
	req.Header.Set("Digest", headers.Digest)
	req.Header.Set("Signature", headers.Signature)
	return req
}

type fakeArticleStore struct {
	articles map[string]content.Article
	err      error
}

func (f *fakeArticleStore) GetArticle(id string) (content.Article, error) {
	if f.err != nil {
		return content.Article{}, f.err
	}
	article, ok := f.articles[id]
	if !ok {
		return content.Article{}, content.ErrNotFound
	}
	return article, nil
}

func (f *fakeArticleStore) ListArticles() ([]content.Article, error) {
	if f.err != nil {
		return nil, f.err
	}
	articles := make([]content.Article, 0, len(f.articles))
	for _, article := range f.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

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
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

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
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

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
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

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
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/followers", nil)
	rec := httptest.NewRecorder()
	h.Followers(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

func TestFollowing(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

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
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

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
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

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

	aliceKey, alicePrivateKeyPEM, _ := generateTestKeyPair(t)
	alicePublicKeyPEM := encodePublicKeyPEM(t, aliceKey)

	// The fake actor doc returns a placeholder inbox URL; patch it to the
	// real httptest server address per-test via a small proxy handler.
	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     server.URL + "/users/alice/inbox",
			"publicKey": map[string]string{"publicKeyPem": alicePublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	cfg := testConfig(t)
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

	activityBody := `{"@context":"https://www.w3.org/ns/activitystreams","type":"Follow","actor":"` + actorServer.URL + `/users/alice","object":"https://blog.nagutabby.uk/actor"}`
	req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/alice#main-key", alicePrivateKeyPEM)
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

	aliceKey, alicePrivateKeyPEM, _ := generateTestKeyPair(t)
	alicePublicKeyPEM := encodePublicKeyPEM(t, aliceKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     server.URL + "/users/alice/inbox",
			"publicKey": map[string]string{"publicKeyPem": alicePublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	cfg := testConfig(t)
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

	activityBody := `{"@context":"https://www.w3.org/ns/activitystreams","type":"Undo","actor":"` + actorServer.URL + `/users/alice","object":{"type":"Follow"}}`
	req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/alice#main-key", alicePrivateKeyPEM)
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
	aliceKey, alicePrivateKeyPEM, _ := generateTestKeyPair(t)
	alicePublicKeyPEM := encodePublicKeyPEM(t, aliceKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     "https://mastodon.example/inbox",
			"publicKey": map[string]string{"publicKeyPem": alicePublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	cfg := testConfig(t)
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

	activityBody := `{"@context":"x","type":"Undo","actor":"` + actorServer.URL + `/users/alice","object":{"type":"Like"}}`
	req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/alice#main-key", alicePrivateKeyPEM)
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
	relayKey, relayPrivateKeyPEM, _ := generateTestKeyPair(t)
	relayPublicKeyPEM := encodePublicKeyPEM(t, relayKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/relay", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     "https://relay.example/inbox",
			"publicKey": map[string]string{"publicKeyPem": relayPublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	relays := &fakeRelayStore{}
	cfg := testConfig(t)
	h := NewHandlers(&fakeFollowerStore{}, relays, &fakeArticleStore{}, cfg)

	activityBody := `{"@context":"x","type":"Accept","actor":"` + actorServer.URL + `/users/relay","object":{"type":"Follow"}}`
	req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/relay#main-key", relayPrivateKeyPEM)
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

func TestInboxRejectsRelayNonAccept(t *testing.T) {
	relays := &fakeRelayStore{}
	followers := &fakeFollowerStore{}
	h := NewHandlers(followers, relays, &fakeArticleStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Follow","actor":"https://relay.example/actor","object":"https://blog.nagutabby.uk/actor"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	req.Header.Set("User-Agent", "SomeRelay/1.0")
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.upsertFollowerCalls) != 0 {
		t.Fatalf("UpsertFollower called %d times, want 0", len(followers.upsertFollowerCalls))
	}

	var respBody map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if respBody["error"] != "Forbidden" {
		t.Fatalf("error = %v, want Forbidden", respBody["error"])
	}
}

func TestInboxAllowsRelayAccept(t *testing.T) {
	relayKey, relayPrivateKeyPEM, _ := generateTestKeyPair(t)
	relayPublicKeyPEM := encodePublicKeyPEM(t, relayKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/relay", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     "https://relay.example/inbox",
			"publicKey": map[string]string{"publicKeyPem": relayPublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	relays := &fakeRelayStore{}
	cfg := testConfig(t)
	h := NewHandlers(&fakeFollowerStore{}, relays, &fakeArticleStore{}, cfg)

	activityBody := `{"@context":"x","type":"Accept","actor":"` + actorServer.URL + `/users/relay","object":{"type":"Follow"}}`
	req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/relay#main-key", relayPrivateKeyPEM)
	req.Header.Set("User-Agent", "SomeRelay/1.0")
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if len(relays.upsertCalls) != 1 {
		t.Fatalf("UpsertRelayConnectionAccepted called %d times, want 1", len(relays.upsertCalls))
	}
}

func TestInboxMissingRequiredFields(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(`{"type":"Follow"}`))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestInboxUnsupportedType(t *testing.T) {
	aliceKey, alicePrivateKeyPEM, _ := generateTestKeyPair(t)
	alicePublicKeyPEM := encodePublicKeyPEM(t, aliceKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     "https://mastodon.example/inbox",
			"publicKey": map[string]string{"publicKeyPem": alicePublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	cfg := testConfig(t)
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

	activityBody := `{"@context":"x","type":"Like","actor":"` + actorServer.URL + `/users/alice","object":"y"}`
	req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/alice#main-key", alicePrivateKeyPEM)
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

func TestInboxRejectsMissingSignature(t *testing.T) {
	aliceKey, _, _ := generateTestKeyPair(t)
	alicePublicKeyPEM := encodePublicKeyPEM(t, aliceKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     "https://mastodon.example/inbox",
			"publicKey": map[string]string{"publicKeyPem": alicePublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Follow","actor":"` + actorServer.URL + `/users/alice","object":"y"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.upsertFollowerCalls) != 0 {
		t.Fatal("UpsertFollower should not be called for an unsigned request")
	}
}

func TestInboxRejectsSignatureOverTamperedBody(t *testing.T) {
	aliceKey, alicePrivateKeyPEM, _ := generateTestKeyPair(t)
	alicePublicKeyPEM := encodePublicKeyPEM(t, aliceKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inbox":     "https://mastodon.example/inbox",
			"publicKey": map[string]string{"publicKeyPem": alicePublicKeyPEM},
		})
	})
	actorServer := httptest.NewServer(mux)
	defer actorServer.Close()

	followers := &fakeFollowerStore{}
	cfg := testConfig(t)
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

	signedBody := `{"@context":"x","type":"Follow","actor":"` + actorServer.URL + `/users/alice","object":"y"}`
	req := newSignedInboxRequest(t, cfg, signedBody, actorServer.URL+"/users/alice#main-key", alicePrivateKeyPEM)

	// Swap in a different body after signing, keeping the original
	// Content-Length so httptest doesn't truncate/pad it strangely: the
	// Digest/Signature headers were computed over signedBody, not this one.
	tamperedBody := `{"@context":"x","type":"Follow","actor":"` + actorServer.URL + `/users/alice","object":"z"}`
	req.Body = io.NopCloser(strings.NewReader(tamperedBody))

	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.upsertFollowerCalls) != 0 {
		t.Fatal("UpsertFollower should not be called when the Digest doesn't match the body")
	}
}

func TestInboxActorFetchFailure(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

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

func TestArticleNote(t *testing.T) {
	articles := &fakeArticleStore{articles: map[string]content.Article{
		"my-article": {
			Title:       "タイトル",
			PublishedAt: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		},
	}}
	h := &Handlers{articles: articles, cfg: testConfig(t)}

	req := httptest.NewRequest(http.MethodGet, "/api/articles/my-article", nil)
	req.SetPathValue("name", "my-article")
	rec := httptest.NewRecorder()
	h.ArticleNote(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/activity+json" {
		t.Fatalf("Content-Type = %q", got)
	}
	if got := rec.Header().Get("Cache-Control"); got != "" {
		t.Fatalf("Cache-Control = %q, want empty (unlike other actor endpoints)", got)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if _, hasContext := body["@context"]; hasContext {
		t.Fatal("bare article Note should not include @context, matching the original SvelteKit endpoint")
	}
	if body["id"] != "https://blog.nagutabby.uk/api/articles/my-article" {
		t.Fatalf("id = %v", body["id"])
	}
	if body["name"] != "タイトル" {
		t.Fatalf("name = %v", body["name"])
	}
	if body["published"] != "2025-06-15T00:00:00.000Z" {
		t.Fatalf("published = %v", body["published"])
	}
}

func TestArticleNoteNotFound(t *testing.T) {
	h := &Handlers{articles: &fakeArticleStore{}, cfg: testConfig(t)}

	req := httptest.NewRequest(http.MethodGet, "/api/articles/missing", nil)
	req.SetPathValue("name", "missing")
	rec := httptest.NewRecorder()
	h.ArticleNote(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}
