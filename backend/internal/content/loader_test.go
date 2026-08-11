package content

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, dir, contentType, name, contents string) {
	t.Helper()
	full := filepath.Join(dir, contentType)
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("failed to create fixture dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(full, name), []byte(contents), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}
}

func newFixtureLoader(t *testing.T) (*Loader, string) {
	t.Helper()
	dir := t.TempDir()

	writeFixture(t, dir, "articles", "older-article.md", `---
title: 古い記事
image: images/old.png
publishedAt: 2024-01-01
updatedAt: 2024-01-02
---
# 古い本文
`)
	writeFixture(t, dir, "articles", "newer-article.md", `---
title: 新しい記事
image: images/new.png
publishedAt: 2025-06-15
updatedAt: 2025-06-16
---
# 新しい本文
`)

	writeFixture(t, dir, "reviews", "some-book.md", `---
title: ある本のレビュー
description: あらすじ
jp_e_code: "1234567890123"
image: images/book.jpg
rating: 5
publishedAt: 2025-03-01
updatedAt: 2025-03-02
---
## 概要
本文です。
`)

	return NewLoader(dir), dir
}

func TestListArticlesSortedByPublishedAtDescending(t *testing.T) {
	loader, _ := newFixtureLoader(t)

	articles, err := loader.ListArticles()
	if err != nil {
		t.Fatalf("ListArticles returned error: %v", err)
	}
	if len(articles) != 2 {
		t.Fatalf("len(articles) = %d, want 2", len(articles))
	}
	if articles[0].ID != "newer-article" || articles[1].ID != "older-article" {
		t.Fatalf("articles not sorted by PublishedAt desc: got [%s, %s]", articles[0].ID, articles[1].ID)
	}
	if articles[0].Image != "/content/articles/images/new.png" {
		t.Fatalf("Image = %q, want transformed path", articles[0].Image)
	}
	if articles[0].Body != "# 新しい本文\n" {
		t.Fatalf("Body = %q", articles[0].Body)
	}
}

func TestGetArticle(t *testing.T) {
	loader, _ := newFixtureLoader(t)

	article, err := loader.GetArticle("older-article")
	if err != nil {
		t.Fatalf("GetArticle returned error: %v", err)
	}
	if article.Title != "古い記事" {
		t.Fatalf("Title = %q, want %q", article.Title, "古い記事")
	}
	if article.PublishedAt.Format("2006-01-02") != "2024-01-01" {
		t.Fatalf("PublishedAt = %v", article.PublishedAt)
	}
}

func TestGetArticleNotFound(t *testing.T) {
	loader, _ := newFixtureLoader(t)

	_, err := loader.GetArticle("does-not-exist")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestGetReview(t *testing.T) {
	loader, _ := newFixtureLoader(t)

	review, err := loader.GetReview("some-book")
	if err != nil {
		t.Fatalf("GetReview returned error: %v", err)
	}
	if review.JPECode != "1234567890123" {
		t.Fatalf("JPECode = %q", review.JPECode)
	}
	if review.Rating != 5 {
		t.Fatalf("Rating = %d, want 5", review.Rating)
	}
	if review.Image != "/content/reviews/images/book.jpg" {
		t.Fatalf("Image = %q, want transformed path", review.Image)
	}
}

func TestListReviews(t *testing.T) {
	loader, _ := newFixtureLoader(t)

	reviews, err := loader.ListReviews()
	if err != nil {
		t.Fatalf("ListReviews returned error: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("len(reviews) = %d, want 1", len(reviews))
	}
}
