package interfaces

import (
	"context"

	"just-kanban/internal/models"
	"just-kanban/pkg/sqlddl"
)

// RefreshTokenRepository is an abstract data storage of refresh tokens
type RefreshTokenRepository interface {
	// Create add new refresh token to data storage
	Create(ctx context.Context, token *models.RefreshToken) error
	// FindByUserID searches for refresh token by its user identifier by provided identifier
	FindByUserID(ctx context.Context, id sqlddl.ID) (*models.RefreshToken, error)
	// FindUserIDByToken searches for user identifier into refresh token record by encoded token string
	FindUserIDByToken(ctx context.Context, token string) (sqlddl.ID, error)
	// FindByToken searches for refresh token record by encoded token string
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	// DeleteByToken removes refresh token record by provided encoded token string
	DeleteByToken(ctx context.Context, token string) error
	// DeleteByUserID removes refresh token record where user id equal provided one
	DeleteByUserID(ctx context.Context, userID sqlddl.ID) error
}
