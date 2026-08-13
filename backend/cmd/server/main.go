package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	contentembed "github.com/nagutabby/sveltekit-blog/backend/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/contact"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db/d1"
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
		if port := os.Getenv("PORT"); port != "" {
			addr = ":" + port
		} else {
			addr = ":8080"
		}
	}

	// CONTENT_DIR overrides the content source with a plain OS directory
	// (useful for local development). Its absence is the production
	// path: the embedded FS, since a relative os.DirFS("content") isn't
	// reliably present in Vercel's Go serverless runtime.
	var contentFS fs.FS = contentembed.FS
	if contentDir := os.Getenv("CONTENT_DIR"); contentDir != "" {
		contentFS = os.DirFS(contentDir)
	}

	siteBaseURL := os.Getenv("SITE_BASE_URL")
	if siteBaseURL == "" {
		siteBaseURL = "https://blog.nagutabby.uk"
	}

	queries, closeDB, err := newQueries()
	if err != nil {
		return err
	}
	defer closeDB()
	contentLoader := content.NewLoader(contentFS)

	federationCfg := federation.Config{
		SiteBaseURL:        siteBaseURL,
		ActorPublicKeyPEM:  os.Getenv("ACTOR_PUBLIC_KEY_PEM"),
		ActorPrivateKeyPEM: os.Getenv("ACTOR_PRIVATE_KEY_PEM"),
	}

	cfg := server.Config{
		Contact: contact.Config{
			APIToken:    os.Getenv("EMAIL_API_TOKEN"),
			FromAddress: os.Getenv("FROM_ADDRESS"),
			BCCAddress:  os.Getenv("BCC_ADDRESS"),
		},
		Federation:           federation.NewHandlers(queries, queries, contentLoader, federationCfg),
		FederationAdmin:      federationadmin.NewService(contentLoader, queries, federationCfg),
		FederationAdminToken: os.Getenv("FEDERATION_ADMIN_TOKEN"),
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

// newQueries builds the db.Querier the server runs against. In production
// this is Cloudflare D1 (reached over its HTTP query API, since the JS
// Workers Binding API isn't available outside a Worker); locally, and
// wherever the D1 environment variables aren't set, it's a plain SQLite
// file, since D1 is SQLite-compatible and needs no separate driver.
func newQueries() (db.Querier, func(), error) {
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	databaseID := os.Getenv("CLOUDFLARE_D1_DATABASE_ID")
	apiToken := os.Getenv("CLOUDFLARE_D1_API_TOKEN")
	if accountID != "" && databaseID != "" && apiToken != "" {
		client := d1.New(d1.Config{
			AccountID:  accountID,
			DatabaseID: databaseID,
			APIToken:   apiToken,
		})
		return client, func() {}, nil
	}

	sqlitePath := os.Getenv("SQLITE_PATH")
	if sqlitePath == "" {
		sqlitePath = "backend.db"
	}
	sqlDB, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		return nil, nil, fmt.Errorf("open local sqlite database %q: %w", sqlitePath, err)
	}
	return db.New(sqlDB), func() { _ = sqlDB.Close() }, nil
}
