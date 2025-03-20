package sql

import (
	"context"
	"database/sql"
	"fmt"

	"just-kanban/internal/models"
	"just-kanban/internal/repositories"
	"just-kanban/pkg/sqlddl"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db}
}

func (repo *RefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	const query = "INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableRefreshTokens,
		sqlddl.ColumnID,
		repositories.ColumnUserID,
		repositories.ColumnToken,
	)
	_, execErr := repo.DB.ExecContext(
		ctx,
		formattedQuery,
		token.ID,
		token.UserID,
		token.Token,
	)
	return execErr
}

func (repo *RefreshTokenRepository) FindByUserID(ctx context.Context, id sqlddl.ID) (*models.RefreshToken, error) {
	const query = "SELECT %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnUserID,
		repositories.ColumnToken,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableRefreshTokens,
	)
	var findToken models.RefreshToken
	row := repo.DB.QueryRowContext(ctx, formattedQuery, id)
	scanErr := row.Scan(
		&findToken.ID,
		&findToken.UserID,
		&findToken.Token,
		&findToken.CreatedAt,
		&findToken.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}
	return &findToken, nil
}

func (repo *RefreshTokenRepository) FindUserIDByToken(ctx context.Context, token string) (sqlddl.ID, error) {
	const query = "SELECT %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.TableRefreshTokens,
	)
	var userId sqlddl.ID
	row := repo.DB.QueryRowContext(ctx, formattedQuery, token)
	scanErr := row.Scan(&userId)
	if scanErr != nil {
		return "", scanErr
	}
	return userId, nil
}

func (repo *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	const query = "SELECT %s, %s, %s, %s, %s FROM %s WHERE %[3]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnUserID,
		repositories.ColumnToken,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableRefreshTokens,
	)
	var findToken models.RefreshToken
	row := repo.DB.QueryRowContext(ctx, formattedQuery, token)
	scanErr := row.Scan(
		&findToken.ID,
		&findToken.UserID,
		&findToken.Token,
		&findToken.CreatedAt,
		&findToken.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}
	return &findToken, nil
}

func (repo *RefreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	const query = "DELETE FROM %s WHERE %s = $1"
	formattedQuery := fmt.Sprintf(query, repositories.TableRefreshTokens, repositories.ColumnToken)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, token)
	return execErr
}

func (repo *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID sqlddl.ID) error {
	const query = "DELETE FROM %s WHERE %s = $1"
	formattedQuery := fmt.Sprintf(query, repositories.TableRefreshTokens, repositories.ColumnUserID)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, userID)
	return execErr
}
