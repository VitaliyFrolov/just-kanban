package middlewares

import (
	"fmt"
	"net/http"

	"just-kanban/pkg/auth"
	"just-kanban/pkg/cors"
	"just-kanban/pkg/router"
	"just-kanban/pkg/tcp"
)

// CORS enables cors functionality for app routing
func CORS(next http.Handler, allowedMethods map[string][]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cors.HeaderAllowOrigin, "*")
		w.Header().Set(cors.HeaderAllowCredentials, "true")
		if r.Method == http.MethodOptions {
			for pattern, methods := range allowedMethods {
				regex := router.PatternToRegex(pattern)
				if regex.MatchString(r.URL.Path) {
					cors.SetHeaderAllowedMethods(w, methods...)
					break
				}
			}
			w.Header().Set(
				cors.HeaderAllowHeaders,
				fmt.Sprintf("%v, %v", tcp.HeaderContentType, auth.TokenHeader),
			)
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
