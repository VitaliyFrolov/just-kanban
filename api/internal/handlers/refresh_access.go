package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"just-kanban/internal/services"
	"just-kanban/pkg/auth/jwt"
)

var noRefreshTokenErr = errors.New("no refresh token")

// RefreshAccessHandler handles http requests for working with methods of services.TokenService
type RefreshAccessHandler struct {
	*services.TokenService
	*services.AuthService
}

func (ha *RefreshAccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		var accessToken string
		accessTokenErr := json.NewDecoder(r.Body).Decode(&accessToken)
		if accessTokenErr != nil {
			http.Error(w, accessTokenErr.Error(), http.StatusBadRequest)
			return
		}
		refreshToken, refreshTokenErr := r.Cookie(jwt.RefreshTokenKey)
		if refreshTokenErr != nil {
			http.Error(w, noRefreshTokenErr.Error(), http.StatusBadRequest)
			return
		}
		userId, userIdErr := ha.FindUserIDByRefresh(ctx, refreshToken.Value)
		if userIdErr != nil {
			http.Error(w, userIdErr.Error(), http.StatusBadRequest)
			return
		}
		findUser, findUserErr := ha.AuthService.UserService.FindByID(ctx, userId)
		if findUserErr != nil {
			http.Error(w, findUserErr.Error(), http.StatusBadRequest)
			return
		}
		accessPair, accessPairErr := ha.CreateAuthPair(ctx, findUser)
		if accessPairErr != nil {
			http.Error(w, accessPairErr.Error(), http.StatusBadRequest)
			return
		}
		encodeErr := json.NewEncoder(w).Encode(accessPair.AccessToken)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
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
			},
		)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
