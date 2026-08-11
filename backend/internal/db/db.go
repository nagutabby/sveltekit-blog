// Package db provides DynamoDB-backed storage for the Follower and
// RelayConnection tables used by ActivityPub federation.
package db

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// TimestampLayout is RFC3339 with a fixed-width (zero-padded) nanosecond
// component, so stored timestamps remain correctly sortable as plain
// strings. time.RFC3339Nano trims trailing zeros, which would otherwise
// misorder a timestamp landing exactly on a whole second against one with a
// nonzero fraction later in the same second. Exported so tools that write
// Follower/RelayConnection items directly (e.g. cmd/migrate-to-dynamo) format
// timestamps the same way.
const TimestampLayout = "2006-01-02T15:04:05.000000000Z"

func nowTimestamp() string {
	return time.Now().UTC().Format(TimestampLayout)
}

// Client is the subset of *dynamodb.Client the Queries methods need,
// defined as an interface so tests can substitute a fake.
type Client interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

// Queries provides the Follower/RelayConnection persistence operations
// federation and federationadmin depend on.
type Queries struct {
	client               Client
	followerTable        string
	relayConnectionTable string
}

func New(client Client, followerTable, relayConnectionTable string) *Queries {
	return &Queries{
		client:               client,
		followerTable:        followerTable,
		relayConnectionTable: relayConnectionTable,
	}
}
