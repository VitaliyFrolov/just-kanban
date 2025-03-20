package handlers

import (
	"net/http"

	"just-kanban/internal/contextkeys"
	"just-kanban/internal/services"
	"just-kanban/pkg/auth/jwt"
	"just-kanban/pkg/sqlddl"
)

// LogoutHandler handles http requests for working with logout methods of services.AuthService
type LogoutHandler struct {
	*services.AuthService
}

// NewLogoutHandler creates new instance of LogoutHandler
func NewLogoutHandler(as *services.AuthService) *LogoutHandler {
	return &LogoutHandler{
		AuthService: as,
	}
}

func (lh *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ctx := r.Context()
		if userID, ok := ctx.Value(contextkeys.KeyUserId).(sqlddl.ID); ok {
			removeErr := lh.AuthService.RemoveRefreshByUserID(ctx, userID)
			if removeErr != nil {
				http.Error(w, removeErr.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(
				w,
				&http.Cookie{
					Name:     jwt.RefreshTokenKey,
					Value:    "",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
					MaxAge:   -1,
				},
			)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
