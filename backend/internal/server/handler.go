package server

import (
	"net/http"

	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/contact/v1/contactv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/content/v1/contentv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/federationadmin/v1/federationadminv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/health/v1/healthv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/internal/contact"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federationadmin"
	"github.com/nagutabby/sveltekit-blog/backend/internal/health"
)

// Config holds the configuration for every service mounted by NewHandler.
// Content and Federation are optional; when nil, the corresponding routes
// are not mounted (used by tests that only care about other services).
type Config struct {
	Contact         contact.Config
	Content         *content.Loader
	Federation      *federation.Handlers
	FederationAdmin *federationadmin.Service
	// AllowedOrigin is the browser origin allowed to call ContactService
	// cross-origin (blog.nagutabby.uk calling api.nagutabby.uk). Empty
	// disables CORS handling (same-origin deployments don't need it).
	AllowedOrigin string
}

// NewHandler builds the top-level HTTP handler for the backend server.
// Public ActivityPub endpoints and further Connect RPC services are
// registered here as they migrate over in later PRs.
func NewHandler(cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)

	healthPath, healthHandler := healthv1connect.NewHealthServiceHandler(health.NewService())
	mux.Handle(healthPath, healthHandler)

	contactPath, contactHandler := contactv1connect.NewContactServiceHandler(contact.NewService(cfg.Contact))
	if cfg.AllowedOrigin != "" {
		contactHandler = withCORS(cfg.AllowedOrigin, contactHandler)
	}
	mux.Handle(contactPath, contactHandler)

	if cfg.Content != nil {
		contentPath, contentHandler := contentv1connect.NewContentServiceHandler(content.NewService(cfg.Content))
		mux.Handle(contentPath, contentHandler)
	}

	if cfg.Federation != nil {
		mux.HandleFunc("GET /.well-known/webfinger", cfg.Federation.Webfinger)
		mux.HandleFunc("GET /actor", cfg.Federation.Actor)
		mux.HandleFunc("GET /actor/followers", cfg.Federation.Followers)
		mux.HandleFunc("GET /actor/following", cfg.Federation.Following)
		mux.HandleFunc("GET /actor/outbox", cfg.Federation.Outbox)
		mux.HandleFunc("POST /actor/inbox", cfg.Federation.Inbox)
		mux.HandleFunc("GET /api/articles/{name}", cfg.Federation.ArticleNote)
	}

	if cfg.FederationAdmin != nil {
		federationAdminPath, federationAdminHandler := federationadminv1connect.NewFederationAdminServiceHandler(cfg.FederationAdmin)
		mux.Handle(federationAdminPath, federationAdminHandler)
	}

	return mux
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
