package sql

import (
	"context"
	"database/sql"
	"fmt"

	"just-kanban/internal/access"
	"just-kanban/internal/models"
	"just-kanban/internal/repositories"
	"just-kanban/pkg/sqlddl"
)

type BoardMemberRepository struct {
	DB *sql.DB
}

func NewBoardMemberRepository(db *sql.DB) *BoardMemberRepository {
	return &BoardMemberRepository{db}
}

func (repo *BoardMemberRepository) Create(ctx context.Context, member *models.BoardMember) error {
	const query = "INSERT INTO %s (%s, %s, %s, %s) VALUES ($1, $2, $3, $4)"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableBoardMembers,
		sqlddl.ColumnID,
		repositories.ColumnUserID,
		repositories.ColumnBoardID,
		repositories.ColumnRole,
	)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, member.ID, member.UserID, member.BoardID, member.Role)
	return execErr
}

func (repo *BoardMemberRepository) ChangeMemberRole(ctx context.Context, id sqlddl.ID, role access.Role) error {
	const query = "UPDATE %s SET %s = $1 WHERE %s = $2"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableBoardMembers,
		repositories.ColumnRole,
		sqlddl.ColumnID,
	)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, role, id)
	return execErr
}

func (repo *BoardMemberRepository) FindByID(ctx context.Context, id sqlddl.ID) (*models.BoardMember, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnUserID,
		repositories.ColumnBoardID,
		repositories.ColumnRole,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableBoardMembers,
	)
	var member models.BoardMember
	row := repo.DB.QueryRowContext(ctx, formattedQuery, id)
	scanErr := row.Scan(
		&member.ID,
		&member.UserID,
		&member.BoardID,
		&member.Role,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	return &member, scanErr
}

func (repo *BoardMemberRepository) FindBoardUser(ctx context.Context, boardID, userID sqlddl.ID) (*models.BoardMember, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s FROM %s WHERE %[2]s = $1 AND %[3]s = $2"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnUserID,
		repositories.ColumnBoardID,
		repositories.ColumnRole,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableBoardMembers,
	)
	var member models.BoardMember
	row := repo.DB.QueryRowContext(ctx, formattedQuery, userID, boardID)
	scanErr := row.Scan(
		&member.ID,
		&member.UserID,
		&member.BoardID,
		&member.Role,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	return &member, scanErr
}

func (repo *BoardMemberRepository) FindBoardMembers(ctx context.Context, boardId sqlddl.ID) ([]models.BoardMember, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s FROM %s WHERE %[3]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnUserID,
		repositories.ColumnBoardID,
		repositories.ColumnRole,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableBoardMembers,
	)
	var members []models.BoardMember
	rows, err := repo.DB.QueryContext(ctx, formattedQuery, boardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member models.BoardMember
		scanErr := rows.Scan(
			&member.ID,
			&member.UserID,
			&member.BoardID,
			&member.Role,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		members = append(members, member)
	}
	return members, nil
}

func (repo *BoardMemberRepository) Delete(ctx context.Context, member *models.BoardMember) error {
	const query = "DELETE FROM %s WHERE %s = $1"
	formattedString := fmt.Sprintf(query, repositories.TableBoardMembers, sqlddl.ColumnID)
	_, execErr := repo.DB.ExecContext(ctx, formattedString, member.UserID)
	return execErr
}
