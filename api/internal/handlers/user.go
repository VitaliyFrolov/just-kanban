package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"just-kanban/internal/config"
	"just-kanban/internal/services"
	"just-kanban/pkg/sqlddl"
	"just-kanban/pkg/validation"
)

// UserHandler handles http requests for working with methods of services.UserService
type UserHandler struct {
	services.UserService
	*validation.Validate
}

// NewUserHandler create new instance of UserHandler
func NewUserHandler(us services.UserService, validator *validation.Validate) *UserHandler {
	return &UserHandler{us, validator}
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.PathValue(config.ParamUserID)
	ctx := r.Context()
	if userIdParam == "" {
		uh.handleMultipleUser(ctx, w, r)
	} else {
		uh.handleSingleUser(ctx, w, r, userIdParam)
	}
}

func (uh *UserHandler) handleMultipleUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := uh.UserService.ListUsers(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(users)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (uh *UserHandler) handleSingleUser(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	userIdParam string,
) {
	userId := sqlddl.ID(userIdParam)
	switch r.Method {
	case http.MethodGet:
		users, err := uh.UserService.FindByID(ctx, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(users)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		var updateData *services.UpdateUserData
		err := json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if validateErr := uh.Validate.Struct(updateData); validateErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validation.FormatValidationErr(validateErr))
			return
		}
		updatedUser, updateErr := uh.UpdateUser(ctx, userId, updateData)
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(updatedUser)
		if encodeErr != nil {
			http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		deleteErr := uh.DeleteUser(ctx, userId)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
