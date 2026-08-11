package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nagutabby/sveltekit-blog/backend/internal/contact"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federationadmin"
	"github.com/nagutabby/sveltekit-blog/backend/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := os.Getenv("BACKEND_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	contentDir := os.Getenv("CONTENT_DIR")
	if contentDir == "" {
		contentDir = "content"
	}

	siteBaseURL := os.Getenv("SITE_BASE_URL")
	if siteBaseURL == "" {
		siteBaseURL = "https://blog.nagutabby.uk"
	}

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	queries := db.New(pool)
	contentLoader := content.NewLoader(contentDir)

	federationCfg := federation.Config{
		SiteBaseURL:        siteBaseURL,
		WebBaseURL:         os.Getenv("WEB_BASE_URL"),
		ActorPublicKeyPEM:  os.Getenv("ACTOR_PUBLIC_KEY_PEM"),
		ActorPrivateKeyPEM: os.Getenv("ACTOR_PRIVATE_KEY_PEM"),
	}

	cfg := server.Config{
		Contact: contact.Config{
			APIToken:    os.Getenv("EMAIL_API_TOKEN"),
			FromAddress: os.Getenv("FROM_ADDRESS"),
			BCCAddress:  os.Getenv("BCC_ADDRESS"),
		},
		Content:         contentLoader,
		Federation:      federation.NewHandlers(queries, queries, contentLoader, federationCfg),
		FederationAdmin: federationadmin.NewService(contentLoader, queries, federationCfg),
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: server.NewHandler(cfg),
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("backend server listening on %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	}
}
