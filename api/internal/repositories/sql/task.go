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

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (repo *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	const query = "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableTasks,
		sqlddl.ColumnID,
		repositories.ColumnName,
		repositories.ColumnDescription,
		repositories.ColumnStatus,
		repositories.ColumnOrder,
		repositories.ColumnBoardID,
		repositories.ColumnCreatorID,
		repositories.ColumnAssigneeID,
	)
	_, execErr := repo.DB.ExecContext(
		ctx,
		formattedQuery,
		task.ID,
		task.Name,
		task.Description,
		task.Status,
		task.Order,
		task.BoardID,
		task.CreatorID,
		task.AssigneeID,
	)
	return execErr
}

func (repo *TaskRepository) Update(ctx context.Context, id sqlddl.ID, d *models.UpdateTask) error {
	execErr := sqlquery.DynamicUpdate(ctx, repo.DB, &sqlquery.DynamicUpdateParams{
		TableName:   repositories.TableTasks,
		WhereColumn: sqlddl.ColumnID,
		WhereValue:  id,
		Changes: map[string]interface{}{
			repositories.ColumnName:        d.Name,
			repositories.ColumnDescription: d.Description,
			repositories.ColumnStatus:      d.Status,
			repositories.ColumnAssigneeID:  d.AssigneeID,
		},
		IsNilValue: func(value interface{}) bool {
			switch v := value.(type) {
			case *models.TaskStatus:
				return v == nil
			case *sqlddl.ID:
				return v == nil
			case *string:
				return v == nil
			default:
				return true
			}
		},
	})
	return execErr
}

func (repo *TaskRepository) FindByID(ctx context.Context, id sqlddl.ID) (*models.Task, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		sqlddl.ColumnID,
		repositories.ColumnBoardID,
		repositories.ColumnName,
		repositories.ColumnDescription,
		repositories.ColumnStatus,
		repositories.ColumnOrder,
		repositories.ColumnCreatorID,
		repositories.ColumnAssigneeID,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableTasks,
	)
	row := repo.DB.QueryRowContext(ctx, formattedQuery, id)
	var findTask models.Task
	scanErr := row.Scan(
		&findTask.ID,
		&findTask.BoardID,
		&findTask.Name,
		&findTask.Description,
		&findTask.Status,
		&findTask.Order,
		&findTask.CreatorID,
		&findTask.AssigneeID,
		&findTask.CreatedAt,
		&findTask.UpdatedAt,
	)
	return &findTask, scanErr
}

func (repo *TaskRepository) FindByOrder(ctx context.Context, boardId sqlddl.ID, order uint) (*models.Task, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1 AND %[2]s = $2"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.ColumnBoardID,
		repositories.ColumnOrder,
		repositories.ColumnName,
		sqlddl.ColumnID,
		repositories.ColumnDescription,
		repositories.ColumnStatus,
		repositories.ColumnCreatorID,
		repositories.ColumnAssigneeID,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableTasks,
	)
	row := repo.DB.QueryRowContext(ctx, formattedQuery, boardId, order)
	var findTask models.Task
	scanErr := row.Scan(
		&findTask.BoardID,
		&findTask.Order,
		&findTask.Name,
		&findTask.ID,
		&findTask.Description,
		&findTask.Status,
		&findTask.CreatorID,
		&findTask.AssigneeID,
		&findTask.CreatedAt,
		&findTask.UpdatedAt,
	)
	return &findTask, scanErr
}

func (repo *TaskRepository) FindByName(ctx context.Context, boardId sqlddl.ID, name string) (*models.Task, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1 AND %[2]s = $2"
	formatterQuery := fmt.Sprintf(
		query,
		repositories.ColumnBoardID,
		repositories.ColumnOrder,
		repositories.ColumnName,
		sqlddl.ColumnID,
		repositories.ColumnDescription,
		repositories.ColumnStatus,
		repositories.ColumnCreatorID,
		repositories.ColumnAssigneeID,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableTasks,
	)
	row := repo.DB.QueryRowContext(ctx, formatterQuery, boardId, name)
	var findTask models.Task
	scanErr := row.Scan(
		&findTask.BoardID,
		&findTask.Name,
		&findTask.ID,
		&findTask.Description,
		&findTask.Status,
		&findTask.Order,
		&findTask.CreatorID,
		&findTask.AssigneeID,
		&findTask.CreatedAt,
		&findTask.UpdatedAt,
	)
	return &findTask, scanErr
}

func (repo *TaskRepository) FindAllByBoardId(ctx context.Context, boardId sqlddl.ID) ([]models.Task, error) {
	const query = "SELECT %s, %s, %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %[1]s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.ColumnBoardID,
		sqlddl.ColumnID,
		repositories.ColumnName,
		repositories.ColumnDescription,
		repositories.ColumnStatus,
		repositories.ColumnOrder,
		repositories.ColumnCreatorID,
		repositories.ColumnAssigneeID,
		sqlddl.ColumnCreatedAt,
		sqlddl.ColumnUpdatedAt,
		repositories.TableTasks,
	)
	rows, rowsErr := repo.DB.QueryContext(ctx, formattedQuery, boardId)
	if rowsErr != nil {
		return nil, rowsErr
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		scanErr := rows.Scan(
			&task.BoardID,
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.Order,
			&task.CreatorID,
			&task.AssigneeID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (repo *TaskRepository) Delete(ctx context.Context, taskId sqlddl.ID) error {
	const query = "DELETE FROM %s WHERE %s = $1"
	formattedQuery := fmt.Sprintf(
		query,
		repositories.TableTasks,
		sqlddl.ColumnID,
	)
	_, execErr := repo.DB.ExecContext(ctx, formattedQuery, taskId)
	return execErr
}
