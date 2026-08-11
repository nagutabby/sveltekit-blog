package db_test

import (
	"context"
	"testing"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	tcdynamodb "github.com/testcontainers/testcontainers-go/modules/dynamodb"

	db "github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

const (
	followerTable        = "Follower"
	relayConnectionTable = "RelayConnection"
)

// setupDB starts a disposable dynamodb-local container, creates the
// Follower/RelayConnection tables, and returns Queries backed by it.
func setupDB(t *testing.T) *db.Queries {
	t.Helper()
	ctx := context.Background()

	container, err := tcdynamodb.Run(ctx, "amazon/dynamodb-local:2.6.1")
	if err != nil {
		t.Fatalf("failed to start dynamodb-local container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate dynamodb-local container: %v", err)
		}
	})

	hostPort, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}
	endpoint := "http://" + hostPort

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}
	client := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	})

	for _, def := range []*dynamodb.CreateTableInput{
		db.FollowerTableDefinition(followerTable),
		db.RelayConnectionTableDefinition(relayConnectionTable),
	} {
		if _, err := client.CreateTable(ctx, def); err != nil {
			t.Fatalf("failed to create table %s: %v", *def.TableName, err)
		}
	}

	return db.New(client, followerTable, relayConnectionTable)
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
	if updated.ActorId != created.ActorId {
		t.Fatalf("UpsertFollower changed the key: got %q, want %q", updated.ActorId, created.ActorId)
	}
	if updated.CreatedAt != created.CreatedAt {
		t.Fatalf("UpsertFollower created a new row: CreatedAt changed from %q to %q", created.CreatedAt, updated.CreatedAt)
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
	if first.LastAcceptedAt == "" {
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
	if second.ActorId != first.ActorId {
		t.Fatalf("UpsertRelayConnectionAccepted changed the key: got %q, want %q", second.ActorId, first.ActorId)
	}
	if second.CreatedAt != first.CreatedAt {
		t.Fatalf("UpsertRelayConnectionAccepted created a new row: CreatedAt changed from %q to %q", first.CreatedAt, second.CreatedAt)
	}
	if second.LastAcceptedAt <= first.LastAcceptedAt {
		t.Fatalf("LastAcceptedAt did not advance on re-accept: first=%v second=%v", first.LastAcceptedAt, second.LastAcceptedAt)
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
