package services

import (
	"context"
	"errors"

	"just-kanban/internal/access"
	"just-kanban/internal/models"
	"just-kanban/internal/repositories/interfaces"
	"just-kanban/pkg/identifier"
	"just-kanban/pkg/sqlddl"
)

type (
	BoardMemberService struct {
		interfaces.BoardMemberRepository
		*BoardService
		UserService
	}
	CreateBoardMemberData struct {
		UserId sqlddl.ID   `json:"user_id" validate:"required"`
		Role   access.Role `json:"role" validate:"required,oneof=owner manager regular"`
	}
	UpdateBoardMemberData struct {
		Role access.Role `json:"role" validate:"required,oneof=owner manager regular"`
	}
)

var (
	memberAlreadyExistsErr = errors.New("member already exists")
	noMemberExistsErr      = errors.New("member does not exist")
)

func NewBoardMemberService(repo interfaces.BoardMemberRepository, bs *BoardService, us UserService) *BoardMemberService {
	return &BoardMemberService{repo, bs, us}
}

// CreateBoardMember adds new member to board, checked before it's possible at all
func (bms *BoardMemberService) CreateBoardMember(ctx context.Context, boardId sqlddl.ID, d *CreateBoardMemberData) (*models.BoardMember, error) {
	if _, findBoardErr := bms.BoardService.FindBoardByID(ctx, boardId); findBoardErr != nil {
		return nil, findBoardErr
	}
	if _, userFindErr := bms.UserService.FindByID(ctx, d.UserId); userFindErr != nil {
		return nil, userFindErr
	}
	if _, findMemberErr := bms.FindBoardMemberByUserID(ctx, boardId, d.UserId); findMemberErr == nil {
		return nil, memberAlreadyExistsErr
	}
	id := sqlddl.ID(identifier.GenerateUUID())
	creationErr := bms.BoardMemberRepository.Create(ctx, &models.BoardMember{
		Model:   models.Model{ID: id},
		BoardID: boardId,
		UserID:  d.UserId,
		Role:    d.Role,
	})
	if creationErr != nil {
		return nil, creationErr
	}
	newBoardMember, searchErr := bms.BoardMemberRepository.FindByID(ctx, id)
	return newBoardMember, searchErr
}

func (bms *BoardMemberService) ChangeBoardMemberRole(ctx context.Context, memberId sqlddl.ID, role access.Role) (*models.BoardMember, error) {
	_, findMemberErr := bms.FindBoardMemberByID(ctx, memberId)
	if findMemberErr != nil {
		return nil, findMemberErr
	}
	updateErr := bms.BoardMemberRepository.ChangeMemberRole(ctx, memberId, role)
	if updateErr != nil {
		return nil, updateErr
	}
	updatedMember, searchErr := bms.BoardMemberRepository.FindByID(ctx, memberId)
	return updatedMember, searchErr
}

// RemoveBoardMember removes member from board
func (bms *BoardMemberService) RemoveBoardMember(ctx context.Context, memberId sqlddl.ID) error {
	_, findMemberErr := bms.FindBoardMemberByID(ctx, memberId)
	if findMemberErr != nil {
		return findMemberErr
	}
	removeErr := bms.BoardRepository.Delete(ctx, memberId)
	return removeErr
}

func (bms *BoardMemberService) FindBoardMemberByID(ctx context.Context, memberId sqlddl.ID) (*models.BoardMember, error) {
	findMember, findMemberErr := bms.BoardMemberRepository.FindByID(ctx, memberId)
	if findMemberErr != nil {
		return nil, noMemberExistsErr
	}
	return findMember, nil
}

func (bms *BoardMemberService) FindBoardMemberByUserID(ctx context.Context, boardId, userId sqlddl.ID) (*models.BoardMember, error) {
	findMember, findMemberErr := bms.BoardMemberRepository.FindBoardUser(ctx, boardId, userId)
	if findMemberErr != nil {
		return nil, noMemberExistsErr
	}
	return findMember, nil
}

func (bms *BoardMemberService) ListBoardMembers(ctx context.Context, boardId sqlddl.ID) ([]models.BoardMember, error) {
	findMembers, searchErr := bms.BoardMemberRepository.FindBoardMembers(ctx, boardId)
	if searchErr != nil {
		return nil, noMemberExistsErr
	}
	return findMembers, nil
}

func (bms *BoardMemberService) IsUserBoardOwner(ctx context.Context, userId, boardId sqlddl.ID) bool {
	findMember, findMemberErr := bms.FindBoardMemberByUserID(ctx, boardId, userId)
	if findMemberErr != nil {
		return false
	}
	return findMember.Role == access.RoleOwner
}

func (bms *BoardMemberService) IsUserBoardManager(ctx context.Context, userId, boardId sqlddl.ID) bool {
	findMember, findMemberErr := bms.FindBoardMemberByUserID(ctx, boardId, userId)
	if findMemberErr != nil {
		return false
	}
	return findMember.Role == access.RoleManager
}

// IsUserAllowedManageBoard checks if user (requester) has access to create members of board
func (bms *BoardMemberService) IsUserAllowedManageBoard(ctx context.Context, userId, boardId sqlddl.ID) bool {
	findMember, findMemberErr := bms.FindBoardMemberByUserID(ctx, boardId, userId)
	if findMemberErr != nil {
		return false
	}
	return findMember.Role == access.RoleManager || findMember.Role == access.RoleOwner
}

func (bms *BoardMemberService) IsUserAllowedDeleteMember(ctx context.Context, userId, memberId sqlddl.ID) bool {
	findMember, findMemberErr := bms.FindBoardMemberByID(ctx, memberId)
	if findMemberErr != nil {
		return false
	}
	return bms.IsUserAllowedManageBoard(ctx, userId, findMember.BoardID) || findMember.UserID == userId
}
