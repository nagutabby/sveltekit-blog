package content

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// parseFrontMatter splits a Markdown file into its YAML frontmatter (as a
// generic map) and the remaining Markdown body. It mirrors the subset of
// gray-matter's behavior this codebase relies on: a "---" delimited YAML
// block at the very start of the file.
func parseFrontMatter(raw []byte) (map[string]any, string, error) {
	content := string(raw)

	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return map[string]any{}, content, nil
	}

	rest := strings.TrimPrefix(strings.TrimPrefix(content, "---\r\n"), "---\n")

	end := strings.Index(rest, "\n---")
	if end == -1 {
		return nil, "", fmt.Errorf("frontmatter is not terminated with \"---\"")
	}

	yamlBlock := rest[:end]
	body := rest[end+len("\n---"):]
	body = strings.TrimPrefix(strings.TrimPrefix(body, "\r\n"), "\n")

	var data map[string]any
	if err := yaml.Unmarshal([]byte(yamlBlock), &data); err != nil {
		return nil, "", fmt.Errorf("parsing frontmatter: %w", err)
	}
	if data == nil {
		data = map[string]any{}
	}

	return data, body, nil
}

func stringField(data map[string]any, key string) string {
	v, ok := data[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

func int32Field(data map[string]any, key string) int32 {
	v, ok := data[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case int:
		return int32(n)
	case int64:
		return int32(n)
	case float64:
		return int32(n)
	default:
		return 0
	}
}

// timeField reads a frontmatter date. yaml.v3 resolves bare "YYYY-MM-DD"
// scalars into time.Time when decoding into a map[string]any, but we fall
// back to string parsing to tolerate other formats.
func timeField(data map[string]any, key string) time.Time {
	v, ok := data[key]
	if !ok {
		return time.Time{}
	}
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		for _, layout := range []string{time.RFC3339, "2006-01-02"} {
			if parsed, err := time.Parse(layout, t); err == nil {
				return parsed
			}
		}
	}
	return time.Time{}
}
