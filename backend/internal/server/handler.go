package server

import (
	"net/http"

	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/contact/v1/contactv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/health/v1/healthv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/internal/contact"
	"github.com/nagutabby/sveltekit-blog/backend/internal/health"
)

// Config holds the configuration for every Connect RPC service mounted by
// NewHandler.
type Config struct {
	Contact contact.Config
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

	return mux
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
