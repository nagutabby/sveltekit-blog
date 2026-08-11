package federationadmin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"

	federationadminv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/federationadmin/v1"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
)

var errUnspecifiedChangeType = errors.New("federationadmin: change_type is required")

// ArticleStore is the subset of *content.Loader PublishArticleActivity
// needs to look up an article's title/publishedAt for Create/Update
// activities. Delete activities don't need it, since by the time a
// Delete fires the article's Markdown source is typically already gone.
type ArticleStore interface {
	GetArticle(id string) (content.Article, error)
}

// RelayStore is the subset of *db.Queries PublishArticleActivity needs to
// broadcast to every connected relay.
type RelayStore interface {
	ListRelayConnections(ctx context.Context) ([]db.RelayConnection, error)
}

// Service implements the blog.federationadmin.v1.FederationAdminService
// Connect RPC service. It ports web's
// api/articles/[name]/{create,update,delete}/+server.ts (activity
// construction + LD-Signature) and api/activitypub/sender/+server.ts
// (broadcast to relays with an HTTP Signature).
type Service struct {
	articles   ArticleStore
	relays     RelayStore
	cfg        federation.Config
	httpClient *http.Client
}

func NewService(articles ArticleStore, relays RelayStore, cfg federation.Config) *Service {
	client := cfg.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	return &Service{articles: articles, relays: relays, cfg: cfg, httpClient: client}
}

func (s *Service) PublishArticleActivity(
	ctx context.Context,
	req *connect.Request[federationadminv1.PublishArticleActivityRequest],
) (*connect.Response[federationadminv1.PublishArticleActivityResponse], error) {
	articleID := req.Msg.GetArticleId()

	act, err := s.buildActivity(articleID, req.Msg.GetChangeType())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	relays, err := s.relays.ListRelayConnections(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if len(relays) == 0 {
		return connect.NewResponse(&federationadminv1.PublishArticleActivityResponse{}), nil
	}

	body, err := json.Marshal(act)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, relay := range relays {
		if err := s.deliver(ctx, relay.Inbox, body); err != nil {
			log.Printf("federationadmin: delivering activity to %s failed: %v", relay.Inbox, err)
			continue
		}
	}

	return connect.NewResponse(&federationadminv1.PublishArticleActivityResponse{}), nil
}

func (s *Service) buildActivity(articleID string, changeType federationadminv1.ChangeType) (activity, error) {
	switch changeType {
	case federationadminv1.ChangeType_CHANGE_TYPE_CREATE, federationadminv1.ChangeType_CHANGE_TYPE_UPDATE:
		article, err := s.articles.GetArticle(articleID)
		if err != nil {
			return activity{}, err
		}

		activityType, suffix := "Create", "create"
		if changeType == federationadminv1.ChangeType_CHANGE_TYPE_UPDATE {
			activityType, suffix = "Update", "update"
		}

		obj := buildNote(s.cfg.SiteBaseURL, articleID, article.Title, article.PublishedAt)
		activityID := s.cfg.SiteBaseURL + "/api/articles/" + articleID + "/" + suffix
		return buildAndSignActivity(activityType, activityID, s.cfg.SiteBaseURL, obj, s.cfg.ActorPrivateKeyPEM)

	case federationadminv1.ChangeType_CHANGE_TYPE_DELETE:
		obj := buildDeleteNote(s.cfg.SiteBaseURL, articleID)
		activityID := s.cfg.SiteBaseURL + "/api/articles/" + articleID + "/delete"
		return buildAndSignActivity("Delete", activityID, s.cfg.SiteBaseURL, obj, s.cfg.ActorPrivateKeyPEM)

	default:
		return activity{}, errUnspecifiedChangeType
	}
}

func (s *Service) deliver(ctx context.Context, inboxURL string, body []byte) error {
	keyID := s.cfg.SiteBaseURL + "/actor#main-key"
	headers, err := federation.SignHTTPRequest(inboxURL, http.MethodPost, string(body), keyID, s.cfg.ActorPrivateKeyPEM)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, inboxURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Date", headers.Date)
	req.Header.Set("Digest", headers.Digest)
	req.Header.Set("Signature", headers.Signature)
	req.Header.Set("Content-Type", "application/activity+json")
	req.Header.Set("Accept", "application/activity+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("federationadmin: delivery failed with status %d", resp.StatusCode)
	}
	return nil
}
