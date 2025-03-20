package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"just-kanban/internal/models"
	"just-kanban/internal/repositories"
	"just-kanban/pkg/sqlddl"
)

type BoardRepository struct {
	DB *sql.DB
}

func NewBoardRepository(db *sql.DB) *BoardRepository {
	return &BoardRepository{db}
}

func (repo *BoardRepository) Create(ctx context.Context, board *models.Board) error {
	const query = "INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableBoards,
		sqlddl.ColumnID,
		repositories.ColumnName,
		repositories.ColumnDescription,
	)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, board.ID, board.Name, board.Description)
	return execErr
}

func (repo *BoardRepository) Update(ctx context.Context, id sqlddl.ID, d *models.UpdateBoard) error {
	const query = "UPDATE %s SET %s WHERE %s = $3"
	var clauses []string
	if d.Name != nil {
		clauses = append(clauses, fmt.Sprintf("%s = $1", repositories.ColumnName))
	}
	if d.Description != nil {
		clauses = append(clauses, fmt.Sprintf("%s = $2", repositories.ColumnDescription))
	}
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableBoards,
		strings.Join(clauses, ", "),
		sqlddl.ColumnID,
	)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, d.Name, d.Description, id)
	return execErr
}

func (repo *BoardRepository) FindByID(ctx context.Context, id sqlddl.ID) (*models.Board, error) {
	const query = "SELECT %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnName,
		repositories.ColumnDescription,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableBoards,
	)
	row := repo.DB.QueryRowContext(ctx, formattedQuery, id)
	var board models.Board
	scanErr := row.Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}
	return &board, nil
}

func (repo *BoardRepository) FindAll(ctx context.Context) ([]models.Board, error) {
	const query = "SELECT %s, %s, %s, %s, %s FROM %s"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnName,
		repositories.ColumnDescription,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableBoards,
	)
	rows, rowsErr := repo.DB.QueryContext(ctx, formattedQuery)
	if rowsErr != nil {
		return nil, rowsErr
	}
	defer rows.Close()
	var boards []models.Board
	for rows.Next() {
		var board models.Board
		scanErr := rows.Scan(
			&board.ID,
			&board.Name,
			&board.Description,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (repo *BoardRepository) FindAllByUserID(ctx context.Context, userId sqlddl.ID) ([]models.Board, error) {
	const query = "SELECT %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnName,
		repositories.ColumnDescription,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableBoards,
	)
	rows, rowsErr := repo.DB.QueryContext(ctx, formattedQuery, userId)
	if rowsErr != nil {
		return nil, rowsErr
	}
	defer rows.Close()
	var boards []models.Board
	for rows.Next() {
		var board models.Board
		scanErr := rows.Scan(
			&board.ID,
			&board.Name,
			&board.Description,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (repo *BoardRepository) Delete(ctx context.Context, id sqlddl.ID) error {
	const query = "DELETE FROM %s WHERE %s = $1"
	formattedQuery := fmt.Sprintf(query, repositories.TableBoards, sqlddl.ColumnID)
	_, err := repo.DB.ExecContext(ctx, formattedQuery, id)
	return err
}
