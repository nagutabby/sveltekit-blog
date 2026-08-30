package content

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
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
		articles = append(articles, articleFromFrontMatter(e))
	}

	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].PublishedAt.After(articles[j].PublishedAt)
	})

	return articles, nil
}

func (l *Loader) GetArticle(id string) (Article, error) {
	entries, err := l.readMarkdownDir("articles")
	if err != nil {
		return Article{}, err
	}
	for _, e := range entries {
		if e.id == id {
			return articleFromFrontMatter(e), nil
		}
	}
	return Article{}, ErrNotFound
}

func (l *Loader) ListReviews() ([]Review, error) {
	entries, err := l.readMarkdownDir("reviews")
	if err != nil {
		return nil, err
	}

	reviews := make([]Review, 0, len(entries))
	for _, e := range entries {
		reviews = append(reviews, reviewFromFrontMatter(e))
	}

	sort.SliceStable(reviews, func(i, j int) bool {
		return reviews[i].PublishedAt.After(reviews[j].PublishedAt)
	})

	return reviews, nil
}

func (l *Loader) GetReview(id string) (Review, error) {
	entries, err := l.readMarkdownDir("reviews")
	if err != nil {
		return Review{}, err
	}
	for _, e := range entries {
		if e.id == id {
			return reviewFromFrontMatter(e), nil
		}
	}
	return Review{}, ErrNotFound
}

type markdownEntry struct {
	id          string
	publishedAt time.Time
	data        map[string]any
	body        string
}

// filenameDatePattern matches the ISO 8601 date-only filename convention
// (YYYY-MM-DD.md), with an optional "-N" suffix to disambiguate multiple
// posts published on the same date.
var filenameDatePattern = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})(-\d+)?\.md$`)

// parsePublishedAtFromFilename derives a post's publish date from its
// filename, which is the source of truth now that frontmatter no longer
// carries a publishedAt field.
func parsePublishedAtFromFilename(name string) (time.Time, error) {
	m := filenameDatePattern.FindStringSubmatch(name)
	if m == nil {
		return time.Time{}, fmt.Errorf("filename %q is not in YYYY-MM-DD[-N].md format", name)
	}
	return time.Parse("2006-01-02", m[1])
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
		publishedAt, err := parsePublishedAtFromFilename(f.Name())
		if err != nil {
			return nil, err
		}
		data, body, err := l.readMarkdownFile(contentType, f.Name())
		if err != nil {
			return nil, err
		}
		entries = append(entries, markdownEntry{
			id:          stringField(data, "id"),
			publishedAt: publishedAt,
			data:        data,
			body:        body,
		})
	}

	return entries, nil
}

func (l *Loader) readMarkdownFile(contentType, filename string) (map[string]any, string, error) {
	filePath := path.Join(contentType, filename)

	raw, err := fs.ReadFile(l.fsys, filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("reading %s: %w", filePath, err)
	}

	return parseFrontMatter(raw)
}

func articleFromFrontMatter(e markdownEntry) Article {
	return Article{
		ID:          e.id,
		Title:       stringField(e.data, "title"),
		Image:       transformImagePath(stringField(e.data, "image"), "articles"),
		Body:        e.body,
		PublishedAt: e.publishedAt,
		UpdatedAt:   timeField(e.data, "updatedAt"),
	}
}

func reviewFromFrontMatter(e markdownEntry) Review {
	return Review{
		ID:          e.id,
		Title:       stringField(e.data, "title"),
		Description: stringField(e.data, "description"),
		JPECode:     stringField(e.data, "jp_e_code"),
		Image:       transformImagePath(stringField(e.data, "image"), "reviews"),
		Rating:      int32Field(e.data, "rating"),
		Body:        e.body,
		PublishedAt: e.publishedAt,
		UpdatedAt:   timeField(e.data, "updatedAt"),
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
