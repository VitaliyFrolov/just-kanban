package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"just-kanban/internal/repositories"
	"just-kanban/internal/services"
	"just-kanban/pkg/auth/jwt"
	"just-kanban/pkg/validation"
)

// RegistrationHandler handles http requests for working with registration methods of services.AuthService
type RegistrationHandler struct {
	*validation.Validate
	*services.AuthService
}

// NewRegistrationHandler creates new instance of RegistrationHandler
func NewRegistrationHandler(as *services.AuthService, validate *validation.Validate) *RegistrationHandler {
	return &RegistrationHandler{AuthService: as, Validate: validate}
}

func (rh *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var createUserData services.CreateUserData
		decodeErr := json.NewDecoder(r.Body).Decode(&createUserData)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}
		if validationErr := rh.Validate.Struct(&createUserData); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validationErr))
			return
		}
		tokens, registrationErr := rh.RegisterUser(r.Context(), &createUserData)
		if errors.Is(registrationErr, services.ErrorUserEmailTaken) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(
				validation.ErrorHTTPResponse{
					Fields: map[string]string{
						repositories.ColumnEmail: registrationErr.Error(),
					},
				})
			return
		}
		if errors.Is(registrationErr, services.ErrorUsernameTaken) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(
				validation.ErrorHTTPResponse{
					Fields: map[string]string{
						repositories.ColumnUsername: registrationErr.Error(),
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
		w.WriteHeader(http.StatusCreated)
		encodeErr := json.NewEncoder(w).Encode(tokens.AccessToken)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
