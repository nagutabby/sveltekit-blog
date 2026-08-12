// Command migrate-to-d1 is a one-time operational tool that copies the
// production Follower/RelayConnection rows from Neon (PostgreSQL) into
// Cloudflare D1, preserving their exact following/connected state and
// timestamps (unlike db.Querier's Upsert* methods, which are written for
// live ActivityPub events and always set following/connected=true).
//
// It is idempotent (upserts on actorId) and defaults to a dry run; pass
// -apply to actually write to D1. Delete this command once the production
// cutover to D1 is complete — it has no reason to exist afterwards.
//
//	SOURCE_DATABASE_URL=postgres://... \
//	CLOUDFLARE_ACCOUNT_ID=... CLOUDFLARE_D1_DATABASE_ID=... CLOUDFLARE_D1_API_TOKEN=... \
//	go run ./cmd/migrate-to-d1 -apply
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/nagutabby/sveltekit-blog/backend/internal/db/d1"
)

type followerRow struct {
	ActorID      string
	Inbox        string
	PublicKeyPem string
	Following    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type relayConnectionRow struct {
	ActorID        string
	Inbox          string
	Connected      bool
	LastAcceptedAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apply := flag.Bool("apply", false, "write to D1; without this flag, only prints what would be migrated")
	flag.Parse()

	ctx := context.Background()

	sourceURL := os.Getenv("SOURCE_DATABASE_URL")
	if sourceURL == "" {
		return fmt.Errorf("SOURCE_DATABASE_URL is required (the Neon/PostgreSQL connection string to migrate from)")
	}
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	databaseID := os.Getenv("CLOUDFLARE_D1_DATABASE_ID")
	apiToken := os.Getenv("CLOUDFLARE_D1_API_TOKEN")
	if accountID == "" || databaseID == "" || apiToken == "" {
		return fmt.Errorf("CLOUDFLARE_ACCOUNT_ID, CLOUDFLARE_D1_DATABASE_ID, and CLOUDFLARE_D1_API_TOKEN are all required")
	}

	conn, err := pgx.Connect(ctx, sourceURL)
	if err != nil {
		return fmt.Errorf("connect to source database: %w", err)
	}
	defer conn.Close(ctx)

	followers, err := readFollowers(ctx, conn)
	if err != nil {
		return fmt.Errorf("read followers: %w", err)
	}
	relayConnections, err := readRelayConnections(ctx, conn)
	if err != nil {
		return fmt.Errorf("read relay connections: %w", err)
	}

	log.Printf("found %d follower row(s) and %d relay connection row(s) to migrate", len(followers), len(relayConnections))
	if !*apply {
		log.Print("dry run (pass -apply to write to D1); previewing rows:")
		for _, f := range followers {
			log.Printf("  Follower actorId=%s following=%v", f.ActorID, f.Following)
		}
		for _, r := range relayConnections {
			log.Printf("  RelayConnection actorId=%s connected=%v", r.ActorID, r.Connected)
		}
		return nil
	}

	client := d1.New(d1.Config{AccountID: accountID, DatabaseID: databaseID, APIToken: apiToken})

	for _, f := range followers {
		if err := writeFollower(ctx, client, f); err != nil {
			return fmt.Errorf("write follower %s: %w", f.ActorID, err)
		}
		log.Printf("migrated follower %s", f.ActorID)
	}
	for _, r := range relayConnections {
		if err := writeRelayConnection(ctx, client, r); err != nil {
			return fmt.Errorf("write relay connection %s: %w", r.ActorID, err)
		}
		log.Printf("migrated relay connection %s", r.ActorID)
	}

	log.Print("migration complete")
	return nil
}

func readFollowers(ctx context.Context, conn *pgx.Conn) ([]followerRow, error) {
	rows, err := conn.Query(ctx, `SELECT "actorId", "inbox", "publicKeyPem", "following", "createdAt", "updatedAt" FROM "Follower"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []followerRow
	for rows.Next() {
		var f followerRow
		if err := rows.Scan(&f.ActorID, &f.Inbox, &f.PublicKeyPem, &f.Following, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func readRelayConnections(ctx context.Context, conn *pgx.Conn) ([]relayConnectionRow, error) {
	rows, err := conn.Query(ctx, `SELECT "actorId", "inbox", "connected", "lastAcceptedAt", "createdAt", "updatedAt" FROM "RelayConnection"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []relayConnectionRow
	for rows.Next() {
		var r relayConnectionRow
		if err := rows.Scan(&r.ActorID, &r.Inbox, &r.Connected, &r.LastAcceptedAt, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// rfc3339 formats t the way the SQLite/D1 schema stores timestamps,
// matching internal/federation's nowTimestamp.
func rfc3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func writeFollower(ctx context.Context, client *d1.Client, f followerRow) error {
	_, err := client.Exec(ctx, `
		INSERT INTO "Follower" ("actorId", "inbox", "publicKeyPem", "following", "createdAt", "updatedAt")
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT ("actorId") DO UPDATE SET
		    "inbox" = excluded."inbox",
		    "publicKeyPem" = excluded."publicKeyPem",
		    "following" = excluded."following",
		    "updatedAt" = excluded."updatedAt"`,
		[]any{f.ActorID, f.Inbox, f.PublicKeyPem, f.Following, rfc3339(f.CreatedAt), rfc3339(f.UpdatedAt)},
	)
	return err
}

func writeRelayConnection(ctx context.Context, client *d1.Client, r relayConnectionRow) error {
	var lastAcceptedAt any
	if r.LastAcceptedAt != nil {
		lastAcceptedAt = rfc3339(*r.LastAcceptedAt)
	}
	_, err := client.Exec(ctx, `
		INSERT INTO "RelayConnection" ("actorId", "inbox", "connected", "lastAcceptedAt", "createdAt", "updatedAt")
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT ("actorId") DO UPDATE SET
		    "inbox" = excluded."inbox",
		    "connected" = excluded."connected",
		    "lastAcceptedAt" = excluded."lastAcceptedAt",
		    "updatedAt" = excluded."updatedAt"`,
		[]any{r.ActorID, r.Inbox, r.Connected, lastAcceptedAt, rfc3339(r.CreatedAt), rfc3339(r.UpdatedAt)},
	)
	return err
}
