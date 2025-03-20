package handlers

import (
	"go.uber.org/mock/gomock"

	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"just-kanban/internal/config"
	"just-kanban/mocks"
	"just-kanban/pkg/validation"
)

func TestUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	validator := validation.NewValidator()
	handler := NewUserHandler(mockUserService, validator)
	t.Run("No records found handling", func(t *testing.T) {
		w := httptest.NewRecorder()
		mockUserService.EXPECT().ListUsers(context.Background()).Return(
			nil,
			errors.New(""),
		)
		handler.ServeHTTP(
			w,
			httptest.NewRequest(http.MethodGet, "/users", nil),
		)
		result := w.Result()
		if result.StatusCode != http.StatusNotFound {
			t.Fatalf("got %d, expected code %d", result.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("Success read list of users", func(t *testing.T) {
		w := httptest.NewRecorder()
		mockUserService.EXPECT().ListUsers(context.Background()).Return(nil, nil)
		handler.ServeHTTP(
			w,
			httptest.NewRequest(http.MethodGet, "/users", nil),
		)
		result := w.Result()
		if result.StatusCode != http.StatusOK {
			t.Fatalf("got %d, expected code %d", result.StatusCode, http.StatusOK)
		}
	})

	t.Run("Success get user", func(t *testing.T) {
		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		req := httptest.NewRequest(http.MethodGet, "/users/uuid", nil)
		mux.HandleFunc("/users/{"+config.ParamUserID+"}", handler.ServeHTTP)
		mockUserService.EXPECT().FindByID(context.Background(), gomock.Any()).Return(nil, nil)
		mux.ServeHTTP(w, req)
		result := w.Result()
		if result.StatusCode != http.StatusOK {
			t.Fatalf("got %d, expected code %d", result.StatusCode, http.StatusOK)
		}
	})

	t.Run("Failed get user", func(t *testing.T) {
		w := httptest.NewRecorder()
		mux := http.NewServeMux()
		req := httptest.NewRequest(http.MethodGet, "/users/uuid", nil)
		mux.HandleFunc("/users/{"+config.ParamUserID+"}", handler.ServeHTTP)
		mockUserService.EXPECT().FindByID(context.Background(), gomock.Any()).Return(nil, errors.New(""))
		mux.ServeHTTP(w, req)
		result := w.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("got %d, expected code %d", result.StatusCode, http.StatusOK)
		}
	})
}
