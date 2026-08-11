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

	"github.com/nagutabby/sveltekit-blog/backend/internal/contact"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	addr := os.Getenv("BACKEND_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	contentDir := os.Getenv("CONTENT_DIR")
	if contentDir == "" {
		// PR6 moves the Markdown sources into backend/content; until then
		// they still live under web/static/content.
		contentDir = "../web/static/content"
	}

	cfg := server.Config{
		Contact: contact.Config{
			APIToken:    os.Getenv("EMAIL_API_TOKEN"),
			FromAddress: os.Getenv("FROM_ADDRESS"),
			BCCAddress:  os.Getenv("BCC_ADDRESS"),
		},
		Content: content.NewLoader(contentDir),
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: server.NewHandler(cfg),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
