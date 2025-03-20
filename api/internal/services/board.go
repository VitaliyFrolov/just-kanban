package services

import (
	"context"
	"errors"

	"just-kanban/internal/models"
	"just-kanban/internal/repositories/interfaces"
	"just-kanban/pkg/identifier"
	"just-kanban/pkg/sqlddl"
)

type (
	BoardService struct {
		interfaces.BoardRepository
		*TaskService
	}

	CreateBoardData struct {
		Name        string `json:"name" validate:"required,min=3,max=255,trimmed"`
		Description string `json:"description" validate:"max=1000,trimmed"`
	}

	UpdateBoardData struct {
		Name        string `json:"name" validate:"omitempty,min=3,max=255,trimmed"`
		Description string `json:"description" validate:"max=1000,trimmed"`
	}
)

var (
	boardNotExistErr = errors.New("board doesn't exist")
)

func NewBoardService(boardRepo interfaces.BoardRepository, taskService *TaskService) *BoardService {
	return &BoardService{boardRepo, taskService}
}

func (bs *BoardService) CreateBoard(ctx context.Context, d *CreateBoardData) (*models.Board, error) {
	id := sqlddl.ID(identifier.GenerateUUID())
	creationErr := bs.BoardRepository.Create(ctx, &models.Board{
		Model:       models.Model{ID: id},
		Name:        d.Name,
		Description: d.Description,
	})
	if creationErr != nil {
		return nil, creationErr
	}
	newBoard, searchErr := bs.BoardRepository.FindByID(ctx, id)
	return newBoard, searchErr
}

func (bs *BoardService) FindBoardByID(ctx context.Context, id sqlddl.ID) (*models.Board, error) {
	findBoard, searchErr := bs.BoardRepository.FindByID(ctx, id)
	if searchErr != nil {
		return nil, searchErr
	}
	return findBoard, nil
}

func (bs *BoardService) UpdateBoard(ctx context.Context, id sqlddl.ID, d *UpdateBoardData) (*models.Board, error) {
	_, searchErr := bs.BoardRepository.FindByID(ctx, id)
	if searchErr != nil {
		return nil, boardNotExistErr
	}
	updateErr := bs.BoardRepository.Update(ctx, id, &models.UpdateBoard{
		Name:        &d.Name,
		Description: &d.Description,
	})
	if updateErr != nil {
		return nil, updateErr
	}
	updatedBoard, searchErr := bs.BoardRepository.FindByID(ctx, id)
	return updatedBoard, searchErr
}

func (bs *BoardService) DeleteBoard(ctx context.Context, boardId sqlddl.ID) error {
	_, searchErr := bs.BoardRepository.FindByID(ctx, boardId)
	if searchErr != nil {
		return searchErr
	}
	deleteErr := bs.BoardRepository.Delete(ctx, boardId)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (bs *BoardService) FindAllBoards(ctx context.Context) ([]models.Board, error) {
	boards, searchErr := bs.BoardRepository.FindAll(ctx)
	if searchErr != nil {
		return nil, searchErr
	}
	return boards, nil
}
