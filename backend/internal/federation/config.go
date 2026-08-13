package federation

import (
	"context"
	"net/http"
	"strings"

	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

// Config holds the settings the federation handlers need to reproduce the
// previous SvelteKit implementation's behavior.
type Config struct {
	// SiteBaseURL is the actor's public origin, e.g. "https://blog.nagutabby.uk".
	SiteBaseURL string
	// WebBaseURL is where /atom.xml (published by web) can be fetched from.
	// Defaults to SiteBaseURL when empty.
	WebBaseURL string
	// ActorPublicKeyPEM is embedded verbatim in the actor document.
	ActorPublicKeyPEM string
	// ActorPrivateKeyPEM signs outgoing Accept activities.
	ActorPrivateKeyPEM string
	// HTTPClient is used for outbound federation requests. Defaults to
	// http.DefaultClient when nil.
	HTTPClient *http.Client
}

func (c Config) webBaseURL() string {
	if c.WebBaseURL != "" {
		return c.WebBaseURL
	}
	return c.SiteBaseURL
}

func (c Config) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c Config) actorURL() string {
	return c.SiteBaseURL + "/actor"
}

func (c Config) actorKeyID() string {
	return c.actorURL() + "#main-key"
}

// NormalizePEM undoes the literal "\n" escaping commonly used when storing a
// multi-line PEM in a single-line environment variable.
func NormalizePEM(pem string) string {
	return strings.ReplaceAll(pem, `\n`, "\n")
}

// FollowerStore is the subset of *db.Queries the Inbox/Followers handlers
// need. Defined as an interface so tests can substitute a fake.
type FollowerStore interface {
	UpsertFollower(ctx context.Context, arg db.UpsertFollowerParams) (db.Follower, error)
	UnfollowByActorID(ctx context.Context, arg db.UnfollowByActorIDParams) (db.Follower, error)
	CountActiveFollowers(ctx context.Context) (int64, error)
}

// RelayStore is the subset of *db.Queries the Inbox handler needs for
// relay Accept bookkeeping.
type RelayStore interface {
	UpsertRelayConnectionAccepted(ctx context.Context, arg db.UpsertRelayConnectionAcceptedParams) (db.RelayConnection, error)
}

// ArticleStore is the subset of *content.Loader the ArticleNote and
// NodeInfo handlers need: GetArticle to render a bare ActivityPub Note
// representation of an article (mirroring web's
// api/articles/[name]/+server.ts; other AP servers dereference an
// activity's object.id/url at this URL), ListArticles to report
// usage.localPosts in the NodeInfo document.
type ArticleStore interface {
	GetArticle(id string) (content.Article, error)
	ListArticles() ([]content.Article, error)
}
