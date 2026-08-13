package federation

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nagutabby/sveltekit-blog/backend/internal/content"
	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

// Handlers implements the public ActivityPub HTTP surface previously
// served by web/src/routes/{.well-known/webfinger,actor,actor/*,
// api/articles/[name]}. These stay plain net/http (not Connect RPC)
// because external federation partners (Mastodon, relays) speak
// ActivityPub JSON-LD over HTTP, not Connect.
type Handlers struct {
	followers FollowerStore
	relays    RelayStore
	articles  ArticleStore
	cfg       Config
}

func NewHandlers(followers FollowerStore, relays RelayStore, articles ArticleStore, cfg Config) *Handlers {
	return &Handlers{followers: followers, relays: relays, articles: articles, cfg: cfg}
}

const activityJSONContentType = "application/activity+json"
const actorCacheControl = "max-age=0, private, must-revalidate"

// nowTimestamp formats the current time the way the sqlite/D1 schema
// stores it: TEXT columns holding RFC3339 with nanosecond precision.
func nowTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func (h *Handlers) Webfinger(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		writePlainText(w, http.StatusBadRequest, "Resource parameter required")
		return
	}

	expected := fmt.Sprintf("acct:article@%s", hostOf(h.cfg.SiteBaseURL))
	if resource != expected {
		writePlainText(w, http.StatusNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/jrd+json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"subject": expected,
		"links": []map[string]string{
			{"rel": "self", "type": activityJSONContentType, "href": h.cfg.actorURL()},
			{"rel": "http://webfinger.net/rel/profile-page", "type": "text/html", "href": h.cfg.SiteBaseURL},
		},
	})
}

func (h *Handlers) Actor(w http.ResponseWriter, r *http.Request) {
	writeActivityJSON(w, http.StatusOK, map[string]any{
		"@context":          []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
		"id":                h.cfg.actorURL(),
		"type":              "Service",
		"preferredUsername": "article",
		"name":              "nagutabbyの考え事",
		"summary":           `<p>ブログ記事を投稿するBotアカウントです。</p><p>運用者: <a href="https://mastodon.social/@nagutabby" target="_blank">@nagutabby</a></p>`,
		"url":               h.cfg.SiteBaseURL,
		"inbox":             h.cfg.actorURL() + "/inbox",
		"outbox":            h.cfg.actorURL() + "/outbox",
		"following":         h.cfg.actorURL() + "/following",
		"followers":         h.cfg.actorURL() + "/followers",
		"discoverable":      true,
		"publicKey": map[string]string{
			"id":           h.cfg.actorKeyID(),
			"owner":        h.cfg.actorURL(),
			"publicKeyPem": NormalizePEM(h.cfg.ActorPublicKeyPEM),
		},
		"icon": map[string]string{
			"type":      "Image",
			"mediaType": "image/png",
			"url":       h.cfg.SiteBaseURL + "/images/Microsoft-Fluentui-Emoji-3d-Cat-3d.500.png",
		},
	})
}

func (h *Handlers) Followers(w http.ResponseWriter, r *http.Request) {
	count, err := h.followers.CountActiveFollowers(r.Context())
	if err != nil {
		writePlainText(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	followersURL := h.cfg.actorURL() + "/followers"
	writeActivityJSON(w, http.StatusOK, map[string]any{
		"@context":   "https://www.w3.org/ns/activitystreams",
		"id":         followersURL,
		"type":       "OrderedCollection",
		"totalItems": count,
		"first":      followersURL + "?page=1",
		"last":       followersURL + "?page=1",
	})
}

func (h *Handlers) Following(w http.ResponseWriter, r *http.Request) {
	followingURL := h.cfg.actorURL() + "/following"
	writeActivityJSON(w, http.StatusOK, map[string]any{
		"@context":   "https://www.w3.org/ns/activitystreams",
		"id":         followingURL,
		"type":       "OrderedCollection",
		"totalItems": 0,
		"first":      followingURL + "?page=1",
		"last":       followingURL + "?page=1",
	})
}

// ArticleNote renders a bare ActivityPub Note representation of an
// article, mirroring web's api/articles/[name]/+server.ts. Activities
// built by internal/federationadmin reference this URL as the object's
// id/url, so other AP servers dereference it here.
func (h *Handlers) ArticleNote(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writePlainText(w, http.StatusNotFound, "Article not found")
		return
	}

	article, err := h.articles.GetArticle(name)
	if err != nil {
		if errors.Is(err, content.ErrNotFound) {
			writePlainText(w, http.StatusNotFound, "Article not found")
			return
		}
		writePlainText(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	articleURL := fmt.Sprintf("%s/api/articles/%s", h.cfg.SiteBaseURL, name)

	w.Header().Set("Content-Type", activityJSONContentType)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":           articleURL,
		"type":         "Note",
		"attributedTo": h.cfg.actorURL(),
		"name":         article.Title,
		"content": fmt.Sprintf(
			`<p>%s</p><a href="%s/articles/%s" target="_blank">%s/articles/%s</a>`,
			article.Title, h.cfg.SiteBaseURL, name, h.cfg.SiteBaseURL, name,
		),
		"published": article.PublishedAt.UTC().Format("2006-01-02T15:04:05.000") + "Z",
		"url":       articleURL,
		"to":        []string{"https://www.w3.org/ns/activitystreams#Public"},
	})
}

type atomFeedEntryCount struct {
	Entries []struct{} `xml:"entry"`
}

func (h *Handlers) Outbox(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, h.cfg.webBaseURL()+"/atom.xml", nil)
	if err != nil {
		writePlainText(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resp, err := h.cfg.httpClient().Do(req)
	if err != nil {
		writePlainText(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writePlainText(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var feed atomFeedEntryCount
	if err := xml.Unmarshal(body, &feed); err != nil {
		writePlainText(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	writeActivityJSON(w, http.StatusOK, map[string]any{
		"@context":   "https://www.w3.org/ns/activitystreams",
		"id":         h.cfg.actorURL() + "/outbox",
		"type":       "OrderedCollection",
		"totalItems": len(feed.Entries),
	})
}

type incomingActivity struct {
	Context json.RawMessage `json:"@context"`
	Type    string          `json:"type"`
	Actor   string          `json:"actor"`
	Object  json.RawMessage `json:"object"`
	raw     json.RawMessage
}

type activityObject struct {
	Type string `json:"type"`
}

func (h *Handlers) Inbox(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var activity incomingActivity
	if err := json.Unmarshal(body, &activity); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	activity.raw = body

	// Relay servers may only deliver the Accept that completes our own
	// Follow/Subscribe to them; they must not push arbitrary activities
	// into /actor/inbox. Previously enforced at the edge by the
	// Cloudflare Worker router (now gone); ported here.
	userAgent := r.Header.Get("User-Agent")
	if isRelayUserAgent(userAgent) && activity.Type != "Accept" {
		writeForbiddenRelayResponse(w, r.URL.Path, userAgent)
		return
	}

	if activity.Context == nil || activity.Type == "" || activity.Actor == "" {
		writeJSONError(w, http.StatusBadRequest, "Invalid activity: missing required fields")
		return
	}

	// Delete is handled before fetchActor/signature verification: by the
	// time an actor announces its own deletion, its actor document is
	// typically already gone (404/410), so fetchActor would always fail
	// and this cleanup would never run. The worst a forged Delete can do
	// is mark a follower unfollowed early, which is low-risk and
	// self-correcting (a real Follow re-adds them), so this is handled
	// as best-effort without a verified signature.
	if activity.Type == "Delete" {
		h.handleDelete(w, r.Context(), activity)
		return
	}

	actorInfo, err := fetchActor(r.Context(), h.cfg.httpClient(), activity.Actor)
	if err != nil {
		writePlainText(w, http.StatusBadRequest, "Could not fetch actor information")
		return
	}

	if err := VerifyHTTPSignature(r, body, actorInfo.PublicKey.PublicKeyPem); err != nil {
		writePlainText(w, http.StatusUnauthorized, "Invalid HTTP Signature")
		return
	}

	switch activity.Type {
	case "Follow":
		h.handleFollow(w, r.Context(), activity, actorInfo)
	case "Undo":
		h.handleUndo(w, r.Context(), activity, actorInfo)
	case "Accept":
		h.handleAccept(w, r.Context(), activity, actorInfo)
	default:
		if acknowledgedActivityTypes[activity.Type] {
			// Recognized ActivityPub activity types this actor has no
			// local state for (it's a single publishing bot, not a
			// timeline/likes UI): acknowledge receipt so senders don't
			// keep retrying delivery, without pretending to act on them.
			w.WriteHeader(http.StatusAccepted)
			return
		}
		writePlainText(w, http.StatusUnprocessableEntity, fmt.Sprintf("%s activity is not supported", activity.Type))
	}
}

// acknowledgedActivityTypes are activity types we understand well enough to
// safely no-op on (this actor keeps no timeline, likes, or shares), as
// opposed to truly unrecognized types that still get a 422.
var acknowledgedActivityTypes = map[string]bool{
	"Create":   true,
	"Update":   true,
	"Like":     true,
	"Announce": true,
	"Flag":     true,
	"Add":      true,
	"Remove":   true,
	"Block":    true,
	"Move":     true,
}

// handleDelete processes an actor announcing its own deletion (the object
// is the actor's own IRI, either bare or as a Tombstone-shaped object) by
// marking any matching follower row unfollowed. Anything else framed as a
// Delete — deleting some other object this server doesn't track — is
// acknowledged without action.
func (h *Handlers) handleDelete(w http.ResponseWriter, ctx context.Context, activity incomingActivity) {
	objectID := deletedObjectID(activity.Object)
	if objectID == "" || objectID != activity.Actor {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	existing, err := h.followers.GetFollowerByActorID(ctx, activity.Actor)
	if err != nil {
		// Not a follower we know about (or already removed); nothing to
		// clean up, but still acknowledge so the sender doesn't retry.
		w.WriteHeader(http.StatusAccepted)
		return
	}

	_, err = h.followers.UnfollowByActorID(ctx, db.UnfollowByActorIDParams{
		ActorId:      activity.Actor,
		Inbox:        existing.Inbox,
		PublicKeyPem: existing.PublicKeyPem,
		UpdatedAt:    nowTimestamp(),
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deletedObjectID extracts the target IRI from a Delete activity's object,
// which per the ActivityPub spec may be either a bare IRI string or an
// object (commonly a Tombstone) carrying an "id" field.
func deletedObjectID(object json.RawMessage) string {
	var asString string
	if err := json.Unmarshal(object, &asString); err == nil {
		return asString
	}

	var asObject struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(object, &asObject); err == nil {
		return asObject.ID
	}

	return ""
}

func (h *Handlers) handleFollow(w http.ResponseWriter, ctx context.Context, activity incomingActivity, actorInfo *remoteActor) {
	now := nowTimestamp()
	_, err := h.followers.UpsertFollower(ctx, db.UpsertFollowerParams{
		ActorId:      activity.Actor,
		Inbox:        actorInfo.Inbox,
		PublicKeyPem: actorInfo.PublicKey.PublicKeyPem,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.sendAccept(ctx, activity.raw, actorInfo.Inbox); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeActivityJSON(w, http.StatusOK, map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       h.cfg.SiteBaseURL + "/activities/" + uuid.NewString(),
		"type":     "Accept",
		"actor":    h.cfg.actorURL(),
		"object":   json.RawMessage(activity.raw),
	})
}

func (h *Handlers) handleUndo(w http.ResponseWriter, ctx context.Context, activity incomingActivity, actorInfo *remoteActor) {
	var object activityObject
	if err := json.Unmarshal(activity.Object, &object); err != nil || object.Type != "Follow" {
		writePlainText(w, http.StatusBadRequest, "Invalid Undo activity")
		return
	}

	_, err := h.followers.UnfollowByActorID(ctx, db.UnfollowByActorIDParams{
		ActorId:      activity.Actor,
		Inbox:        actorInfo.Inbox,
		PublicKeyPem: actorInfo.PublicKey.PublicKeyPem,
		UpdatedAt:    nowTimestamp(),
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.sendAccept(ctx, activity.raw, actorInfo.Inbox); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeActivityJSON(w, http.StatusOK, map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"type":     "Accept",
		"actor":    h.cfg.actorURL(),
		"object":   json.RawMessage(activity.raw),
	})
}

func (h *Handlers) handleAccept(w http.ResponseWriter, ctx context.Context, activity incomingActivity, actorInfo *remoteActor) {
	var object activityObject
	if err := json.Unmarshal(activity.Object, &object); err != nil || (object.Type != "Follow" && object.Type != "Subscribe") {
		writePlainText(w, http.StatusBadRequest, "Invalid Accept activity")
		return
	}

	now := nowTimestamp()
	_, err := h.relays.UpsertRelayConnectionAccepted(ctx, db.UpsertRelayConnectionAcceptedParams{
		ActorId:        activity.Actor,
		Inbox:          actorInfo.Inbox,
		LastAcceptedAt: sql.NullString{String: now, Valid: true},
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// sendAccept builds and delivers a signed Accept activity to targetInbox,
// mirroring web's sendAcceptActivity().
func (h *Handlers) sendAccept(ctx context.Context, activity json.RawMessage, targetInbox string) error {
	accept := map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       h.cfg.SiteBaseURL + "/activities/" + uuid.NewString(),
		"type":     "Accept",
		"actor":    h.cfg.actorURL(),
		"object":   activity,
	}

	body, err := json.Marshal(accept)
	if err != nil {
		return err
	}

	headers, err := SignHTTPRequest(targetInbox, http.MethodPost, string(body), h.cfg.actorKeyID(), h.cfg.ActorPrivateKeyPEM)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetInbox, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Date", headers.Date)
	req.Header.Set("Digest", headers.Digest)
	req.Header.Set("Signature", headers.Signature)
	req.Header.Set("Content-Type", activityJSONContentType)
	req.Header.Set("Accept", activityJSONContentType)

	resp, err := h.cfg.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("federation: accept delivery to %s failed with status %d", targetInbox, resp.StatusCode)
	}

	return nil
}

func writeActivityJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", activityJSONContentType)
	w.Header().Set("Cache-Control", actorCacheControl)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}

func writePlainText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}

func isRelayUserAgent(userAgent string) bool {
	return strings.Contains(strings.ToLower(userAgent), "relay")
}

func writeForbiddenRelayResponse(w http.ResponseWriter, path, userAgent string) {
	writeActivityJSON(w, http.StatusForbidden, map[string]any{
		"error":   "Forbidden",
		"status":  http.StatusForbidden,
		"message": "Relay server access is not permitted",
		"details": map[string]string{
			"reason":    "This server does not accept requests from relay servers",
			"path":      path,
			"userAgent": userAgent,
		},
	})
}

func hostOf(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return parsed.Host
}
