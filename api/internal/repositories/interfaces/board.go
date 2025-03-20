package interfaces

import (
	"context"

	"just-kanban/internal/models"
	"just-kanban/pkg/sqlddl"
)

// BoardRepository is an abstract data storage of  boards
type BoardRepository interface {
	// Create adds new board to storage
	Create(ctx context.Context, board *models.Board) error
	// Update changes data of board exists into storage
	Update(ctx context.Context, id sqlddl.ID, d *models.UpdateBoard) error
	// FindByID searches board with provided id
	FindByID(ctx context.Context, id sqlddl.ID) (*models.Board, error)
	// FindAll searches all existing boards
	FindAll(ctx context.Context) ([]models.Board, error)
	// Delete removes board data from storage
	Delete(ctx context.Context, id sqlddl.ID) error
}
