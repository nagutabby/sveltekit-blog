// Package content embeds the Markdown source under backend/content
// (articles/, reviews/) into the compiled binary, so
// internal/content.Loader can read it regardless of the runtime's working
// directory or filesystem — notably Vercel's Go builder, which compiles
// main.go into a serverless function without shipping arbitrary non-Go
// files alongside it, so a relative os.ReadFile("content/...") silently
// finds nothing in production. Mirrors backend/db's embedding of goose
// migrations (db/migrations.go) for the same reason.
package content

import "embed"

//go:embed articles reviews
var FS embed.FS
