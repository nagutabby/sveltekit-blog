package health

import (
	"context"
	"testing"

	"connectrpc.com/connect"

	healthv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/health/v1"
)

func TestServiceCheck(t *testing.T) {
	svc := NewService()

	resp, err := svc.Check(context.Background(), connect.NewRequest(&healthv1.CheckRequest{}))
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if got := resp.Msg.GetStatus(); got != "ok" {
		t.Fatalf("status = %q, want %q", got, "ok")
	}
}
