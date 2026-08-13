package content

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
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

// Loader reads Markdown + frontmatter from an fs.FS laid out like
// web/static/content: articles/*.md and reviews/*.md. In production this
// is the backend/content package's embedded FS (see its doc comment for
// why: a plain OS directory isn't reliably present at runtime on Vercel),
// but any fs.FS works — tests pass os.DirFS pointed at a fixture
// directory instead.
type Loader struct {
	fsys fs.FS
}

func NewLoader(fsys fs.FS) *Loader {
	return &Loader{fsys: fsys}
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
	files, err := fs.ReadDir(l.fsys, contentType)
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
	filePath := path.Join(contentType, id+".md")

	raw, err := fs.ReadFile(l.fsys, filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("reading %s: %w", filePath, err)
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
		return fmt.Sprintf("/content/%s/images/%s", contentType, path.Base(imagePath))
	}
	return imagePath
}
