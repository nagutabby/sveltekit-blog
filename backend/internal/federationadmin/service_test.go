package federationadmin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"

	federationadminv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/federationadmin/v1"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
)

type fakeArticleStore struct {
	articles map[string]content.Article
}

func (f *fakeArticleStore) GetArticle(id string) (content.Article, error) {
	article, ok := f.articles[id]
	if !ok {
		return content.Article{}, content.ErrNotFound
	}
	return article, nil
}

type fakeRelayStore struct {
	relays []db.RelayConnection
	err    error
}

func (f *fakeRelayStore) ListRelayConnections(_ context.Context) ([]db.RelayConnection, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.relays, nil
}

func testFederationConfig(t *testing.T) federation.Config {
	t.Helper()
	_, privateKeyPEM := generateTestKeyPair(t)
	return federation.Config{
		SiteBaseURL:        "https://blog.nagutabby.uk",
		ActorPrivateKeyPEM: privateKeyPEM,
	}
}

func TestPublishArticleActivityCreateDeliversToAllRelays(t *testing.T) {
	var receivedSignatures []string
	var receivedBodies []map[string]any
	relayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedSignatures = append(receivedSignatures, r.Header.Get("Signature"))
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		receivedBodies = append(receivedBodies, body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer relayServer.Close()

	articles := &fakeArticleStore{articles: map[string]content.Article{
		"my-article": {Title: "タイトル", PublishedAt: time.Now()},
	}}
	relays := &fakeRelayStore{relays: []db.RelayConnection{
		{ActorId: "relay-1", Inbox: relayServer.URL + "/inbox-1"},
		{ActorId: "relay-2", Inbox: relayServer.URL + "/inbox-2"},
	}}

	svc := NewService(articles, relays, testFederationConfig(t))

	resp, err := svc.PublishArticleActivity(context.Background(), connect.NewRequest(&federationadminv1.PublishArticleActivityRequest{
		ArticleId:  "my-article",
		ChangeType: federationadminv1.ChangeType_CHANGE_TYPE_CREATE,
	}))
	if err != nil {
		t.Fatalf("PublishArticleActivity returned error: %v", err)
	}
	_ = resp

	if len(receivedSignatures) != 2 {
		t.Fatalf("relay deliveries = %d, want 2", len(receivedSignatures))
	}
	for i, sig := range receivedSignatures {
		if sig == "" {
			t.Fatalf("delivery %d had no Signature header", i)
		}
	}
	for i, body := range receivedBodies {
		if body["type"] != "Create" {
			t.Fatalf("delivery %d type = %v, want Create", i, body["type"])
		}
		object := body["object"].(map[string]any)
		if object["name"] != "タイトル" {
			t.Fatalf("delivery %d object.name = %v", i, object["name"])
		}
	}
}

func TestPublishArticleActivityDeleteDoesNotLookUpArticle(t *testing.T) {
	relayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer relayServer.Close()

	// No articles registered: a Delete must not need to look one up,
	// since by the time it fires the Markdown source is typically gone.
	articles := &fakeArticleStore{articles: map[string]content.Article{}}
	relays := &fakeRelayStore{relays: []db.RelayConnection{{ActorId: "relay-1", Inbox: relayServer.URL}}}

	svc := NewService(articles, relays, testFederationConfig(t))

	_, err := svc.PublishArticleActivity(context.Background(), connect.NewRequest(&federationadminv1.PublishArticleActivityRequest{
		ArticleId:  "already-deleted-article",
		ChangeType: federationadminv1.ChangeType_CHANGE_TYPE_DELETE,
	}))
	if err != nil {
		t.Fatalf("PublishArticleActivity returned error: %v", err)
	}
}

func TestPublishArticleActivityCreateArticleNotFound(t *testing.T) {
	articles := &fakeArticleStore{articles: map[string]content.Article{}}
	relays := &fakeRelayStore{}

	svc := NewService(articles, relays, testFederationConfig(t))

	_, err := svc.PublishArticleActivity(context.Background(), connect.NewRequest(&federationadminv1.PublishArticleActivityRequest{
		ArticleId:  "missing",
		ChangeType: federationadminv1.ChangeType_CHANGE_TYPE_CREATE,
	}))
	if err == nil {
		t.Fatal("expected an error for a missing article, got nil")
	}
}

func TestPublishArticleActivityNoRelaysIsANoop(t *testing.T) {
	articles := &fakeArticleStore{articles: map[string]content.Article{
		"my-article": {Title: "タイトル", PublishedAt: time.Now()},
	}}
	relays := &fakeRelayStore{relays: nil}

	svc := NewService(articles, relays, testFederationConfig(t))

	_, err := svc.PublishArticleActivity(context.Background(), connect.NewRequest(&federationadminv1.PublishArticleActivityRequest{
		ArticleId:  "my-article",
		ChangeType: federationadminv1.ChangeType_CHANGE_TYPE_UPDATE,
	}))
	if err != nil {
		t.Fatalf("PublishArticleActivity returned error: %v", err)
	}
}

func TestPublishArticleActivityContinuesAfterOneRelayFails(t *testing.T) {
	var deliveredToWorking bool

	failingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failingServer.Close()

	workingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deliveredToWorking = true
		w.WriteHeader(http.StatusAccepted)
	}))
	defer workingServer.Close()

	articles := &fakeArticleStore{articles: map[string]content.Article{
		"my-article": {Title: "タイトル", PublishedAt: time.Now()},
	}}
	relays := &fakeRelayStore{relays: []db.RelayConnection{
		{ActorId: "broken-relay", Inbox: failingServer.URL},
		{ActorId: "working-relay", Inbox: workingServer.URL},
	}}

	svc := NewService(articles, relays, testFederationConfig(t))

	_, err := svc.PublishArticleActivity(context.Background(), connect.NewRequest(&federationadminv1.PublishArticleActivityRequest{
		ArticleId:  "my-article",
		ChangeType: federationadminv1.ChangeType_CHANGE_TYPE_UPDATE,
	}))
	if err != nil {
		t.Fatalf("PublishArticleActivity returned error: %v", err)
	}
	if !deliveredToWorking {
		t.Fatal("expected delivery to the working relay to still happen after the broken one failed")
	}
}
