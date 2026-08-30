package federation

import (
	"sort"
	"strconv"
	"strings"
)

// ldJSONContentType is what a strict JSON-LD client (or one mirroring our
// own outbound Content-Type back at us) sends/expects instead of
// application/activity+json. Both describe the same ActivityStreams JSON
// payload; only the media type label differs.
const ldJSONContentType = `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`

// negotiateAPContentType picks the Content-Type an ActivityStreams JSON
// response should be served with, based on the request's Accept header.
// An empty header (no preference stated) defaults to
// application/activity+json, the media type this server has always used
// and that most ActivityPub implementations send/expect. The second
// return value is false when the client stated a preference this server
// has no representation for, so the caller can respond 406.
func negotiateAPContentType(acceptHeader string) (string, bool) {
	if strings.TrimSpace(acceptHeader) == "" {
		return activityJSONContentType, true
	}

	for _, mediaType := range parseAcceptMediaTypes(acceptHeader) {
		switch mediaType {
		case "*/*", "application/*", "application/json", "application/activity+json":
			return activityJSONContentType, true
		case "application/ld+json":
			return ldJSONContentType, true
		}
	}

	return "", false
}

// parseAcceptMediaTypes parses an Accept header into bare media types
// (parameters like charset/profile stripped, but q honored), ordered by
// descending preference. It's a deliberately small subset of RFC 9110
// content negotiation: no wildcard-subtype matching beyond "*/*" and
// "application/*", which is all the media types this server ever emits
// need.
func parseAcceptMediaTypes(header string) []string {
	type entry struct {
		mediaType string
		q         float64
	}

	var entries []entry
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		segments := strings.Split(part, ";")
		mediaType := strings.ToLower(strings.TrimSpace(segments[0]))
		if mediaType == "" {
			continue
		}

		q := 1.0
		for _, param := range segments[1:] {
			param = strings.TrimSpace(param)
			if value, ok := strings.CutPrefix(param, "q="); ok {
				if parsed, err := strconv.ParseFloat(value, 64); err == nil {
					q = parsed
				}
			}
		}
		if q <= 0 {
			continue
		}

		entries = append(entries, entry{mediaType: mediaType, q: q})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].q > entries[j].q
	})

	mediaTypes := make([]string, len(entries))
	for i, e := range entries {
		mediaTypes[i] = e.mediaType
	}
	return mediaTypes
}
