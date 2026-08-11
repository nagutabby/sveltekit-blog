// Package migrations embeds the goose SQL migrations that define the
// backend's Postgres schema, so they can be applied both by the goose CLI
// (pointed at this directory) and programmatically in tests.
package migrations

import "embed"

//go:embed migrations/*.sql
var FS embed.FS
