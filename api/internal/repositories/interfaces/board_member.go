package interfaces

import (
	"context"
	"just-kanban/internal/access"

	"just-kanban/internal/models"
	"just-kanban/pkg/sqlddl"
)

// BoardMemberRepository is an abstracted data storage of boards member
type BoardMemberRepository interface {
	// Create adds new board to data storage
	Create(ctx context.Context, member *models.BoardMember) error
	// ChangeMemberRole updates field of member record which contains his role info
	ChangeMemberRole(ctx context.Context, memberId sqlddl.ID, role access.Role) error
	// FindByID searches board by provided id
	FindByID(ctx context.Context, memberId sqlddl.ID) (*models.BoardMember, error)
	// FindBoardUser searches for board user by provided board and user identifiers
	FindBoardUser(ctx context.Context, boardID, userID sqlddl.ID) (*models.BoardMember, error)
	// FindBoardMembers searches for all member of a boards by provided board identifier
	FindBoardMembers(ctx context.Context, boardId sqlddl.ID) ([]models.BoardMember, error)
	// Delete removes board member from data storage
	Delete(ctx context.Context, member *models.BoardMember) error
}
