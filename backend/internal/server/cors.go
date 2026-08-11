package server

import "net/http"

// withCORS allows the given origin to call a Connect RPC handler from a
// browser. Only ContactService needs this: it's the only service called
// directly from the browser across origins (blog.nagutabby.uk calling
// api.nagutabby.uk). Federation endpoints are server-to-server and
// Content/FederationAdmin are internal-only, so neither needs CORS.
func withCORS(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
