// migrate-to-dynamo is a one-shot CLI that copies every existing
// Follower/RelayConnection row from Postgres into DynamoDB. Run it once
// during cutover, after the DynamoDB tables exist (infra `cdk deploy`) and
// before backend traffic is switched over to the DynamoDB-backed Lambda.
//
// It writes items directly via PutItem (not internal/db's Upsert* methods)
// because it must preserve the exact historical Following/Connected state
// and timestamps; the Upsert* methods always set Following/Connected to
// true, which would be wrong for followers who had already unfollowed.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jackc/pgx/v5"

	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	pgConn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}
	defer pgConn.Close(ctx)

	followerTable := envOr("FOLLOWER_TABLE_NAME", "Follower")
	relayConnectionTable := envOr("RELAY_CONNECTION_TABLE_NAME", "RelayConnection")

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	dynamoClient := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	followerCount, err := migrateFollowers(ctx, pgConn, dynamoClient, followerTable)
	if err != nil {
		return fmt.Errorf("migrating followers: %w", err)
	}
	log.Printf("migrated %d follower(s) to table %s", followerCount, followerTable)

	relayCount, err := migrateRelayConnections(ctx, pgConn, dynamoClient, relayConnectionTable)
	if err != nil {
		return fmt.Errorf("migrating relay connections: %w", err)
	}
	log.Printf("migrated %d relay connection(s) to table %s", relayCount, relayConnectionTable)

	return nil
}

func migrateFollowers(ctx context.Context, pgConn *pgx.Conn, client *dynamodb.Client, table string) (int, error) {
	rows, err := pgConn.Query(ctx, `SELECT "actorId", inbox, "publicKeyPem", following, "createdAt", "updatedAt" FROM "Follower"`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var actorID, inbox, publicKeyPEM string
		var following bool
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&actorID, &inbox, &publicKeyPEM, &following, &createdAt, &updatedAt); err != nil {
			return count, err
		}

		item, err := attributevalue.MarshalMap(db.Follower{
			ActorId:      actorID,
			Inbox:        inbox,
			PublicKeyPem: publicKeyPEM,
			Following:    following,
			CreatedAt:    createdAt.UTC().Format(db.TimestampLayout),
			UpdatedAt:    updatedAt.UTC().Format(db.TimestampLayout),
		})
		if err != nil {
			return count, err
		}
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{TableName: &table, Item: item}); err != nil {
			return count, fmt.Errorf("actorId %s: %w", actorID, err)
		}
		count++
	}
	return count, rows.Err()
}

func migrateRelayConnections(ctx context.Context, pgConn *pgx.Conn, client *dynamodb.Client, table string) (int, error) {
	rows, err := pgConn.Query(ctx, `SELECT "actorId", inbox, connected, "lastAcceptedAt", "createdAt", "updatedAt" FROM "RelayConnection"`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var actorID, inbox string
		var connected bool
		var lastAcceptedAt *time.Time
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&actorID, &inbox, &connected, &lastAcceptedAt, &createdAt, &updatedAt); err != nil {
			return count, err
		}

		relayConnection := db.RelayConnection{
			ActorId:   actorID,
			Inbox:     inbox,
			Connected: connected,
			CreatedAt: createdAt.UTC().Format(db.TimestampLayout),
			UpdatedAt: updatedAt.UTC().Format(db.TimestampLayout),
		}
		if lastAcceptedAt != nil {
			relayConnection.LastAcceptedAt = lastAcceptedAt.UTC().Format(db.TimestampLayout)
		}

		item, err := attributevalue.MarshalMap(relayConnection)
		if err != nil {
			return count, err
		}
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{TableName: &table, Item: item}); err != nil {
			return count, fmt.Errorf("actorId %s: %w", actorID, err)
		}
		count++
	}
	return count, rows.Err()
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
