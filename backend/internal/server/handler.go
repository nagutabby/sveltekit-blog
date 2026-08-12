package server

import (
	"net/http"

	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/contact/v1/contactv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/federationadmin/v1/federationadminv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/health/v1/healthv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/internal/contact"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federationadmin"
	"github.com/nagutabby/sveltekit-blog/backend/internal/health"
)

// Config holds the configuration for every service mounted by NewHandler.
// Federation is optional; when nil, its routes are not mounted (used by
// tests that only care about other services).
type Config struct {
	Contact         contact.Config
	Federation      *federation.Handlers
	FederationAdmin *federationadmin.Service

	// FederationAdminToken guards FederationAdmin: requests must send
	// "Authorization: Bearer <FederationAdminToken>". Required (and
	// non-empty) once FederationAdmin is exposed publicly rather than
	// reached only from localhost/a private network, since this service
	// can trigger signed ActivityPub deliveries to every relay.
	FederationAdminToken string
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
	mux.Handle(contactPath, contactHandler)

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
		mux.Handle(federationAdminPath, requireBearerToken(cfg.FederationAdminToken, federationAdminHandler))
	}

	return mux
}

// requireBearerToken rejects any request whose Authorization header isn't
// "Bearer <token>". An empty token always rejects (fail closed) rather
// than leaving the route open, since an unset token most likely means
// misconfiguration rather than "no auth needed".
func requireBearerToken(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token == "" || r.Header.Get("Authorization") != "Bearer "+token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
