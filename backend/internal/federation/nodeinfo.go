package federation

import (
	"encoding/json"
	"net/http"
)

// NodeInfo (http://nodeinfo.diaspora.software/) is the standard discovery
// mechanism Fediverse tooling (Mastodon's instance directory, fedidb.org,
// federation health crawlers, ...) uses to learn what software an
// ActivityPub server runs and roughly how active it is. Without it, this
// server is invisible to that tooling even though it otherwise federates
// fine.
const (
	nodeInfoRel20 = "http://nodeinfo.diaspora.software/ns/schema/2.0"
	nodeInfoRel21 = "http://nodeinfo.diaspora.software/ns/schema/2.1"

	nodeInfo20ContentType = `application/json; profile="` + nodeInfoRel20 + `#"`
	nodeInfo21ContentType = `application/json; profile="` + nodeInfoRel21 + `#"`

	softwareName       = "sveltekit-blog"
	softwareVersion    = "1.0.0"
	softwareRepository = "https://github.com/nagutabby/sveltekit-blog"
)

// WellKnownNodeInfo serves .well-known/nodeinfo, pointing discovery clients
// at the versioned documents below.
func (h *Handlers) WellKnownNodeInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"links": []map[string]string{
			{"rel": nodeInfoRel21, "href": h.cfg.SiteBaseURL + "/nodeinfo/2.1"},
			{"rel": nodeInfoRel20, "href": h.cfg.SiteBaseURL + "/nodeinfo/2.0"},
		},
	})
}

// NodeInfo21 serves /nodeinfo/2.1.
func (h *Handlers) NodeInfo21(w http.ResponseWriter, r *http.Request) {
	h.writeNodeInfo(w, "2.1", nodeInfo21ContentType, true)
}

// NodeInfo20 serves /nodeinfo/2.0. The 2.0 schema doesn't have a
// software.repository field, unlike 2.1.
func (h *Handlers) NodeInfo20(w http.ResponseWriter, r *http.Request) {
	h.writeNodeInfo(w, "2.0", nodeInfo20ContentType, false)
}

func (h *Handlers) writeNodeInfo(w http.ResponseWriter, version, contentType string, includeRepository bool) {
	software := map[string]string{
		"name":    softwareName,
		"version": softwareVersion,
	}
	if includeRepository {
		software["repository"] = softwareRepository
	}

	localPosts := 0
	if articles, err := h.articles.ListArticles(); err == nil {
		localPosts = len(articles)
	}

	w.Header().Set("Content-Type", contentType)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"version":  version,
		"software": software,
		// This actor is a single publishing bot, not a multi-user
		// instance: ActivityPub, no other inbound/outbound protocols.
		"protocols":        []string{"activitypub"},
		"services":         map[string][]string{"inbound": {}, "outbound": {}},
		"openRegistration": false,
		"usage": map[string]any{
			"users":      map[string]int{"total": 1, "activeMonth": 1, "activeHalfyear": 1},
			"localPosts": localPosts,
		},
		"metadata": map[string]any{},
	})
}
