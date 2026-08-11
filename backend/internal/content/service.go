package content

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"

	contentv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/content/v1"
)

// Service implements the blog.content.v1.ContentService Connect RPC
// service on top of a Loader.
type Service struct {
	loader *Loader
}

func NewService(loader *Loader) *Service {
	return &Service{loader: loader}
}

func (s *Service) ListArticles(
	_ context.Context,
	_ *connect.Request[contentv1.ListArticlesRequest],
) (*connect.Response[contentv1.ListArticlesResponse], error) {
	articles, err := s.loader.ListArticles()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pb := make([]*contentv1.Article, 0, len(articles))
	for _, a := range articles {
		pb = append(pb, articleToProto(a))
	}

	return connect.NewResponse(&contentv1.ListArticlesResponse{Articles: pb}), nil
}

func (s *Service) GetArticle(
	_ context.Context,
	req *connect.Request[contentv1.GetArticleRequest],
) (*connect.Response[contentv1.GetArticleResponse], error) {
	article, err := s.loader.GetArticle(req.Msg.GetId())
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&contentv1.GetArticleResponse{Article: articleToProto(article)}), nil
}

func (s *Service) ListReviews(
	_ context.Context,
	_ *connect.Request[contentv1.ListReviewsRequest],
) (*connect.Response[contentv1.ListReviewsResponse], error) {
	reviews, err := s.loader.ListReviews()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pb := make([]*contentv1.Review, 0, len(reviews))
	for _, r := range reviews {
		pb = append(pb, reviewToProto(r))
	}

	return connect.NewResponse(&contentv1.ListReviewsResponse{Reviews: pb}), nil
}

func (s *Service) GetReview(
	_ context.Context,
	req *connect.Request[contentv1.GetReviewRequest],
) (*connect.Response[contentv1.GetReviewResponse], error) {
	review, err := s.loader.GetReview(req.Msg.GetId())
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&contentv1.GetReviewResponse{Review: reviewToProto(review)}), nil
}

func articleToProto(a Article) *contentv1.Article {
	return &contentv1.Article{
		Id:          a.ID,
		Title:       a.Title,
		Image:       a.Image,
		Body:        a.Body,
		PublishedAt: formatTime(a.PublishedAt),
		UpdatedAt:   formatTime(a.UpdatedAt),
	}
}

func reviewToProto(r Review) *contentv1.Review {
	return &contentv1.Review{
		Id:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		JpECode:     r.JPECode,
		Image:       r.Image,
		Rating:      r.Rating,
		Body:        r.Body,
		PublishedAt: formatTime(r.PublishedAt),
		UpdatedAt:   formatTime(r.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
