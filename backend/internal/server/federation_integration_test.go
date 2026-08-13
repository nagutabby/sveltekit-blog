package server_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"

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

	dbPath := filepath.Join(t.TempDir(), fmt.Sprintf("test-%d.db", time.Now().UnixNano()))
	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	queries := db.New(sqlDB)

	// Fake remote Mastodon actor + inbox. inboxURL is filled in once the
	// server is listening, since the actor document needs to advertise
	// the server's own address as its inbox.
	aliceKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	aliceKeyDER, err := x509.MarshalPKIXPublicKey(&aliceKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal alice's public key: %v", err)
	}
	alicePublicKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: aliceKeyDER}))
	alicePrivateKeyPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(aliceKey),
	}))

	var inboxURL string
	var acceptDelivered bool
	mux := http.NewServeMux()
	mux.HandleFunc("/users/alice", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"inbox": "` + inboxURL + `", "publicKey": {"publicKeyPem": ` + strconv.Quote(alicePublicKeyPEM) + `}}`))
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

	// The inbox verifies the sender's HTTP Signature, so the Follow must be
	// signed with alice's key the same way a real remote server would sign
	// it — host/date/digest all keyed to backend's actual (random-port)
	// address, since that's what the server sees as the request's Host.
	activityBody := `{"@context":"https://www.w3.org/ns/activitystreams","type":"Follow","actor":"` + remote.URL + `/users/alice","object":"https://blog.nagutabby.uk/actor"}`
	signedHeaders, err := federation.SignHTTPRequest(backend.URL+"/actor/inbox", http.MethodPost, activityBody, remote.URL+"/users/alice#main-key", alicePrivateKeyPEM)
	if err != nil {
		t.Fatalf("failed to sign Follow activity: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, backend.URL+"/actor/inbox", strings.NewReader(activityBody))
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/activity+json")
	req.Header.Set("Date", signedHeaders.Date)
	req.Header.Set("Digest", signedHeaders.Digest)
	req.Header.Set("Signature", signedHeaders.Signature)

	resp, err := http.DefaultClient.Do(req)
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
