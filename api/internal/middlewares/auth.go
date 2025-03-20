package middlewares

import (
	"context"
	"net/http"
	"strings"

	"just-kanban/internal/config"
	"just-kanban/internal/contextkeys"
	"just-kanban/internal/services"
	"just-kanban/pkg/auth"
	"just-kanban/pkg/auth/jwt"
	"just-kanban/pkg/sqlddl"
)

// Auth proxies request and check them on auth credentials.
// If credentials provided add id of authenticated user to request context
func Auth(next http.Handler, env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(auth.TokenHeader)
		if authHeader == "" {
			http.Error(w, auth.UnauthorizedErr.Error(), http.StatusUnauthorized)
			return
		}
		splitHeader := strings.Split(authHeader, " ")
		authType := splitHeader[0]
		accessToken := splitHeader[1]
		if authType != jwt.AuthTypeBearer || accessToken == "" {
			http.Error(w, auth.UnauthorizedErr.Error(), http.StatusUnauthorized)
			return
		}
		var accessTokenClaims services.AccessTokenClaims
		_, accessTokenParseErr := jwt.ParseWithClaims(
			&accessTokenClaims,
			accessToken,
			env.JWTSecret,
		)
		if accessTokenParseErr != nil {
			http.Error(w, auth.UnauthorizedErr.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), contextkeys.KeyUserId, sqlddl.ID(accessTokenClaims.Subject))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
