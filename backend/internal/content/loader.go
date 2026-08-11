package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Article mirrors web's Article type (src/lib/types/blog.ts), minus HTML
// rendering: Body is the raw Markdown source. Rendering it to HTML stays a
// web-side concern (marked + KaTeX + GFM heading IDs).
type Article struct {
	ID          string
	Title       string
	Image       string
	Body        string
	PublishedAt time.Time
	UpdatedAt   time.Time
}

// Review mirrors web's Review type (src/lib/types/blog.ts).
type Review struct {
	ID          string
	Title       string
	Description string
	JPECode     string
	Image       string
	Rating      int32
	Body        string
	PublishedAt time.Time
	UpdatedAt   time.Time
}

var ErrNotFound = fmt.Errorf("content not found")

// Loader reads Markdown + frontmatter from a directory laid out like
// web/static/content: <baseDir>/articles/*.md and <baseDir>/reviews/*.md.
type Loader struct {
	baseDir string
}

func NewLoader(baseDir string) *Loader {
	return &Loader{baseDir: baseDir}
}

func (l *Loader) ListArticles() ([]Article, error) {
	entries, err := l.readMarkdownDir("articles")
	if err != nil {
		return nil, err
	}

	articles := make([]Article, 0, len(entries))
	for _, e := range entries {
		articles = append(articles, articleFromFrontMatter(e.id, e.data, e.body))
	}

	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].PublishedAt.After(articles[j].PublishedAt)
	})

	return articles, nil
}

func (l *Loader) GetArticle(id string) (Article, error) {
	data, body, err := l.readMarkdownFile("articles", id)
	if err != nil {
		return Article{}, err
	}
	return articleFromFrontMatter(id, data, body), nil
}

func (l *Loader) ListReviews() ([]Review, error) {
	entries, err := l.readMarkdownDir("reviews")
	if err != nil {
		return nil, err
	}

	reviews := make([]Review, 0, len(entries))
	for _, e := range entries {
		reviews = append(reviews, reviewFromFrontMatter(e.id, e.data, e.body))
	}

	sort.SliceStable(reviews, func(i, j int) bool {
		return reviews[i].PublishedAt.After(reviews[j].PublishedAt)
	})

	return reviews, nil
}

func (l *Loader) GetReview(id string) (Review, error) {
	data, body, err := l.readMarkdownFile("reviews", id)
	if err != nil {
		return Review{}, err
	}
	return reviewFromFrontMatter(id, data, body), nil
}

type markdownEntry struct {
	id   string
	data map[string]any
	body string
}

func (l *Loader) readMarkdownDir(contentType string) ([]markdownEntry, error) {
	dir := filepath.Join(l.baseDir, contentType)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s directory: %w", contentType, err)
	}

	var entries []markdownEntry
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		id := strings.TrimSuffix(f.Name(), ".md")
		data, body, err := l.readMarkdownFile(contentType, id)
		if err != nil {
			return nil, err
		}
		entries = append(entries, markdownEntry{id: id, data: data, body: body})
	}

	return entries, nil
}

func (l *Loader) readMarkdownFile(contentType, id string) (map[string]any, string, error) {
	path := filepath.Join(l.baseDir, contentType, id+".md")

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("reading %s: %w", path, err)
	}

	return parseFrontMatter(raw)
}

func articleFromFrontMatter(id string, data map[string]any, body string) Article {
	return Article{
		ID:          id,
		Title:       stringField(data, "title"),
		Image:       transformImagePath(stringField(data, "image"), "articles"),
		Body:        body,
		PublishedAt: timeField(data, "publishedAt"),
		UpdatedAt:   timeField(data, "updatedAt"),
	}
}

func reviewFromFrontMatter(id string, data map[string]any, body string) Review {
	return Review{
		ID:          id,
		Title:       stringField(data, "title"),
		Description: stringField(data, "description"),
		JPECode:     stringField(data, "jp_e_code"),
		Image:       transformImagePath(stringField(data, "image"), "reviews"),
		Rating:      int32Field(data, "rating"),
		Body:        body,
		PublishedAt: timeField(data, "publishedAt"),
		UpdatedAt:   timeField(data, "updatedAt"),
	}
}

// transformImagePath mirrors web's $lib/markdown.ts transformImagePath: a
// frontmatter image path like "images/foo.png" is rewritten to the URL the
// web static file server exposes it at.
func transformImagePath(imagePath, contentType string) string {
	if strings.HasPrefix(imagePath, "images/") {
		return fmt.Sprintf("/content/%s/images/%s", contentType, filepath.Base(imagePath))
	}
	return imagePath
}
