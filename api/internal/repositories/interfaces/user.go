package interfaces

import (
	"context"

	"just-kanban/internal/models"
	"just-kanban/pkg/sqlddl"
)

// UserRepository is an abstract data storage of users
type UserRepository interface {
	// Create adds new user record to data storage
	Create(ctx context.Context, user *models.User) error
	// Update changes data of user record into data storage
	Update(ctx context.Context, id sqlddl.ID, d *models.UpdateUser) error
	// FindByID searches for user record by provided id
	FindByID(ctx context.Context, id sqlddl.ID) (*models.User, error)
	// FindByUsername searches for user record by provided username
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	// FindByEmail searches for user record by provided email
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	// FindAll searches for all users
	FindAll(ctx context.Context) ([]models.User, error)
	// Delete removes user record from data storage
	Delete(ctx context.Context, id sqlddl.ID) error
}
