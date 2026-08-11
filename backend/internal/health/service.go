package health

import (
	"context"

	"connectrpc.com/connect"

	healthv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/health/v1"
)

// Service implements the blog.health.v1.HealthService Connect RPC service.
// It exists to smoke-test the web<->backend Connect RPC pipeline; domain
// services (Contact, Content, FederationAdmin) are added in later PRs.
type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Check(
	_ context.Context,
	_ *connect.Request[healthv1.CheckRequest],
) (*connect.Response[healthv1.CheckResponse], error) {
	return connect.NewResponse(&healthv1.CheckResponse{Status: "ok"}), nil
}
