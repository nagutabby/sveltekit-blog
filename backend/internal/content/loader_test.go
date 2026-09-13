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

	writeFixture(t, dir, "articles", "2024-01-01.md", `---
id: older-article
title: 古い記事
image: images/old.png
is_draft: false
---
# 古い本文
`)
	writeFixture(t, dir, "articles", "2025-06-15.md", `---
id: newer-article
title: 新しい記事
image: images/new.png
is_draft: false
---
# 新しい本文
`)

	writeFixture(t, dir, "reviews", "2025-03-01.md", `---
id: some-book
title: ある本のレビュー
description: あらすじ
jp_e_code: "1234567890123"
image: images/book.jpg
rating: 5
is_draft: false
---
## 概要
本文です。
`)

	return NewLoader(os.DirFS(dir)), dir
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

func TestListArticlesResolvesFilenameSuffixCollisions(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "articles", "2024-05-01.md", `---
id: first-on-that-day
title: 一件目
image: images/a.png
is_draft: false
---
# 一件目本文
`)
	writeFixture(t, dir, "articles", "2024-05-01-2.md", `---
id: second-on-that-day
title: 二件目
image: images/b.png
is_draft: false
---
# 二件目本文
`)

	loader := NewLoader(os.DirFS(dir))
	articles, err := loader.ListArticles()
	if err != nil {
		t.Fatalf("ListArticles returned error: %v", err)
	}
	if len(articles) != 2 {
		t.Fatalf("len(articles) = %d, want 2", len(articles))
	}
	for _, a := range articles {
		if a.PublishedAt.Format("2006-01-02") != "2024-05-01" {
			t.Fatalf("PublishedAt = %v, want 2024-05-01", a.PublishedAt)
		}
	}
}

func TestDraftArticlesAreExcluded(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "articles", "2024-01-01.md", `---
id: published-article
title: 公開済み
image: images/a.png
is_draft: false
---
# 本文
`)
	writeFixture(t, dir, "articles", "2024-01-02.md", `---
id: draft-article
title: 下書き
image: images/b.png
is_draft: true
---
# 本文
`)
	writeFixture(t, dir, "articles", "2024-01-03.md", `---
id: missing-is-draft
title: is_draft未指定
image: images/c.png
---
# 本文
`)

	loader := NewLoader(os.DirFS(dir))

	articles, err := loader.ListArticles()
	if err != nil {
		t.Fatalf("ListArticles returned error: %v", err)
	}
	if len(articles) != 1 || articles[0].ID != "published-article" {
		t.Fatalf("ListArticles = %+v, want only published-article", articles)
	}

	if _, err := loader.GetArticle("draft-article"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetArticle(draft-article) err = %v, want ErrNotFound", err)
	}
	if _, err := loader.GetArticle("missing-is-draft"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetArticle(missing-is-draft) err = %v, want ErrNotFound (is_draft must default to true)", err)
	}
}

func TestListArticlesRejectsInvalidFilename(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "articles", "not-a-date.md", `---
id: bad-filename
title: 不正なファイル名
---
本文
`)

	loader := NewLoader(os.DirFS(dir))
	if _, err := loader.ListArticles(); err == nil {
		t.Fatal("ListArticles should return an error for a non-date filename")
	}
}
