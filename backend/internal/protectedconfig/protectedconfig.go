// Package protectedconfig resolves secret-shaped config values either from
// AWS Secrets Manager (when a *_SECRET_ARN env var is set, i.e. running on
// Lambda) or directly from a plain env var (Railway/local dev, where
// Secrets Manager isn't wired up). This lets internal/app work unchanged
// across both deployment targets.
package protectedconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Client is the subset of *secretsmanager.Client this package needs,
// defined as an interface so tests can substitute a fake.
type Client interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type ActorKeys struct {
	PublicKeyPEM  string `json:"ACTOR_PUBLIC_KEY_PEM"`
	PrivateKeyPEM string `json:"ACTOR_PRIVATE_KEY_PEM"`
}

// LoadActorKeys reads the ActivityPub actor's key pair from the secret
// named by ACTOR_KEYS_SECRET_ARN (a JSON object with ACTOR_PUBLIC_KEY_PEM/
// ACTOR_PRIVATE_KEY_PEM keys), or from the ACTOR_PUBLIC_KEY_PEM/
// ACTOR_PRIVATE_KEY_PEM env vars directly if that ARN isn't set.
func LoadActorKeys(ctx context.Context, client Client) (ActorKeys, error) {
	arn := os.Getenv("ACTOR_KEYS_SECRET_ARN")
	if arn == "" {
		return ActorKeys{
			PublicKeyPEM:  os.Getenv("ACTOR_PUBLIC_KEY_PEM"),
			PrivateKeyPEM: os.Getenv("ACTOR_PRIVATE_KEY_PEM"),
		}, nil
	}
	value, err := getSecretString(ctx, client, arn)
	if err != nil {
		return ActorKeys{}, err
	}
	var keys ActorKeys
	if err := json.Unmarshal([]byte(value), &keys); err != nil {
		return ActorKeys{}, fmt.Errorf("protectedconfig: parsing ActorKeys secret: %w", err)
	}
	return keys, nil
}

// LoadEmailAPIToken reads the contact form's email API token from the
// secret named by EMAIL_API_TOKEN_SECRET_ARN, or from the EMAIL_API_TOKEN
// env var directly if that ARN isn't set.
func LoadEmailAPIToken(ctx context.Context, client Client) (string, error) {
	arn := os.Getenv("EMAIL_API_TOKEN_SECRET_ARN")
	if arn == "" {
		return os.Getenv("EMAIL_API_TOKEN"), nil
	}
	return getSecretString(ctx, client, arn)
}

func getSecretString(ctx context.Context, client Client, arn string) (string, error) {
	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &arn})
	if err != nil {
		return "", fmt.Errorf("protectedconfig: fetching secret %s: %w", arn, err)
	}
	if out.SecretString == nil {
		return "", fmt.Errorf("protectedconfig: secret %s has no SecretString", arn)
	}
	return *out.SecretString, nil
}
