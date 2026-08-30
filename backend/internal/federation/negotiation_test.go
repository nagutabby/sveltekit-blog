package federation

import "testing"

func TestNegotiateAPContentTypeDefaultsWhenNoAcceptHeader(t *testing.T) {
	contentType, ok := negotiateAPContentType("")
	if !ok {
		t.Fatal("expected ok=true for an empty Accept header")
	}
	if contentType != activityJSONContentType {
		t.Fatalf("contentType = %q, want %q", contentType, activityJSONContentType)
	}
}

func TestNegotiateAPContentTypePrefersActivityJSON(t *testing.T) {
	for _, accept := range []string{
		"application/activity+json",
		"*/*",
		"application/json",
		"text/html, application/activity+json;q=0.9",
	} {
		t.Run(accept, func(t *testing.T) {
			contentType, ok := negotiateAPContentType(accept)
			if !ok {
				t.Fatalf("expected ok=true for Accept: %q", accept)
			}
			if contentType != activityJSONContentType {
				t.Fatalf("contentType = %q, want %q", contentType, activityJSONContentType)
			}
		})
	}
}

func TestNegotiateAPContentTypeHonorsLDJSON(t *testing.T) {
	for _, accept := range []string{
		`application/ld+json; profile="https://www.w3.org/ns/activitystreams"`,
		"application/ld+json",
	} {
		t.Run(accept, func(t *testing.T) {
			contentType, ok := negotiateAPContentType(accept)
			if !ok {
				t.Fatalf("expected ok=true for Accept: %q", accept)
			}
			if contentType != ldJSONContentType {
				t.Fatalf("contentType = %q, want %q", contentType, ldJSONContentType)
			}
		})
	}
}

func TestNegotiateAPContentTypePicksHighestQValue(t *testing.T) {
	// ld+json is listed first but at a lower q, so activity+json should win.
	contentType, ok := negotiateAPContentType("application/ld+json;q=0.5, application/activity+json;q=0.9")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if contentType != activityJSONContentType {
		t.Fatalf("contentType = %q, want %q", contentType, activityJSONContentType)
	}
}

func TestNegotiateAPContentTypeRejectsIncompatibleAccept(t *testing.T) {
	_, ok := negotiateAPContentType("text/html")
	if ok {
		t.Fatal("expected ok=false when the client only accepts text/html")
	}
}

func TestNegotiateAPContentTypeIgnoresZeroQValues(t *testing.T) {
	// application/activity+json is explicitly disabled (q=0); only ld+json
	// remains viable.
	contentType, ok := negotiateAPContentType("application/activity+json;q=0, application/ld+json")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if contentType != ldJSONContentType {
		t.Fatalf("contentType = %q, want %q", contentType, ldJSONContentType)
	}
}
