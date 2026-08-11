// AWS Lambda entrypoint. Wraps the same http.Handler cmd/server uses
// (via internal/app) with httpadapter instead of rewriting routing logic
// for Lambda's event format.
package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/nagutabby/sveltekit-blog/backend/internal/app"
)

func main() {
	handler, err := app.NewHandler(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	lambda.Start(httpadapter.New(handler).ProxyWithContext)
}
