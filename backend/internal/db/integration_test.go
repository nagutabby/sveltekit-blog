package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	migrations "github.com/nagutabby/sveltekit-blog/backend/db"
	db "github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

// setupDB starts a disposable Postgres container, applies the goose
// migrations embedded in backend/db, and returns sqlc Queries backed by it.
func setupDB(t *testing.T) *db.Queries {
	t.Helper()
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

	return db.New(pool)
}

func TestFollowerLifecycle(t *testing.T) {
	queries := setupDB(t)
	ctx := context.Background()

	const actorID = "https://mastodon.example/users/alice"

	created, err := queries.UpsertFollower(ctx, db.UpsertFollowerParams{
		ActorId:      actorID,
		Inbox:        "https://mastodon.example/users/alice/inbox",
		PublicKeyPem: "PEM-1",
	})
	if err != nil {
		t.Fatalf("UpsertFollower (create) returned error: %v", err)
	}
	if !created.Following {
		t.Fatalf("created.Following = false, want true")
	}

	count, err := queries.CountActiveFollowers(ctx)
	if err != nil {
		t.Fatalf("CountActiveFollowers returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("CountActiveFollowers = %d, want 1", count)
	}

	// Re-following with a refreshed key must update in place, not duplicate.
	updated, err := queries.UpsertFollower(ctx, db.UpsertFollowerParams{
		ActorId:      actorID,
		Inbox:        "https://mastodon.example/users/alice/inbox",
		PublicKeyPem: "PEM-2",
	})
	if err != nil {
		t.Fatalf("UpsertFollower (update) returned error: %v", err)
	}
	if updated.ID != created.ID {
		t.Fatalf("UpsertFollower created a new row: got id %d, want %d", updated.ID, created.ID)
	}
	if updated.PublicKeyPem != "PEM-2" {
		t.Fatalf("PublicKeyPem = %q, want %q", updated.PublicKeyPem, "PEM-2")
	}

	countAfterReFollow, err := queries.CountActiveFollowers(ctx)
	if err != nil {
		t.Fatalf("CountActiveFollowers returned error: %v", err)
	}
	if countAfterReFollow != 1 {
		t.Fatalf("CountActiveFollowers after re-follow = %d, want 1 (no duplicate row)", countAfterReFollow)
	}

	unfollowed, err := queries.UnfollowByActorID(ctx, db.UnfollowByActorIDParams{
		ActorId:      actorID,
		Inbox:        "https://mastodon.example/users/alice/inbox",
		PublicKeyPem: "PEM-2",
	})
	if err != nil {
		t.Fatalf("UnfollowByActorID returned error: %v", err)
	}
	if unfollowed.Following {
		t.Fatalf("unfollowed.Following = true, want false")
	}

	countAfterUnfollow, err := queries.CountActiveFollowers(ctx)
	if err != nil {
		t.Fatalf("CountActiveFollowers returned error: %v", err)
	}
	if countAfterUnfollow != 0 {
		t.Fatalf("CountActiveFollowers after unfollow = %d, want 0", countAfterUnfollow)
	}

	got, err := queries.GetFollowerByActorID(ctx, actorID)
	if err != nil {
		t.Fatalf("GetFollowerByActorID returned error: %v", err)
	}
	if got.ActorId != actorID {
		t.Fatalf("ActorId = %q, want %q", got.ActorId, actorID)
	}
}

func TestUnfollowByActorIDWithoutExistingFollowerFails(t *testing.T) {
	queries := setupDB(t)
	ctx := context.Background()

	_, err := queries.UnfollowByActorID(ctx, db.UnfollowByActorIDParams{
		ActorId:      "https://mastodon.example/users/never-followed",
		Inbox:        "https://mastodon.example/users/never-followed/inbox",
		PublicKeyPem: "PEM",
	})
	if err == nil {
		t.Fatal("expected an error unfollowing an actor that never followed, got nil")
	}
}

func TestRelayConnectionLifecycle(t *testing.T) {
	queries := setupDB(t)
	ctx := context.Background()

	const actorID = "https://relay.example/actor"

	first, err := queries.UpsertRelayConnectionAccepted(ctx, db.UpsertRelayConnectionAcceptedParams{
		ActorId: actorID,
		Inbox:   "https://relay.example/inbox",
	})
	if err != nil {
		t.Fatalf("UpsertRelayConnectionAccepted (create) returned error: %v", err)
	}
	if !first.Connected {
		t.Fatalf("Connected = false, want true")
	}
	if !first.LastAcceptedAt.Valid {
		t.Fatalf("LastAcceptedAt is not set")
	}

	time.Sleep(10 * time.Millisecond)

	second, err := queries.UpsertRelayConnectionAccepted(ctx, db.UpsertRelayConnectionAcceptedParams{
		ActorId: actorID,
		Inbox:   "https://relay.example/inbox",
	})
	if err != nil {
		t.Fatalf("UpsertRelayConnectionAccepted (re-accept) returned error: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("UpsertRelayConnectionAccepted created a new row: got id %d, want %d", second.ID, first.ID)
	}
	if !second.LastAcceptedAt.Time.After(first.LastAcceptedAt.Time) {
		t.Fatalf("LastAcceptedAt did not advance on re-accept: first=%v second=%v", first.LastAcceptedAt.Time, second.LastAcceptedAt.Time)
	}

	list, err := queries.ListRelayConnections(ctx)
	if err != nil {
		t.Fatalf("ListRelayConnections returned error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("ListRelayConnections returned %d rows, want 1", len(list))
	}

	got, err := queries.GetRelayConnectionByActorID(ctx, actorID)
	if err != nil {
		t.Fatalf("GetRelayConnectionByActorID returned error: %v", err)
	}
	if got.ActorId != actorID {
		t.Fatalf("ActorId = %q, want %q", got.ActorId, actorID)
	}
}
