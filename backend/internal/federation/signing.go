package federation

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ParseRSAPrivateKey accepts both PKCS#1 ("BEGIN RSA PRIVATE KEY") and
// PKCS#8 ("BEGIN PRIVATE KEY") PEM encodings, mirroring what Node's
// crypto.sign() accepts transparently.
func ParseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("federation: failed to decode PEM block")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("federation: parsing private key: %w", err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("federation: private key is not RSA")
	}
	return rsaKey, nil
}

// SignedRequestHeaders is the Go equivalent of web's signRequest() in
// src/lib/signRequest.ts: an HTTP Signature (draft-cavage) over
// (request-target), host, date and digest.
type SignedRequestHeaders struct {
	Date      string
	Digest    string
	Signature string
}

func SignHTTPRequest(targetURL, method, body, keyID, privateKeyPEM string) (SignedRequestHeaders, error) {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return SignedRequestHeaders{}, fmt.Errorf("federation: parsing target URL: %w", err)
	}

	key, err := ParseRSAPrivateKey(NormalizePEM(privateKeyPEM))
	if err != nil {
		return SignedRequestHeaders{}, err
	}

	date := time.Now().UTC().Format(http.TimeFormat)

	sum := sha256.Sum256([]byte(body))
	digest := "SHA-256=" + base64.StdEncoding.EncodeToString(sum[:])

	signString := fmt.Sprintf(
		"(request-target): %s %s\nhost: %s\ndate: %s\ndigest: %s",
		strings.ToLower(method), parsed.Path, parsed.Hostname(), date, digest,
	)

	hashed := sha256.Sum256([]byte(signString))
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	if err != nil {
		return SignedRequestHeaders{}, fmt.Errorf("federation: signing request: %w", err)
	}

	signature := fmt.Sprintf(
		`keyId="%s",algorithm="rsa-sha256",headers="(request-target) host date digest",signature="%s"`,
		keyID, base64.StdEncoding.EncodeToString(sigBytes),
	)

	return SignedRequestHeaders{Date: date, Digest: digest, Signature: signature}, nil
}
