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
	"regexp"
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
		strings.ToLower(method), parsed.Path, parsed.Host, date, digest,
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

// ParseRSAPublicKey accepts both PKIX ("BEGIN PUBLIC KEY", what actor
// documents publish) and PKCS#1 ("BEGIN RSA PUBLIC KEY") PEM encodings.
func ParseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("federation: failed to decode PEM block")
	}

	if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		rsaKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("federation: public key is not RSA")
		}
		return rsaKey, nil
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("federation: parsing public key: %w", err)
	}
	return pub, nil
}

// maxSignatureAge bounds how far the signed Date header may drift from
// wall-clock time, in either direction. This blocks replay of an
// old-but-validly-signed inbox delivery long after it was first sent,
// while staying generous enough for real network/clock skew between
// federated servers.
const maxSignatureAge = 1 * time.Hour

var signatureParamPattern = regexp.MustCompile(`(\w+)="([^"]*)"`)

// parseSignatureHeader parses a draft-cavage HTTP Signature header, e.g.
// `keyId="...",algorithm="rsa-sha256",headers="(request-target) host date digest",signature="..."`.
func parseSignatureHeader(header string) (map[string]string, error) {
	matches := signatureParamPattern.FindAllStringSubmatch(header, -1)
	if matches == nil {
		return nil, errors.New("federation: malformed Signature header")
	}

	params := make(map[string]string, len(matches))
	for _, m := range matches {
		params[m[1]] = m[2]
	}

	if params["keyId"] == "" || params["signature"] == "" {
		return nil, errors.New("federation: Signature header missing keyId or signature")
	}

	return params, nil
}

// VerifyHTTPSignature checks the inbound request's draft-cavage HTTP
// Signature (the receiving side of SignHTTPRequest) against the sending
// actor's publicKeyPem, and confirms the Digest header actually matches
// body. It requires the digest to be part of the signed header set
// whenever the request carries a body, so a signature can't be replayed
// over a tampered payload.
func VerifyHTTPSignature(r *http.Request, body []byte, publicKeyPEM string) error {
	sigHeader := r.Header.Get("Signature")
	if sigHeader == "" {
		return errors.New("federation: missing Signature header")
	}

	params, err := parseSignatureHeader(sigHeader)
	if err != nil {
		return err
	}

	if algorithm := params["algorithm"]; algorithm != "" && algorithm != "rsa-sha256" {
		return fmt.Errorf("federation: unsupported signature algorithm %q", algorithm)
	}

	headerList := strings.Fields(params["headers"])
	if len(headerList) == 0 {
		headerList = []string{"date"}
	}

	dateHeader := r.Header.Get("Date")
	if dateHeader == "" {
		return errors.New("federation: missing Date header")
	}
	parsedDate, err := time.Parse(http.TimeFormat, dateHeader)
	if err != nil {
		return fmt.Errorf("federation: invalid Date header: %w", err)
	}
	if age := time.Since(parsedDate); age > maxSignatureAge || age < -maxSignatureAge {
		return errors.New("federation: Date header outside acceptable window")
	}

	if len(body) > 0 {
		if !containsHeaderFold(headerList, "digest") {
			return errors.New("federation: Digest header must be part of the signed headers")
		}
		digestHeader := r.Header.Get("Digest")
		if digestHeader == "" {
			return errors.New("federation: missing Digest header")
		}
		sum := sha256.Sum256(body)
		wantDigest := "SHA-256=" + base64.StdEncoding.EncodeToString(sum[:])
		if digestHeader != wantDigest {
			return errors.New("federation: Digest header does not match body")
		}
	}

	signingString, err := buildSigningString(r, headerList)
	if err != nil {
		return err
	}

	pubKey, err := ParseRSAPublicKey(NormalizePEM(publicKeyPEM))
	if err != nil {
		return fmt.Errorf("federation: parsing actor public key: %w", err)
	}

	sigBytes, err := base64.StdEncoding.DecodeString(params["signature"])
	if err != nil {
		return fmt.Errorf("federation: decoding signature: %w", err)
	}

	hashed := sha256.Sum256([]byte(signingString))
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sigBytes); err != nil {
		return errors.New("federation: signature verification failed")
	}

	return nil
}

func containsHeaderFold(headers []string, target string) bool {
	for _, h := range headers {
		if strings.EqualFold(h, target) {
			return true
		}
	}
	return false
}

// buildSigningString reconstructs the exact string the sender signed, per
// the "headers" param of its Signature header, from the values actually
// present on the received request.
func buildSigningString(r *http.Request, headerList []string) (string, error) {
	lines := make([]string, 0, len(headerList))
	for _, h := range headerList {
		lower := strings.ToLower(h)
		var value string
		switch lower {
		case "(request-target)":
			value = strings.ToLower(r.Method) + " " + r.RequestURI
		case "host":
			value = r.Host
		default:
			value = r.Header.Get(h)
			if value == "" {
				return "", fmt.Errorf("federation: signed header %q is missing from the request", h)
			}
		}
		lines = append(lines, lower+": "+value)
	}
	return strings.Join(lines, "\n"), nil
}
