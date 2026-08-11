package content

import (
	"context"
	"testing"

	"connectrpc.com/connect"

	contentv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/content/v1"
)

func TestServiceListArticles(t *testing.T) {
	loader, _ := newFixtureLoader(t)
	svc := NewService(loader)

	resp, err := svc.ListArticles(context.Background(), connect.NewRequest(&contentv1.ListArticlesRequest{}))
	if err != nil {
		t.Fatalf("ListArticles returned error: %v", err)
	}
	articles := resp.Msg.GetArticles()
	if len(articles) != 2 {
		t.Fatalf("len(articles) = %d, want 2", len(articles))
	}
	if articles[0].GetId() != "newer-article" {
		t.Fatalf("articles[0].Id = %q, want %q", articles[0].GetId(), "newer-article")
	}
	if articles[0].GetPublishedAt() != "2025-06-15T00:00:00Z" {
		t.Fatalf("PublishedAt = %q", articles[0].GetPublishedAt())
	}
}

func TestServiceGetArticleNotFound(t *testing.T) {
	loader, _ := newFixtureLoader(t)
	svc := NewService(loader)

	_, err := svc.GetArticle(context.Background(), connect.NewRequest(&contentv1.GetArticleRequest{Id: "missing"}))
	if err == nil {
		t.Fatal("expected an error for a missing article, got nil")
	}
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("error code = %v, want %v", connect.CodeOf(err), connect.CodeNotFound)
	}
}

func TestServiceGetReview(t *testing.T) {
	loader, _ := newFixtureLoader(t)
	svc := NewService(loader)

	resp, err := svc.GetReview(context.Background(), connect.NewRequest(&contentv1.GetReviewRequest{Id: "some-book"}))
	if err != nil {
		t.Fatalf("GetReview returned error: %v", err)
	}
	if resp.Msg.GetReview().GetRating() != 5 {
		t.Fatalf("Rating = %d, want 5", resp.Msg.GetReview().GetRating())
	}
}

func TestServiceListReviews(t *testing.T) {
	loader, _ := newFixtureLoader(t)
	svc := NewService(loader)

	resp, err := svc.ListReviews(context.Background(), connect.NewRequest(&contentv1.ListReviewsRequest{}))
	if err != nil {
		t.Fatalf("ListReviews returned error: %v", err)
	}
	if len(resp.Msg.GetReviews()) != 1 {
		t.Fatalf("len(reviews) = %d, want 1", len(resp.Msg.GetReviews()))
	}
}
