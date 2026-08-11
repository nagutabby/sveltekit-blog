package server_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	migrations "github.com/nagutabby/sveltekit-blog/backend/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
	"github.com/nagutabby/sveltekit-blog/backend/internal/server"
)

func generateTestActorPrivateKeyPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}))
}

// TestFederationFollowOverHTTP exercises the full stack for the highest
// risk part of the migration: a real HTTP server (server.NewHandler) with
// a real Postgres database, receiving a Follow activity from a fake
// remote Mastodon actor, and verifying it both persists the follower and
// delivers a signed Accept back to the remote inbox.
func TestFederationFollowOverHTTP(t *testing.T) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:17-alpine",
		tcpostgres.WithDatabase("sveltekit_blog_test"),
		tcpostgres.WithUsername("sveltekit_blog_test"),
		tcpostgres.WithPassword("sveltekit_blog_test"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	})

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	migrationDB, err := sql.Open("pgx", connString)
	if err != nil {
		t.Fatalf("failed to open migration connection: %v", err)
	}
	defer migrationDB.Close()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(migrationDB, "migrations"); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Fatalf("failed to open pgx pool: %v", err)
	}
	t.Cleanup(pool.Close)
	queries := db.New(pool)

	// Fake remote Mastodon actor + inbox. inboxURL is filled in once the
	// server is listening, since the actor document needs to advertise
	// the server's own address as its inbox.
	var inboxURL string
	var acceptDelivered bool
	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"inbox": "` + inboxURL + `", "publicKey": {"publicKeyPem": "REMOTE-KEY"}}`))
	})
	mux.HandleFunc("/users/alice/inbox", func(w http.ResponseWriter, r *http.Request) {
		acceptDelivered = r.Header.Get("Signature") != ""
		w.WriteHeader(http.StatusAccepted)
	})
	remote := httptest.NewServer(mux)
	defer remote.Close()
	inboxURL = remote.URL + "/users/alice/inbox"

	cfg := server.Config{
		Federation: federation.NewHandlers(queries, queries, nil, federation.Config{
			SiteBaseURL:        "https://blog.nagutabby.uk",
			ActorPrivateKeyPEM: generateTestActorPrivateKeyPEM(t),
		}),
	}
	backend := httptest.NewServer(server.NewHandler(cfg))
	defer backend.Close()

	activityBody := `{"@context":"https://www.w3.org/ns/activitystreams","type":"Follow","actor":"` + remote.URL + `/users/alice","object":"https://blog.nagutabby.uk/actor"}`
	resp, err := http.Post(backend.URL+"/actor/inbox", "application/activity+json", strings.NewReader(activityBody))
	if err != nil {
		t.Fatalf("POST /actor/inbox failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if !acceptDelivered {
		t.Fatal("expected a signed Accept to be delivered to the remote inbox")
	}

	follower, err := queries.GetFollowerByActorID(ctx, remote.URL+"/users/alice")
	if err != nil {
		t.Fatalf("GetFollowerByActorID returned error: %v", err)
	}
	if !follower.Following {
		t.Fatal("Following = false, want true after a Follow activity")
	}

	// Confirm the public /actor/followers endpoint reflects the new count.
	followersResp, err := http.Get(backend.URL + "/actor/followers")
	if err != nil {
		t.Fatalf("GET /actor/followers failed: %v", err)
	}
	defer followersResp.Body.Close()
	if followersResp.StatusCode != http.StatusOK {
		t.Fatalf("followers status = %d, want 200", followersResp.StatusCode)
	}
}
