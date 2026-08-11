package federation

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"

	"github.com/nagutabby/sveltekit-blog/backend/internal/db"
)

// Handlers implements the public ActivityPub HTTP surface previously
// served by web/src/routes/{.well-known/webfinger,actor,actor/*}. These
// stay plain net/http (not Connect RPC) because external federation
// partners (Mastodon, relays) speak ActivityPub JSON-LD over HTTP, not
// Connect.
type Handlers struct {
	followers FollowerStore
	relays    RelayStore
	cfg       Config
}

func NewHandlers(followers FollowerStore, relays RelayStore, cfg Config) *Handlers {
	return &Handlers{followers: followers, relays: relays, cfg: cfg}
}

const activityJSONContentType = "application/activity+json"
const actorCacheControl = "max-age=0, private, must-revalidate"

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
			"publicKeyPem": normalizePEM(h.cfg.ActorPublicKeyPEM),
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

	if activity.Context == nil || activity.Type == "" || activity.Actor == "" {
		writeJSONError(w, http.StatusBadRequest, "Invalid activity: missing required fields")
		return
	}

	actorInfo, err := fetchActor(r.Context(), h.cfg.httpClient(), activity.Actor)
	if err != nil {
		writePlainText(w, http.StatusBadRequest, "Could not fetch actor information")
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
		writePlainText(w, http.StatusUnprocessableEntity, fmt.Sprintf("%s activity is not supported", activity.Type))
	}
}

func (h *Handlers) handleFollow(w http.ResponseWriter, ctx context.Context, activity incomingActivity, actorInfo *remoteActor) {
	_, err := h.followers.UpsertFollower(ctx, db.UpsertFollowerParams{
		ActorId:      activity.Actor,
		Inbox:        actorInfo.Inbox,
		PublicKeyPem: actorInfo.PublicKey.PublicKeyPem,
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

	_, err := h.relays.UpsertRelayConnectionAccepted(ctx, db.UpsertRelayConnectionAcceptedParams{
		ActorId: activity.Actor,
		Inbox:   actorInfo.Inbox,
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

	headers, err := signHTTPRequest(targetInbox, http.MethodPost, string(body), h.cfg.actorKeyID(), h.cfg.ActorPrivateKeyPEM)
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

func hostOf(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return parsed.Host
}
