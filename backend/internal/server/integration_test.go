package server_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"

	healthv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/health/v1"
	"github.com/nagutabby/sveltekit-blog/backend/gen/blog/health/v1/healthv1connect"
	"github.com/nagutabby/sveltekit-blog/backend/internal/server"
)

// TestHealthServiceOverHTTP exercises the full stack: a real HTTP server
// serving server.NewHandler(server.Config{}), hit by a Connect RPC client generated from
// proto/blog/health/v1/health.proto. It exists to prove the buf-generated
// Go client/server pipeline is wired correctly end to end.
func TestHealthServiceOverHTTP(t *testing.T) {
	srv := httptest.NewServer(server.NewHandler(server.Config{}))
	defer srv.Close()

	client := healthv1connect.NewHealthServiceClient(srv.Client(), srv.URL)

	resp, err := client.Check(context.Background(), connect.NewRequest(&healthv1.CheckRequest{}))
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if got := resp.Msg.GetStatus(); got != "ok" {
		t.Fatalf("status = %q, want %q", got, "ok")
	}
}
