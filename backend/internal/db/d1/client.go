// Package d1 implements db.Querier against Cloudflare D1's HTTP query API
// (https://api.cloudflare.com/client/v4/accounts/{account}/d1/database/{db}/query),
// for use outside a Cloudflare Worker where the JS Workers Binding API
// isn't available. Local development and tests use sqlc's generated
// database/sql implementation directly against a real SQLite file
// instead (see internal/db), since D1 is SQLite-compatible.
package d1

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

// Config holds the Cloudflare account/database identifiers and API token
// needed to reach D1's HTTP query API.
type Config struct {
	AccountID  string
	DatabaseID string
	APIToken   string

	// HTTPClient is used for outgoing requests; defaults to http.DefaultClient.
	HTTPClient *http.Client

	// BaseURL overrides the Cloudflare API base URL
	// (https://api.cloudflare.com/client/v4); tests point this at an
	// httptest server instead of the real Cloudflare API.
	BaseURL string
}

const defaultBaseURL = "https://api.cloudflare.com/client/v4"

// Client implements db.Querier by issuing SQL statements over D1's HTTP
// query API.
type Client struct {
	cfg        Config
	httpClient *http.Client
	baseURL    string
}

var _ db.Querier = (*Client)(nil)

func New(cfg Config) *Client {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{cfg: cfg, httpClient: httpClient, baseURL: baseURL}
}

type d1QueryRequest struct {
	SQL    string `json:"sql"`
	Params []any  `json:"params"`
}

type d1QueryResult struct {
	Results []map[string]any `json:"results"`
	Success bool             `json:"success"`
}

type d1Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type d1Response struct {
	Result  []d1QueryResult `json:"result"`
	Success bool            `json:"success"`
	Errors  []d1Error       `json:"errors"`
}

// exec runs a single SQL statement and returns its result rows, each as a
// map keyed by column name (D1's HTTP API shape).
func (c *Client) exec(ctx context.Context, stmt string, params []any) ([]map[string]any, error) {
	body, err := json.Marshal(d1QueryRequest{SQL: stmt, Params: params})
	if err != nil {
		return nil, fmt.Errorf("d1: marshal request: %w", err)
	}

	url := fmt.Sprintf(
		"%s/accounts/%s/d1/database/%s/query",
		c.baseURL, c.cfg.AccountID, c.cfg.DatabaseID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("d1: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("d1: request failed: %w", err)
	}
	defer resp.Body.Close()

	var parsed d1Response
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("d1: decode response (status %d): %w", resp.StatusCode, err)
	}
	if !parsed.Success || resp.StatusCode >= 300 {
		if len(parsed.Errors) > 0 {
			return nil, fmt.Errorf("d1: query failed: %s (code %d)", parsed.Errors[0].Message, parsed.Errors[0].Code)
		}
		return nil, fmt.Errorf("d1: query failed with status %d", resp.StatusCode)
	}
	if len(parsed.Result) == 0 {
		return nil, nil
	}
	return parsed.Result[0].Results, nil
}

// Exec runs an arbitrary SQL statement against D1 and returns its result
// rows. Exported for one-off operational tooling (e.g. cmd/migrate-to-d1)
// that needs to bypass db.Querier's upsert semantics — which force
// following/connected to true, appropriate for a live Follow/Accept but
// not for replaying historical rows that may have following=false.
func (c *Client) Exec(ctx context.Context, stmt string, params []any) ([]map[string]any, error) {
	return c.exec(ctx, stmt, params)
}

// execOne runs stmt and returns its first result row, or sql.ErrNoRows if
// it returned none, matching the :one semantics sqlc generates elsewhere.
func (c *Client) execOne(ctx context.Context, stmt string, params []any) (map[string]any, error) {
	rows, err := c.exec(ctx, stmt, params)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}
	return rows[0], nil
}

func rowString(row map[string]any, col string) string {
	s, _ := row[col].(string)
	return s
}

func rowNullString(row map[string]any, col string) sql.NullString {
	v, ok := row[col]
	if !ok || v == nil {
		return sql.NullString{}
	}
	s, ok := v.(string)
	if !ok {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func rowBool(row map[string]any, col string) bool {
	switch v := row[col].(type) {
	case bool:
		return v
	case float64:
		return v != 0
	default:
		return false
	}
}

func rowInt64(row map[string]any, col string) int64 {
	f, _ := row[col].(float64)
	return int64(f)
}

// nullStringParam converts a sql.NullString into the value D1's JSON query
// API expects for a bound parameter (a JSON string, or null).
func nullStringParam(v sql.NullString) any {
	if !v.Valid {
		return nil
	}
	return v.String
}

func scanFollower(row map[string]any) db.Follower {
	return db.Follower{
		ID:           rowInt64(row, "id"),
		ActorId:      rowString(row, "actorId"),
		Inbox:        rowString(row, "inbox"),
		PublicKeyPem: rowString(row, "publicKeyPem"),
		Following:    rowBool(row, "following"),
		CreatedAt:    rowString(row, "createdAt"),
		UpdatedAt:    rowString(row, "updatedAt"),
	}
}

func scanRelayConnection(row map[string]any) db.RelayConnection {
	return db.RelayConnection{
		ID:             rowInt64(row, "id"),
		ActorId:        rowString(row, "actorId"),
		Inbox:          rowString(row, "inbox"),
		Connected:      rowBool(row, "connected"),
		LastAcceptedAt: rowNullString(row, "lastAcceptedAt"),
		CreatedAt:      rowString(row, "createdAt"),
		UpdatedAt:      rowString(row, "updatedAt"),
	}
}

func (c *Client) CountActiveFollowers(ctx context.Context) (int64, error) {
	row, err := c.execOne(ctx, `SELECT count(*) AS count FROM "Follower" WHERE "following" = 1`, nil)
	if err != nil {
		return 0, err
	}
	return rowInt64(row, "count"), nil
}

func (c *Client) GetFollowerByActorID(ctx context.Context, actorid string) (db.Follower, error) {
	row, err := c.execOne(ctx, `SELECT id, actorId, inbox, publicKeyPem, "following", createdAt, updatedAt FROM "Follower" WHERE "actorId" = ?`, []any{actorid})
	if err != nil {
		return db.Follower{}, err
	}
	return scanFollower(row), nil
}

func (c *Client) UpsertFollower(ctx context.Context, arg db.UpsertFollowerParams) (db.Follower, error) {
	row, err := c.execOne(ctx, `
		INSERT INTO "Follower" ("actorId", "inbox", "publicKeyPem", "following", "createdAt", "updatedAt")
		VALUES (?, ?, ?, 1, ?, ?)
		ON CONFLICT ("actorId") DO UPDATE SET
		    "following" = 1,
		    "inbox" = excluded."inbox",
		    "publicKeyPem" = excluded."publicKeyPem",
		    "updatedAt" = excluded."updatedAt"
		RETURNING id, actorId, inbox, publicKeyPem, "following", createdAt, updatedAt`,
		[]any{arg.ActorId, arg.Inbox, arg.PublicKeyPem, arg.CreatedAt, arg.UpdatedAt},
	)
	if err != nil {
		return db.Follower{}, err
	}
	return scanFollower(row), nil
}

func (c *Client) UnfollowByActorID(ctx context.Context, arg db.UnfollowByActorIDParams) (db.Follower, error) {
	row, err := c.execOne(ctx, `
		UPDATE "Follower"
		SET "following" = 0,
		    "inbox" = ?,
		    "publicKeyPem" = ?,
		    "updatedAt" = ?
		WHERE "actorId" = ?
		RETURNING id, actorId, inbox, publicKeyPem, "following", createdAt, updatedAt`,
		[]any{arg.Inbox, arg.PublicKeyPem, arg.UpdatedAt, arg.ActorId},
	)
	if err != nil {
		return db.Follower{}, err
	}
	return scanFollower(row), nil
}

func (c *Client) GetRelayConnectionByActorID(ctx context.Context, actorid string) (db.RelayConnection, error) {
	row, err := c.execOne(ctx, `SELECT id, actorId, inbox, connected, lastAcceptedAt, createdAt, updatedAt FROM "RelayConnection" WHERE "actorId" = ?`, []any{actorid})
	if err != nil {
		return db.RelayConnection{}, err
	}
	return scanRelayConnection(row), nil
}

func (c *Client) ListRelayConnections(ctx context.Context) ([]db.RelayConnection, error) {
	rows, err := c.exec(ctx, `SELECT id, actorId, inbox, connected, lastAcceptedAt, createdAt, updatedAt FROM "RelayConnection" ORDER BY "id"`, nil)
	if err != nil {
		return nil, err
	}
	items := make([]db.RelayConnection, 0, len(rows))
	for _, row := range rows {
		items = append(items, scanRelayConnection(row))
	}
	return items, nil
}

func (c *Client) UpsertRelayConnectionAccepted(ctx context.Context, arg db.UpsertRelayConnectionAcceptedParams) (db.RelayConnection, error) {
	row, err := c.execOne(ctx, `
		INSERT INTO "RelayConnection" ("actorId", "inbox", "connected", "lastAcceptedAt", "createdAt", "updatedAt")
		VALUES (?, ?, 1, ?, ?, ?)
		ON CONFLICT ("actorId") DO UPDATE SET
		    "connected" = 1,
		    "inbox" = excluded."inbox",
		    "lastAcceptedAt" = excluded."lastAcceptedAt",
		    "updatedAt" = excluded."updatedAt"
		RETURNING id, actorId, inbox, connected, lastAcceptedAt, createdAt, updatedAt`,
		[]any{arg.ActorId, arg.Inbox, nullStringParam(arg.LastAcceptedAt), arg.CreatedAt, arg.UpdatedAt},
	)
	if err != nil {
		return db.RelayConnection{}, err
	}
	return scanRelayConnection(row), nil
}
