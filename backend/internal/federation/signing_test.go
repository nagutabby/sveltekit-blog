package federation

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
)

func generateTestKeyPair(t *testing.T) (*rsa.PrivateKey, string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	pkcs1 := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal PKCS8 key: %v", err)
	}
	pkcs8 := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})

	return key, string(pkcs1), string(pkcs8)
}

func TestParseRSAPrivateKeyAcceptsPKCS1AndPKCS8(t *testing.T) {
	original, pkcs1, pkcs8 := generateTestKeyPair(t)

	fromPKCS1, err := ParseRSAPrivateKey(pkcs1)
	if err != nil {
		t.Fatalf("parsing PKCS1: %v", err)
	}
	if !fromPKCS1.Equal(original) {
		t.Fatal("PKCS1-parsed key does not match original")
	}

	fromPKCS8, err := ParseRSAPrivateKey(pkcs8)
	if err != nil {
		t.Fatalf("parsing PKCS8: %v", err)
	}
	if !fromPKCS8.Equal(original) {
		t.Fatal("PKCS8-parsed key does not match original")
	}
}

func TestParseRSAPrivateKeyRejectsGarbage(t *testing.T) {
	if _, err := ParseRSAPrivateKey("not a pem"); err == nil {
		t.Fatal("expected an error for a non-PEM string, got nil")
	}
}

func TestSignHTTPRequestProducesVerifiableSignature(t *testing.T) {
	key, pkcs1, _ := generateTestKeyPair(t)
	pkcsWithEscapedNewlines := strings.ReplaceAll(pkcs1, "\n", `\n`)

	body := `{"type":"Accept"}`
	targetURL := "https://mastodon.example/users/alice/inbox"
	keyID := "https://blog.nagutabby.uk/actor#main-key"

	headers, err := SignHTTPRequest(targetURL, "POST", body, keyID, pkcsWithEscapedNewlines)
	if err != nil {
		t.Fatalf("SignHTTPRequest returned error: %v", err)
	}

	if !strings.HasPrefix(headers.Digest, "SHA-256=") {
		t.Fatalf("Digest = %q, want SHA-256= prefix", headers.Digest)
	}
	sum := sha256.Sum256([]byte(body))
	wantDigest := "SHA-256=" + base64.StdEncoding.EncodeToString(sum[:])
	if headers.Digest != wantDigest {
		t.Fatalf("Digest = %q, want %q", headers.Digest, wantDigest)
	}

	sigMatch := regexp.MustCompile(`signature="([^"]+)"`).FindStringSubmatch(headers.Signature)
	if sigMatch == nil {
		t.Fatalf("could not extract signature from header: %q", headers.Signature)
	}
	sigBytes, err := base64.StdEncoding.DecodeString(sigMatch[1])
	if err != nil {
		t.Fatalf("failed to decode signature: %v", err)
	}

	parsed, _ := url.Parse(targetURL)
	signString := "(request-target): post " + parsed.Path + "\n" +
		"host: " + parsed.Host + "\n" +
		"date: " + headers.Date + "\n" +
		"digest: " + headers.Digest
	hashed := sha256.Sum256([]byte(signString))

	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hashed[:], sigBytes); err != nil {
		t.Fatalf("signature does not verify: %v", err)
	}

	if !strings.Contains(headers.Signature, `keyId="`+keyID+`"`) {
		t.Fatalf("Signature header missing expected keyId: %q", headers.Signature)
	}
}

func encodeTestPublicKeyPEM(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

// signedInboxTestRequest builds a request the same way newSignedInboxRequest
// in handlers_test.go does, but standalone so signing_test.go can exercise
// VerifyHTTPSignature directly without going through Handlers.Inbox.
func signedInboxTestRequest(t *testing.T, body, keyID, privateKeyPEM string) *http.Request {
	t.Helper()
	const siteBaseURL = "https://blog.nagutabby.uk"
	headers, err := SignHTTPRequest(siteBaseURL+"/actor/inbox", http.MethodPost, body, keyID, privateKeyPEM)
	if err != nil {
		t.Fatalf("failed to sign request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader(body))
	req.Host = "blog.nagutabby.uk"
	req.Header.Set("Date", headers.Date)
	req.Header.Set("Digest", headers.Digest)
	req.Header.Set("Signature", headers.Signature)
	return req
}

func TestVerifyHTTPSignatureAccepts(t *testing.T) {
	key, privateKeyPEM, _ := generateTestKeyPair(t)
	publicKeyPEM := encodeTestPublicKeyPEM(t, key)

	body := `{"type":"Follow"}`
	req := signedInboxTestRequest(t, body, "https://mastodon.example/users/alice#main-key", privateKeyPEM)

	if err := VerifyHTTPSignature(req, []byte(body), publicKeyPEM); err != nil {
		t.Fatalf("VerifyHTTPSignature returned error for a validly signed request: %v", err)
	}
}

func TestVerifyHTTPSignatureRejectsMissingHeader(t *testing.T) {
	key, _, _ := generateTestKeyPair(t)
	publicKeyPEM := encodeTestPublicKeyPEM(t, key)

	req := httptest.NewRequest(http.MethodPost, "/actor/inbox", strings.NewReader("{}"))
	if err := VerifyHTTPSignature(req, []byte("{}"), publicKeyPEM); err == nil {
		t.Fatal("expected an error for a request with no Signature header")
	}
}

func TestVerifyHTTPSignatureRejectsWrongKey(t *testing.T) {
	_, signerPrivateKeyPEM, _ := generateTestKeyPair(t)
	otherKey, _, _ := generateTestKeyPair(t)
	otherPublicKeyPEM := encodeTestPublicKeyPEM(t, otherKey)

	body := `{"type":"Follow"}`
	req := signedInboxTestRequest(t, body, "https://mastodon.example/users/alice#main-key", signerPrivateKeyPEM)

	if err := VerifyHTTPSignature(req, []byte(body), otherPublicKeyPEM); err == nil {
		t.Fatal("expected an error when verifying against a public key that didn't sign the request")
	}
}

func TestVerifyHTTPSignatureRejectsBodyMismatch(t *testing.T) {
	key, privateKeyPEM, _ := generateTestKeyPair(t)
	publicKeyPEM := encodeTestPublicKeyPEM(t, key)

	signedBody := `{"type":"Follow"}`
	req := signedInboxTestRequest(t, signedBody, "https://mastodon.example/users/alice#main-key", privateKeyPEM)

	tamperedBody := `{"type":"Delete"}`
	if err := VerifyHTTPSignature(req, []byte(tamperedBody), publicKeyPEM); err == nil {
		t.Fatal("expected an error when the actual body doesn't match the signed Digest")
	}
}

func TestVerifyHTTPSignatureRejectsStaleDate(t *testing.T) {
	key, privateKeyPEM, _ := generateTestKeyPair(t)
	publicKeyPEM := encodeTestPublicKeyPEM(t, key)

	body := `{"type":"Follow"}`
	req := signedInboxTestRequest(t, body, "https://mastodon.example/users/alice#main-key", privateKeyPEM)

	staleDate := req.Header.Get("Date")
	parsed, err := http.ParseTime(staleDate)
	if err != nil {
		t.Fatalf("failed to parse Date header: %v", err)
	}
	req.Header.Set("Date", parsed.Add(-2*time.Hour).Format(http.TimeFormat))

	if err := VerifyHTTPSignature(req, []byte(body), publicKeyPEM); err == nil {
		t.Fatal("expected an error for a Date header far outside the acceptable window")
	}
}

func TestParseRSAPublicKeyRoundTrips(t *testing.T) {
	key, _, _ := generateTestKeyPair(t)
	publicKeyPEM := encodeTestPublicKeyPEM(t, key)

	parsed, err := ParseRSAPublicKey(publicKeyPEM)
	if err != nil {
		t.Fatalf("ParseRSAPublicKey returned error: %v", err)
	}
	if !parsed.Equal(&key.PublicKey) {
		t.Fatal("parsed public key does not match original")
	}
}
