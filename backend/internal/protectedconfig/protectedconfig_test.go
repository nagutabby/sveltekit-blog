package protectedconfig_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

	"github.com/nagutabby/sveltekit-blog/backend/internal/protectedconfig"
)

type fakeClient struct {
	secrets map[string]string
}

func (f *fakeClient) GetSecretValue(_ context.Context, params *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	value, ok := f.secrets[*params.SecretId]
	if !ok {
		return nil, &types.ResourceNotFoundException{}
	}
	return &secretsmanager.GetSecretValueOutput{SecretString: &value}, nil
}

func TestLoadActorKeysFallsBackToEnvVarsWithoutArn(t *testing.T) {
	t.Setenv("ACTOR_KEYS_SECRET_ARN", "")
	t.Setenv("ACTOR_PUBLIC_KEY_PEM", "PUBLIC-FROM-ENV")
	t.Setenv("ACTOR_PRIVATE_KEY_PEM", "PRIVATE-FROM-ENV")

	keys, err := protectedconfig.LoadActorKeys(context.Background(), &fakeClient{})
	if err != nil {
		t.Fatalf("LoadActorKeys returned error: %v", err)
	}
	if keys.PublicKeyPEM != "PUBLIC-FROM-ENV" || keys.PrivateKeyPEM != "PRIVATE-FROM-ENV" {
		t.Fatalf("keys = %+v, want values from env vars", keys)
	}
}

func TestLoadActorKeysReadsFromSecretsManagerWhenArnIsSet(t *testing.T) {
	const arn = "arn:aws:secretsmanager:us-east-1:123456789012:secret:actor-keys"
	t.Setenv("ACTOR_KEYS_SECRET_ARN", arn)

	client := &fakeClient{secrets: map[string]string{
		arn: `{"ACTOR_PUBLIC_KEY_PEM":"PUBLIC-FROM-SECRET","ACTOR_PRIVATE_KEY_PEM":"PRIVATE-FROM-SECRET"}`,
	}}

	keys, err := protectedconfig.LoadActorKeys(context.Background(), client)
	if err != nil {
		t.Fatalf("LoadActorKeys returned error: %v", err)
	}
	if keys.PublicKeyPEM != "PUBLIC-FROM-SECRET" || keys.PrivateKeyPEM != "PRIVATE-FROM-SECRET" {
		t.Fatalf("keys = %+v, want values from the secret", keys)
	}
}

func TestLoadEmailAPITokenFallsBackToEnvVarWithoutArn(t *testing.T) {
	t.Setenv("EMAIL_API_TOKEN_SECRET_ARN", "")
	t.Setenv("EMAIL_API_TOKEN", "TOKEN-FROM-ENV")

	token, err := protectedconfig.LoadEmailAPIToken(context.Background(), &fakeClient{})
	if err != nil {
		t.Fatalf("LoadEmailAPIToken returned error: %v", err)
	}
	if token != "TOKEN-FROM-ENV" {
		t.Fatalf("token = %q, want %q", token, "TOKEN-FROM-ENV")
	}
}

func TestLoadEmailAPITokenReadsFromSecretsManagerWhenArnIsSet(t *testing.T) {
	const arn = "arn:aws:secretsmanager:us-east-1:123456789012:secret:email-api-token"
	t.Setenv("EMAIL_API_TOKEN_SECRET_ARN", arn)

	client := &fakeClient{secrets: map[string]string{arn: "TOKEN-FROM-SECRET"}}

	token, err := protectedconfig.LoadEmailAPIToken(context.Background(), client)
	if err != nil {
		t.Fatalf("LoadEmailAPIToken returned error: %v", err)
	}
	if token != "TOKEN-FROM-SECRET" {
		t.Fatalf("token = %q, want %q", token, "TOKEN-FROM-SECRET")
	}
}
