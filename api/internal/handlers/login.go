package handlers

import (
	"encoding/json"
	"net/http"

	"just-kanban/internal/services"
	"just-kanban/pkg/auth/jwt"
	"just-kanban/pkg/validation"
)

// LoginHandler handles http requests for working with login methods of services.AuthService
type LoginHandler struct {
	*validation.Validate
	*services.AuthService
}

// NewLoginHandler creates new instance of LoginHandler
func NewLoginHandler(as *services.AuthService, validate *validation.Validate) *LoginHandler {
	return &LoginHandler{AuthService: as, Validate: validate}
}

func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var loginData services.LoginData
		decodeErr := json.NewDecoder(r.Body).Decode(&loginData)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		if validationErr := lh.Validate.Struct(&loginData); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validationErr))
			return
		}
		tokens, loginErr := lh.Login(r.Context(), &loginData)
		if loginErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.ErrorHTTPResponse{
				Fields: map[string]string{
					"root": loginErr.Error(),
				},
			})
			return
		}
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     jwt.RefreshTokenKey,
				Value:    tokens.RefreshToken,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			},
		)
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(tokens.AccessToken)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
