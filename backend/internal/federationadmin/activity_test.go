package federationadmin

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"testing"
	"time"

	"github.com/gowebpki/jcs"
)

func generateTestKeyPair(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	pkcs1 := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	return key, string(pkcs1)
}

func TestBuildNote(t *testing.T) {
	publishedAt := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	n := buildNote("https://blog.nagutabby.uk", "my-article", "タイトル", publishedAt)

	if n.ID != "https://blog.nagutabby.uk/api/articles/my-article" {
		t.Fatalf("ID = %q", n.ID)
	}
	if n.Type != "Note" {
		t.Fatalf("Type = %q", n.Type)
	}
	if n.AttributedTo != "https://blog.nagutabby.uk/actor" {
		t.Fatalf("AttributedTo = %q", n.AttributedTo)
	}
	if n.Name != "タイトル" {
		t.Fatalf("Name = %q", n.Name)
	}
	if n.Published != "2025-06-15T00:00:00.000Z" {
		t.Fatalf("Published = %q", n.Published)
	}
	if n.URL != n.ID {
		t.Fatalf("URL = %q, want equal to ID", n.URL)
	}
	if len(n.To) != 1 || n.To[0] != "https://www.w3.org/ns/activitystreams#Public" {
		t.Fatalf("To = %v", n.To)
	}
}

func TestBuildDeleteNoteIsMinimal(t *testing.T) {
	n := buildDeleteNote("https://blog.nagutabby.uk", "my-article")

	if n.ID != "https://blog.nagutabby.uk/api/articles/my-article" {
		t.Fatalf("ID = %q", n.ID)
	}
	if n.Type != "Note" {
		t.Fatalf("Type = %q", n.Type)
	}
	if n.AttributedTo != "" || n.Name != "" || n.Content != "" || n.Published != "" || n.URL != "" || n.To != nil {
		t.Fatalf("delete note should only have id/type set, got %+v", n)
	}

	// Confirm the JSON actually omits the empty fields, matching web's
	// {id, type: "Note"} literal.
	raw, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(decoded) != 2 {
		t.Fatalf("expected exactly {id, type}, got %v", decoded)
	}
}

func TestBuildAndSignActivityCreate(t *testing.T) {
	key, privateKeyPEM := generateTestKeyPair(t)

	obj := buildNote("https://blog.nagutabby.uk", "my-article", "タイトル", time.Now())
	act, err := buildAndSignActivity("Create", "https://blog.nagutabby.uk/api/articles/my-article/create", "https://blog.nagutabby.uk", obj, privateKeyPEM)
	if err != nil {
		t.Fatalf("buildAndSignActivity returned error: %v", err)
	}

	if act.Type != "Create" {
		t.Fatalf("Type = %q", act.Type)
	}
	if act.Actor != "https://blog.nagutabby.uk/actor" {
		t.Fatalf("Actor = %q", act.Actor)
	}
	if act.Signature.Type != "RsaSignature2017" {
		t.Fatalf("Signature.Type = %q", act.Signature.Type)
	}
	if act.Signature.Creator != "https://blog.nagutabby.uk/actor#main-key" {
		t.Fatalf("Signature.Creator = %q", act.Signature.Creator)
	}

	verifyLDSignature(t, key, act)
}

func TestBuildAndSignActivityDelete(t *testing.T) {
	_, privateKeyPEM := generateTestKeyPair(t)

	obj := buildDeleteNote("https://blog.nagutabby.uk", "my-article")
	act, err := buildAndSignActivity("Delete", "https://blog.nagutabby.uk/api/articles/my-article/delete", "https://blog.nagutabby.uk", obj, privateKeyPEM)
	if err != nil {
		t.Fatalf("buildAndSignActivity returned error: %v", err)
	}
	if act.Type != "Delete" {
		t.Fatalf("Type = %q", act.Type)
	}
	if act.Object.Name != "" {
		t.Fatalf("expected the delete object to have no Name, got %q", act.Object.Name)
	}
}

// verifyLDSignature independently reconstructs the JCS-canonicalized
// {@context,type,actor,object,created} payload and verifies it against
// the RsaSignature2017 block using the known-good public key, proving the
// signature is not just present but cryptographically correct.
func verifyLDSignature(t *testing.T, key *rsa.PrivateKey, act activity) {
	t.Helper()

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
		Created: act.Signature.Created,
	}

	raw, err := json.Marshal(dataToSign)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	canonical, err := jcs.Transform(raw)
	if err != nil {
		t.Fatalf("failed to canonicalize: %v", err)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(act.Signature.SignatureValue)
	if err != nil {
		t.Fatalf("failed to decode signature: %v", err)
	}

	hashed := sha256.Sum256(canonical)
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hashed[:], sigBytes); err != nil {
		t.Fatalf("LD-Signature does not verify: %v", err)
	}
}
