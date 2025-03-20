package sql

import (
	"context"
	"database/sql"
	"fmt"

	"just-kanban/internal/models"
	"just-kanban/internal/repositories"
	"just-kanban/pkg/sqlddl"
	"just-kanban/pkg/sqlquery"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (repo *UserRepository) Create(ctx context.Context, user *models.User) error {
	const query = "INSERT INTO %s(%s, %s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableUsers,
		sqlddl.ColumnID,
		repositories.ColumnEmail,
		repositories.ColumnPassword,
		repositories.ColumnUsername,
		repositories.ColumnFirstName,
		repositories.ColumnsLastName,
		repositories.ColumnAvatar,
	)
	_, execErr := repo.DB.ExecContext(
		ctx,
		formattedQuery,
		user.ID,
		user.Email,
		user.Password,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Avatar,
	)
	return execErr
}

// Update partial change data of user record and save it to db
func (repo *UserRepository) Update(ctx context.Context, id sqlddl.ID, d *models.UpdateUser) error {
	execErr := sqlquery.DynamicUpdate(ctx, repo.DB, &sqlquery.DynamicUpdateParams{
		TableName:   repositories.TableUsers,
		WhereColumn: sqlddl.ColumnID,
		WhereValue:  id,
		Changes: map[string]interface{}{
			repositories.ColumnFirstName: d.FirstName,
			repositories.ColumnsLastName: d.LastName,
			repositories.ColumnAvatar:    d.Avatar,
		},
	})
	return execErr
}

func (repo *UserRepository) FindByID(ctx context.Context, id sqlddl.ID) (*models.User, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnEmail,
		repositories.ColumnPassword,
		repositories.ColumnAvatar,
		repositories.ColumnUsername,
		repositories.ColumnFirstName,
		repositories.ColumnsLastName,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableUsers,
	)
	var findUser models.User
	row := repo.DB.QueryRowContext(ctx, formattedQuery, id)
	scanErr := row.Scan(
		&findUser.ID,
		&findUser.Email,
		&findUser.Password,
		&findUser.Avatar,
		&findUser.Username,
		&findUser.FirstName,
		&findUser.LastName,
		&findUser.CreatedAt,
		&findUser.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}
	return &findUser, nil
}

func (repo *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.ColumnUsername,
		sqlddl.ColumnID,
		repositories.ColumnEmail,
		repositories.ColumnPassword,
		repositories.ColumnAvatar,
		repositories.ColumnFirstName,
		repositories.ColumnsLastName,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableUsers,
	)
	var findUser models.User
	row := repo.DB.QueryRowContext(ctx, formattedQuery, username)
	scanErr := row.Scan(
		&findUser.Username,
		&findUser.ID,
		&findUser.Email,
		&findUser.Password,
		&findUser.Avatar,
		&findUser.FirstName,
		&findUser.LastName,
		&findUser.CreatedAt,
		&findUser.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}
	return &findUser, nil
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	fmt.Println("SEARCH BY EMAIL")
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %[2]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnEmail,
		repositories.ColumnPassword,
		repositories.ColumnAvatar,
		repositories.ColumnUsername,
		repositories.ColumnFirstName,
		repositories.ColumnsLastName,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableUsers,
	)
	var findUser models.User
	row := repo.DB.QueryRowContext(ctx, formattedQuery, email)
	scanErr := row.Scan(
		&findUser.ID,
		&findUser.Email,
		&findUser.Password,
		&findUser.Avatar,
		&findUser.Username,
		&findUser.FirstName,
		&findUser.LastName,
		&findUser.CreatedAt,
		&findUser.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}
	return &findUser, nil
}

func (repo *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnEmail,
		repositories.ColumnPassword,
		repositories.ColumnAvatar,
		repositories.ColumnUsername,
		repositories.ColumnFirstName,
		repositories.ColumnsLastName,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableUsers,
	)
	var findUsers []models.User
	rows, rowsErr := repo.DB.QueryContext(ctx, formattedQuery)
	if rowsErr != nil {
		return nil, rowsErr
	}
	defer rows.Close()
	for rows.Next() {
		var findUser models.User
		scanErr := rows.Scan(
			&findUser.ID,
			&findUser.Email,
			&findUser.Password,
			&findUser.Avatar,
			&findUser.Username,
			&findUser.FirstName,
			&findUser.LastName,
			&findUser.CreatedAt,
			&findUser.UpdatedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		findUsers = append(findUsers, findUser)
	}
	return findUsers, nil
}

func (repo *UserRepository) Delete(ctx context.Context, id sqlddl.ID) error {
	const query = "DELETE FROM %s WHERE %s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableUsers,
		sqlddl.ColumnID,
	)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, id)
	if execErr != nil {
		return execErr
	}
	return nil
}
