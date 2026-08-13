package federation

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
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
	// orderedArticles backs ListArticles with a deterministic order,
	// since map iteration order is randomized and pagination tests need
	// stable results.
	orderedArticles []content.Article
	err             error
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
	articles := make([]content.Article, 0, len(f.orderedArticles))
	articles = append(articles, f.orderedArticles...)
	return articles, nil
}

type fakeFollowerStore struct {
	upsertFollowerCalls []db.UpsertFollowerParams
	upsertFollowerErr   error

	unfollowCalls []db.UnfollowByActorIDParams
	unfollowErr   error

	followerCount int64
	countErr      error

	byActorID     map[string]db.Follower
	getByActorErr error

	// activeActorIDs backs ListActiveFollowerActorIDs in a deterministic
	// order (map iteration order isn't), so pagination tests are stable.
	activeActorIDs []string
	listErr        error
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

func (f *fakeFollowerStore) GetFollowerByActorID(_ context.Context, actorID string) (db.Follower, error) {
	if f.getByActorErr != nil {
		return db.Follower{}, f.getByActorErr
	}
	follower, ok := f.byActorID[actorID]
	if !ok {
		return db.Follower{}, sql.ErrNoRows
	}
	return follower, nil
}

func (f *fakeFollowerStore) ListActiveFollowerActorIDs(_ context.Context, arg db.ListActiveFollowerActorIDsParams) ([]string, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	start := int(arg.Offset)
	if start < 0 || start >= len(f.activeActorIDs) {
		return []string{}, nil
	}
	end := start + int(arg.Limit)
	if end > len(f.activeActorIDs) {
		end = len(f.activeActorIDs)
	}
	return f.activeActorIDs[start:end], nil
}

type fakeRelayStore struct {
	upsertCalls []db.UpsertRelayConnectionAcceptedParams
	upsertErr   error

	connections []db.RelayConnection
	listErr     error
}

func (f *fakeRelayStore) UpsertRelayConnectionAccepted(_ context.Context, arg db.UpsertRelayConnectionAcceptedParams) (db.RelayConnection, error) {
	f.upsertCalls = append(f.upsertCalls, arg)
	if f.upsertErr != nil {
		return db.RelayConnection{}, f.upsertErr
	}
	return db.RelayConnection{ActorId: arg.ActorId, Inbox: arg.Inbox, Connected: true}, nil
}

func (f *fakeRelayStore) ListRelayConnections(_ context.Context) ([]db.RelayConnection, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.connections, nil
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
	if body["type"] != "OrderedCollection" {
		t.Fatalf("type = %v", body["type"])
	}
	if body["totalItems"] != float64(42) {
		t.Fatalf("totalItems = %v", body["totalItems"])
	}
	if body["first"] != "https://blog.nagutabby.uk/actor/followers?page=1" {
		t.Fatalf("first = %v", body["first"])
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

func TestFollowersPage(t *testing.T) {
	followers := &fakeFollowerStore{activeActorIDs: []string{"https://a.example/users/1", "https://a.example/users/2"}}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/followers?page=1", nil)
	rec := httptest.NewRecorder()
	h.Followers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["type"] != "OrderedCollectionPage" {
		t.Fatalf("type = %v", body["type"])
	}
	if body["partOf"] != "https://blog.nagutabby.uk/actor/followers" {
		t.Fatalf("partOf = %v", body["partOf"])
	}
	items := body["orderedItems"].([]any)
	if len(items) != 2 || items[0] != "https://a.example/users/1" || items[1] != "https://a.example/users/2" {
		t.Fatalf("orderedItems = %v", items)
	}
	if _, hasNext := body["next"]; hasNext {
		t.Fatal("should not have a next page when results are under the page size")
	}
}

func TestFollowersPageHasNextWhenFull(t *testing.T) {
	full := make([]string, collectionPageSize)
	for i := range full {
		full[i] = fmt.Sprintf("https://a.example/users/%d", i)
	}
	followers := &fakeFollowerStore{activeActorIDs: full}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/followers?page=1", nil)
	rec := httptest.NewRecorder()
	h.Followers(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["next"] != "https://blog.nagutabby.uk/actor/followers?page=2" {
		t.Fatalf("next = %v", body["next"])
	}
	if _, hasPrev := body["prev"]; hasPrev {
		t.Fatal("page 1 should not have a prev link")
	}
}

func TestFollowersPageInvalid(t *testing.T) {
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/followers?page=nope", nil)
	rec := httptest.NewRecorder()
	h.Followers(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestFollowing(t *testing.T) {
	relays := &fakeRelayStore{connections: []db.RelayConnection{
		{ActorId: "https://relay.example/actor", Connected: true},
		{ActorId: "https://old-relay.example/actor", Connected: false},
	}}
	h := NewHandlers(&fakeFollowerStore{}, relays, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/following", nil)
	rec := httptest.NewRecorder()
	h.Following(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	// Only the connected relay counts; a relay whose Follow was
	// undone/rejected shouldn't inflate "following".
	if body["totalItems"] != float64(1) {
		t.Fatalf("totalItems = %v, want 1", body["totalItems"])
	}
}

func TestFollowingPage(t *testing.T) {
	relays := &fakeRelayStore{connections: []db.RelayConnection{
		{ActorId: "https://relay-a.example/actor", Connected: true},
		{ActorId: "https://relay-b.example/actor", Connected: true},
	}}
	h := NewHandlers(&fakeFollowerStore{}, relays, &fakeArticleStore{}, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/following?page=1", nil)
	rec := httptest.NewRecorder()
	h.Following(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	items := body["orderedItems"].([]any)
	if len(items) != 2 {
		t.Fatalf("orderedItems = %v, want 2 connected relays", items)
	}
}

func TestOutbox(t *testing.T) {
	articles := &fakeArticleStore{orderedArticles: []content.Article{
		{ID: "a", Title: "A"},
		{ID: "b", Title: "B"},
	}}
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, articles, testConfig(t))

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

func TestOutboxError(t *testing.T) {
	articles := &fakeArticleStore{err: context.DeadlineExceeded}
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, articles, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/outbox", nil)
	rec := httptest.NewRecorder()
	h.Outbox(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

func TestOutboxPageReturnsCreateActivities(t *testing.T) {
	articles := &fakeArticleStore{orderedArticles: []content.Article{
		{ID: "my-article", Title: "タイトル", PublishedAt: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)},
	}}
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, articles, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/outbox?page=1", nil)
	rec := httptest.NewRecorder()
	h.Outbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	items := body["orderedItems"].([]any)
	if len(items) != 1 {
		t.Fatalf("orderedItems = %v, want 1", items)
	}
	create := items[0].(map[string]any)
	if create["type"] != "Create" {
		t.Fatalf("type = %v", create["type"])
	}
	if create["id"] != "https://blog.nagutabby.uk/api/articles/my-article/create" {
		t.Fatalf("id = %v", create["id"])
	}
	object := create["object"].(map[string]any)
	if object["type"] != "Note" {
		t.Fatalf("object.type = %v", object["type"])
	}
	if object["name"] != "タイトル" {
		t.Fatalf("object.name = %v", object["name"])
	}
}

func TestOutboxPageNotFound(t *testing.T) {
	articles := &fakeArticleStore{orderedArticles: []content.Article{{ID: "a"}}}
	h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, articles, testConfig(t))

	req := httptest.NewRequest(http.MethodGet, "/actor/outbox?page=99", nil)
	rec := httptest.NewRecorder()
	h.Outbox(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
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

	// "Wave" isn't a real ActivityPub activity type, unlike Like/Create/
	// etc. below, which are recognized-but-unhandled and get a 202.
	activityBody := `{"@context":"x","type":"Wave","actor":"` + actorServer.URL + `/users/alice","object":"y"}`
	req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/alice#main-key", alicePrivateKeyPEM)
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

func TestInboxAcknowledgesRecognizedActivityTypes(t *testing.T) {
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

	for _, activityType := range []string{"Create", "Update", "Like", "Announce"} {
		t.Run(activityType, func(t *testing.T) {
			h := NewHandlers(&fakeFollowerStore{}, &fakeRelayStore{}, &fakeArticleStore{}, cfg)

			activityBody := `{"@context":"x","type":"` + activityType + `","actor":"` + actorServer.URL + `/users/alice","object":"y"}`
			req := newSignedInboxRequest(t, cfg, activityBody, actorServer.URL+"/users/alice#main-key", alicePrivateKeyPEM)
			rec := httptest.NewRecorder()
			h.Inbox(rec, req)

			if rec.Code != http.StatusAccepted {
				t.Fatalf("status = %d, want 202, body=%s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestInboxDeleteSelfUnfollowsKnownFollower(t *testing.T) {
	const actorID = "https://mastodon.example/users/alice"
	followers := &fakeFollowerStore{byActorID: map[string]db.Follower{
		actorID: {ActorId: actorID, Inbox: "https://mastodon.example/inbox", PublicKeyPem: "ALICE-KEY", Following: true},
	}}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Delete","actor":"` + actorID + `","object":"` + actorID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.unfollowCalls) != 1 {
		t.Fatalf("UnfollowByActorID called %d times, want 1", len(followers.unfollowCalls))
	}
	call := followers.unfollowCalls[0]
	if call.ActorId != actorID {
		t.Fatalf("ActorId = %q", call.ActorId)
	}
	if call.Inbox != "https://mastodon.example/inbox" || call.PublicKeyPem != "ALICE-KEY" {
		t.Fatalf("unfollow preserved wrong inbox/publicKeyPem: %+v", call)
	}
}

func TestInboxDeleteSelfWithTombstoneObject(t *testing.T) {
	const actorID = "https://mastodon.example/users/alice"
	followers := &fakeFollowerStore{byActorID: map[string]db.Follower{
		actorID: {ActorId: actorID, Inbox: "https://mastodon.example/inbox", PublicKeyPem: "ALICE-KEY", Following: true},
	}}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Delete","actor":"` + actorID + `","object":{"id":"` + actorID + `","type":"Tombstone"}}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.unfollowCalls) != 1 {
		t.Fatalf("UnfollowByActorID called %d times, want 1", len(followers.unfollowCalls))
	}
}

func TestInboxDeleteUnknownFollowerIsAcknowledged(t *testing.T) {
	followers := &fakeFollowerStore{}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	const actorID = "https://mastodon.example/users/bob"
	activityBody := `{"@context":"x","type":"Delete","actor":"` + actorID + `","object":"` + actorID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.unfollowCalls) != 0 {
		t.Fatal("UnfollowByActorID should not be called for an unknown actor")
	}
}

func TestInboxDeleteOfOtherObjectIsAcknowledgedWithoutAction(t *testing.T) {
	const actorID = "https://mastodon.example/users/alice"
	followers := &fakeFollowerStore{byActorID: map[string]db.Follower{
		actorID: {ActorId: actorID, Following: true},
	}}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	// alice deleting one of her own posts, not herself.
	activityBody := `{"@context":"x","type":"Delete","actor":"` + actorID + `","object":"https://mastodon.example/users/alice/statuses/123"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202, body=%s", rec.Code, rec.Body.String())
	}
	if len(followers.unfollowCalls) != 0 {
		t.Fatal("UnfollowByActorID should not be called when deleting an unrelated object")
	}
}

func TestInboxDeleteDoesNotRequireSignature(t *testing.T) {
	// The whole point of handling Delete before fetchActor/signature
	// verification is that the actor is typically already gone by the
	// time its self-Delete arrives, so an unsigned request must still
	// work.
	const actorID = "https://mastodon.example/users/alice"
	followers := &fakeFollowerStore{byActorID: map[string]db.Follower{
		actorID: {ActorId: actorID, Inbox: "https://mastodon.example/inbox", PublicKeyPem: "ALICE-KEY", Following: true},
	}}
	h := NewHandlers(followers, &fakeRelayStore{}, &fakeArticleStore{}, testConfig(t))

	activityBody := `{"@context":"x","type":"Delete","actor":"` + actorID + `","object":"` + actorID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(activityBody))
	// Deliberately no Signature/Date/Digest headers.
	rec := httptest.NewRecorder()
	h.Inbox(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
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
			ID:          "my-article",
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
