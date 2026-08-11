// Package app wires together the backend's dependencies (DynamoDB, content
// loader, federation config) from environment variables into a single
// http.Handler. It is shared by cmd/server (a normal HTTP server, used on
// Railway) and cmd/lambda (wrapped for AWS Lambda via httpadapter), so the
// two entrypoints don't duplicate the wiring.
package app

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/nagutabby/sveltekit-blog/backend/internal/contact"
	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
	"github.com/nagutabby/sveltekit-blog/backend/internal/federationadmin"
	"github.com/nagutabby/sveltekit-blog/backend/internal/protectedconfig"
	"github.com/nagutabby/sveltekit-blog/backend/internal/server"
)

func NewHandler(ctx context.Context) (http.Handler, error) {
	contentDir := os.Getenv("CONTENT_DIR")
	if contentDir == "" {
		contentDir = "content"
	}

	siteBaseURL := os.Getenv("SITE_BASE_URL")
	if siteBaseURL == "" {
		siteBaseURL = "https://blog.nagutabby.uk"
	}

	followerTable := os.Getenv("FOLLOWER_TABLE_NAME")
	if followerTable == "" {
		followerTable = "Follower"
	}
	relayConnectionTable := os.Getenv("RELAY_CONNECTION_TABLE_NAME")
	if relayConnectionTable == "" {
		relayConnectionTable = "RelayConnection"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	dynamoClient := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		// dynamodb-localでのローカル開発用。本番ではDYNAMODB_ENDPOINTを設定しない。
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	queries := db.New(dynamoClient, followerTable, relayConnectionTable)
	contentLoader := content.NewLoader(contentDir)

	secretsClient := secretsmanager.NewFromConfig(awsCfg)
	actorKeys, err := protectedconfig.LoadActorKeys(ctx, secretsClient)
	if err != nil {
		return nil, err
	}
	emailAPIToken, err := protectedconfig.LoadEmailAPIToken(ctx, secretsClient)
	if err != nil {
		return nil, err
	}

	federationCfg := federation.Config{
		SiteBaseURL:        siteBaseURL,
		WebBaseURL:         os.Getenv("WEB_BASE_URL"),
		ActorPublicKeyPEM:  actorKeys.PublicKeyPEM,
		ActorPrivateKeyPEM: actorKeys.PrivateKeyPEM,
	}

	cfg := server.Config{
		Contact: contact.Config{
			APIToken:    emailAPIToken,
			FromAddress: os.Getenv("FROM_ADDRESS"),
			BCCAddress:  os.Getenv("BCC_ADDRESS"),
		},
		Content:         contentLoader,
		Federation:      federation.NewHandlers(queries, queries, contentLoader, federationCfg),
		FederationAdmin: federationadmin.NewService(contentLoader, queries, federationCfg),
		AllowedOrigin:   os.Getenv("CORS_ALLOWED_ORIGIN"),
	}

	return server.NewHandler(cfg), nil
}
