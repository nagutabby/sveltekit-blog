package content

import (
	"testing"

	contentembed "github.com/nagutabby/sveltekit-blog/backend/content"
)

// TestRealContentParsesWithoutError guards against both frontmatter
// mistakes in actual blog posts (e.g. an unquoted title containing a
// colon breaking YAML) and the embedded FS silently missing files: it
// runs against backend/content's embed.FS (the same one production
// uses), so it exercises every article and review that ships in this
// repository through the exact path production takes.
func TestRealContentParsesWithoutError(t *testing.T) {
	loader := NewLoader(contentembed.FS)

	articles, err := loader.ListArticles()
	if err != nil {
		t.Fatalf("ListArticles returned error: %v", err)
	}
	if len(articles) == 0 {
		t.Fatal("expected at least one real article, got 0")
	}
	seenArticleIDs := make(map[string]bool)
	for _, a := range articles {
		if a.Title == "" {
			t.Errorf("article %q has an empty title", a.ID)
		}
		if a.PublishedAt.IsZero() {
			t.Errorf("article %q has no publishedAt", a.ID)
		}
		if a.ID == "" {
			t.Error("article has no id (missing frontmatter id?)")
		}
		if seenArticleIDs[a.ID] {
			t.Errorf("duplicate article id %q", a.ID)
		}
		seenArticleIDs[a.ID] = true
	}

	reviews, err := loader.ListReviews()
	if err != nil {
		t.Fatalf("ListReviews returned error: %v", err)
	}
	if len(reviews) == 0 {
		t.Fatal("expected at least one real review, got 0")
	}
	seenReviewIDs := make(map[string]bool)
	for _, r := range reviews {
		if r.Title == "" {
			t.Errorf("review %q has an empty title", r.ID)
		}
		if r.Rating < 1 || r.Rating > 5 {
			t.Errorf("review %q has an out-of-range rating: %d", r.ID, r.Rating)
		}
		if r.ID == "" {
			t.Error("review has no id (missing frontmatter id?)")
		}
		if seenReviewIDs[r.ID] {
			t.Errorf("duplicate review id %q", r.ID)
		}
		seenReviewIDs[r.ID] = true
	}
}
