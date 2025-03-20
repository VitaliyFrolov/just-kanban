package interfaces

import (
	"context"

	"just-kanban/internal/models"
	"just-kanban/pkg/sqlddl"
)

// TaskRepository is an abstract storage of project's tasks
type TaskRepository interface {
	// Create adds new task record to data storage
	Create(ctx context.Context, task *models.Task) error
	// Update changes task record into data storage, where id equal provided
	Update(ctx context.Context, taskId sqlddl.ID, d *models.UpdateTask) error
	// Delete removes task record from data storage
	Delete(ctx context.Context, taskId sqlddl.ID) error
	// FindByID searches for task by provided id
	FindByID(ctx context.Context, taskId sqlddl.ID) (*models.Task, error)
	// FindByOrder searches for task by order on project board
	FindByOrder(ctx context.Context, boardId sqlddl.ID, order uint) (*models.Task, error)
	// FindByName searches for task by name on project board
	FindByName(ctx context.Context, boardId sqlddl.ID, name string) (*models.Task, error)
	// FindAllByBoardId searches for all project board's tasks
	FindAllByBoardId(ctx context.Context, boardId sqlddl.ID) ([]models.Task, error)
}
