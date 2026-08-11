package federationadmin

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gowebpki/jcs"

	"github.com/nagutabby/sveltekit-blog/backend/internal/federation"
)

type note struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	AttributedTo string   `json:"attributedTo,omitempty"`
	Name         string   `json:"name,omitempty"`
	Content      string   `json:"content,omitempty"`
	Published    string   `json:"published,omitempty"`
	URL          string   `json:"url,omitempty"`
	To           []string `json:"to,omitempty"`
}

type ldSignature struct {
	Type           string `json:"type"`
	Creator        string `json:"creator"`
	Created        string `json:"created"`
	SignatureValue string `json:"signatureValue"`
}

type activity struct {
	Context   []string    `json:"@context"`
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Actor     string      `json:"actor"`
	Published string      `json:"published"`
	To        []string    `json:"to"`
	Object    note        `json:"object"`
	Signature ldSignature `json:"signature"`
}

// buildNote mirrors the "note" object built by web's
// api/articles/[name]/{create,update}/+server.ts: attributedTo/name/
// content/published/url/to are all populated from article metadata.
func buildNote(siteBaseURL, articleID, title string, publishedAt time.Time) note {
	articleURL := fmt.Sprintf("%s/api/articles/%s", siteBaseURL, articleID)
	return note{
		ID:           articleURL,
		Type:         "Note",
		AttributedTo: siteBaseURL + "/actor",
		Name:         title,
		Content: fmt.Sprintf(
			`<p>%s</p><a href="%s/articles/%s" target="_blank">%s/articles/%s</a>`,
			title, siteBaseURL, articleID, siteBaseURL, articleID,
		),
		Published: publishedAt.UTC().Format(isoMillis) + "Z",
		URL:       articleURL,
		To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
	}
}

// buildDeleteNote mirrors web's api/articles/[name]/delete/+server.ts,
// which (unlike create/update) only ever sends {id, type: "Note"} since
// the article's metadata is typically already gone by the time it fires.
func buildDeleteNote(siteBaseURL, articleID string) note {
	return note{
		ID:   fmt.Sprintf("%s/api/articles/%s", siteBaseURL, articleID),
		Type: "Note",
	}
}

const isoMillis = "2006-01-02T15:04:05.000"

func nowISOMillis() string {
	return time.Now().UTC().Format(isoMillis) + "Z"
}

// buildAndSignActivity mirrors web's signActivity() + the Create/Update/
// Delete activity construction in api/articles/[name]/{create,update,delete}.
func buildAndSignActivity(activityType, activityID, siteBaseURL string, object note, privateKeyPEM string) (activity, error) {
	act := activity{
		Context:   []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
		ID:        activityID,
		Type:      activityType,
		Actor:     siteBaseURL + "/actor",
		Published: nowISOMillis(),
		To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
		Object:    object,
	}

	sig, err := signActivity(act, privateKeyPEM)
	if err != nil {
		return activity{}, err
	}
	act.Signature = sig

	return act, nil
}

// signActivity is the Go port of web's src/lib/signRequest.ts
// signActivity(): it JCS-canonicalizes {@context, type, actor, object,
// created} (RFC 8785, matching the npm "canonicalize" package) and signs
// the result with RSA-SHA256, producing an LDSignatures-style
// RsaSignature2017 block.
func signActivity(act activity, privateKeyPEM string) (ldSignature, error) {
	created := nowISOMillis()

	dataToSign := struct {
		Context []string `json:"@context"`
		Type    string   `json:"type"`
		Actor   string   `json:"actor"`
		Object  note     `json:"object"`
		Created string   `json:"created"`
	}{
		Context: act.Context,
		Type:    act.Type,
		Actor:   act.Actor,
		Object:  act.Object,
		Created: created,
	}

	raw, err := json.Marshal(dataToSign)
	if err != nil {
		return ldSignature{}, fmt.Errorf("federationadmin: marshaling activity for signing: %w", err)
	}

	canonical, err := jcs.Transform(raw)
	if err != nil {
		return ldSignature{}, fmt.Errorf("federationadmin: canonicalizing activity: %w", err)
	}

	key, err := federation.ParseRSAPrivateKey(federation.NormalizePEM(privateKeyPEM))
	if err != nil {
		return ldSignature{}, err
	}

	hashed := sha256.Sum256(canonical)
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return ldSignature{}, fmt.Errorf("federationadmin: signing activity: %w", err)
	}

	return ldSignature{
		Type:           "RsaSignature2017",
		Creator:        act.Actor + "#main-key",
		Created:        created,
		SignatureValue: base64.StdEncoding.EncodeToString(sigBytes),
	}, nil
}
