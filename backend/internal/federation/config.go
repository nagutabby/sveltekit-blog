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
	// ActorPublicKeyPEM is embedded verbatim in the actor document.
	ActorPublicKeyPEM string
	// ActorPrivateKeyPEM signs outgoing Accept activities.
	ActorPrivateKeyPEM string
	// HTTPClient is used for outbound federation requests. Defaults to
	// http.DefaultClient when nil.
	HTTPClient *http.Client
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
	// GetFollowerByActorID backs Inbox's handling of a remote actor
	// deleting itself: it looks up the existing row so UnfollowByActorID
	// can be called without clobbering inbox/publicKeyPem with blanks
	// (the actor is gone by the time its Delete arrives, so there's
	// nothing to re-fetch them from).
	GetFollowerByActorID(ctx context.Context, actorID string) (db.Follower, error)
	// ListActiveFollowerActorIDs backs paginated GET /actor/followers
	// pages (?page=N), returning real actorId items instead of just a
	// totalItems count.
	ListActiveFollowerActorIDs(ctx context.Context, arg db.ListActiveFollowerActorIDsParams) ([]string, error)
}

// RelayStore is the subset of *db.Queries the Inbox handler needs for
// relay Accept bookkeeping, and Following needs to list what this actor
// follows: this codebase never sends its own Follow to a relay (that's
// done out-of-band by whoever operates the instance), but it does record
// the Accept once a relay confirms it, so RelayConnection rows with
// connected=true are exactly this actor's "following" set.
type RelayStore interface {
	UpsertRelayConnectionAccepted(ctx context.Context, arg db.UpsertRelayConnectionAcceptedParams) (db.RelayConnection, error)
	ListRelayConnections(ctx context.Context) ([]db.RelayConnection, error)
}

// ArticleStore is the subset of *content.Loader the ArticleNote, NodeInfo,
// and Outbox handlers need: GetArticle to render a bare ActivityPub Note
// representation of an article (mirroring web's
// api/articles/[name]/+server.ts; other AP servers dereference an
// activity's object.id/url at this URL), ListArticles to report
// usage.localPosts in the NodeInfo document and to render the outbox's
// actual Create activities.
type ArticleStore interface {
	GetArticle(id string) (content.Article, error)
	ListArticles() ([]content.Article, error)
}
