package services

import (
	"context"
	"errors"

	"just-kanban/internal/contextkeys"
	"just-kanban/internal/models"
	"just-kanban/internal/repositories/interfaces"
	"just-kanban/pkg/identifier"
	"just-kanban/pkg/sqlddl"
)

var (
	taskAlreadyExistsErr      = errors.New("task with this name already exists")
	taskWithOrderNotExistsErr = errors.New("task with this order does not exist")
)

type (
	TaskService struct {
		interfaces.TaskRepository
	}
	CreateTaskData struct {
		Name        string    `json:"name" validate:"required,min=3,max=255,trimmed"`
		Description string    `json:"description" validate:"max=1000,trimmed"`
		BoardID     sqlddl.ID `json:"board_id" validate:"required"`
		AssigneeID  sqlddl.ID `json:"assignee_id"`
	}
	UpdateTaskData struct {
		Name        string            `json:"name" validate:"omitempty,min=3,max=255,trimmed"`
		Description string            `json:"description" validate:"omitempty,max=1000,trimmed"`
		AssigneeID  sqlddl.ID         `json:"assignee_id"`
		Status      models.TaskStatus `json:"status"`
	}
)

func (utd *UpdateTaskData) ToUpdateTaskModel() *models.UpdateTask {
	var model models.UpdateTask
	if utd.Name != "" {
		model.Name = &utd.Name
	}
	if utd.Description != "" {
		model.Description = &utd.Description
	}
	if utd.Status != 0 {
		model.Status = &utd.Status
	}
	if utd.AssigneeID != "" {
		model.AssigneeID = &utd.AssigneeID
	}
	return &model
}

func NewTaskService(taskRepository interfaces.TaskRepository) *TaskService {
	return &TaskService{taskRepository}
}

func (ts *TaskService) CreateTask(ctx context.Context, d *CreateTaskData) (*models.Task, error) {
	userId, userIdErr := contextkeys.GetUserId(ctx)
	if userIdErr != nil {
		return nil, userIdErr
	}
	_, notExistErr := ts.TaskRepository.FindByName(ctx, d.BoardID, d.Name)
	if notExistErr == nil {
		return nil, taskAlreadyExistsErr
	}
	id := sqlddl.ID(identifier.GenerateUUID())
	var assigneeId = userId
	if d.AssigneeID != "" {
		assigneeId = d.AssigneeID
	}
	boardTasks, boardTasksErr := ts.FindAllByBoardId(ctx, d.BoardID)
	if boardTasksErr != nil {
		return nil, boardTasksErr
	}
	order := ts.findMaxTasksOrder(boardTasks) + 1
	creationErr := ts.TaskRepository.Create(ctx, &models.Task{
		Model:       models.Model{ID: id},
		Name:        d.Name,
		Description: d.Description,
		BoardID:     d.BoardID,
		CreatorID:   userId,
		AssigneeID:  assigneeId,
		Status:      models.TaskStatusBacklog,
		Order:       order,
	})
	if creationErr != nil {
		return nil, creationErr
	}
	createdTask, searchErr := ts.TaskRepository.FindByID(ctx, id)
	return createdTask, searchErr
}

func (ts *TaskService) UpdateTask(ctx context.Context, taskId sqlddl.ID, d *UpdateTaskData) (*models.Task, error) {
	_, userIdErr := contextkeys.GetUserId(ctx)
	if userIdErr != nil {
		return nil, userIdErr
	}
	updateErr := ts.TaskRepository.Update(ctx, taskId, d.ToUpdateTaskModel())
	if updateErr != nil {
		return nil, updateErr
	}
	updatedTask, searchErr := ts.TaskRepository.FindByID(ctx, taskId)
	return updatedTask, searchErr
}

func (ts *TaskService) FindByID(ctx context.Context, id sqlddl.ID) (*models.Task, error) {
	task, searchErr := ts.TaskRepository.FindByID(ctx, id)
	return task, searchErr
}

func (ts *TaskService) FindByName(ctx context.Context, boardID sqlddl.ID, name string) (*models.Task, error) {
	task, searchErr := ts.TaskRepository.FindByName(ctx, boardID, name)
	return task, searchErr
}

func (ts *TaskService) FindByOrder(ctx context.Context, boardID sqlddl.ID, order uint) (*models.Task, error) {
	task, searchErr := ts.TaskRepository.FindByOrder(ctx, boardID, order)
	if searchErr != nil {
		return nil, taskWithOrderNotExistsErr
	}
	return task, searchErr
}

func (ts *TaskService) FindAllByBoardId(ctx context.Context, boardId sqlddl.ID) ([]models.Task, error) {
	tasks, searchErr := ts.TaskRepository.FindAllByBoardId(ctx, boardId)
	return tasks, searchErr
}

func (ts *TaskService) findMaxTasksOrder(tasks []models.Task) int {
	maxOrder := 0
	for _, task := range tasks {
		if task.Order > maxOrder {
			maxOrder = task.Order
		}
	}
	return maxOrder
}
