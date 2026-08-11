// create-tables creates the Follower/RelayConnection tables against a
// dynamodb-local instance for local development (docker compose up -d).
// The real AWS tables are managed by infra/lib/dynamodb-stack.ts.
package main

import (
	"context"
	"errors"
	"log"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	// dynamodb-localのみを対象にするツールなので、実際のAWS認証情報解決
	// (IMDS/SSO等)には頼らずダミー認証情報を直接使う。デフォルトの認証情報
	// チェーンを辿ると、ローカル環境によっては解決に時間がかかることがある。
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		return err
	}
	client := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	})

	for _, def := range []*dynamodb.CreateTableInput{
		db.FollowerTableDefinition("Follower"),
		db.RelayConnectionTableDefinition("RelayConnection"),
	} {
		_, err := client.CreateTable(ctx, def)
		var inUse *types.ResourceInUseException
		if errors.As(err, &inUse) {
			log.Printf("table %s already exists, skipping", *def.TableName)
			continue
		}
		if err != nil {
			return err
		}
		log.Printf("created table %s", *def.TableName)
	}
	return nil
}
