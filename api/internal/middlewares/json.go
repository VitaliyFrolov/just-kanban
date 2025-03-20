package middlewares

import (
	"net/http"

	"just-kanban/pkg/tcp"
)

// JSONResponse proxies requests and write them headers required for json responses to client.
// If another type of response data has to be provided, then need to rewrite related headers on place
func JSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(tcp.HeaderContentType, tcp.ContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}
