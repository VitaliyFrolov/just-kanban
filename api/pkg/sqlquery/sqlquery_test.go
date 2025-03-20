package sqlquery

import (
	"github.com/DATA-DOG/go-sqlmock"

	"context"
	"testing"
)

func TestDynamicUpdate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	t.Run("Update single record's field", func(t *testing.T) {
		expectation := mock.ExpectExec("UPDATE table SET status = \\$1 WHERE id = \\$2")
		expectation.WithArgs("Updated", 1)
		expectation.WillReturnResult(sqlmock.NewResult(1, 1))
		status := "Updated"
		err := DynamicUpdate(context.Background(), db, &DynamicUpdateParams{
			TableName:   "table",
			WhereColumn: "id",
			WhereValue:  1,
			Changes: map[string]interface{}{
				"status": &status,
			},
			IsNilValue: func(value interface{}) bool {
				return false
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if expectErr := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(expectErr)
		}
	})
	t.Run("Update multiple record's field", func(t *testing.T) {
		expectation := mock.ExpectExec("UPDATE table SET (status = \\$1, order = \\$2) WHERE id = \\$3")
		expectation.WithArgs("Updated", 2, 1)
		expectation.WillReturnResult(sqlmock.NewResult(1, 1))
		status := "Updated"
		order := 2
		err := DynamicUpdate(context.Background(), db, &DynamicUpdateParams{
			TableName:   "table",
			WhereColumn: "id",
			WhereValue:  1,
			Changes: map[string]interface{}{
				"status": &status,
				"order":  &order,
			},
			IsNilValue: func(value interface{}) bool {
				return false
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if expectErr := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(expectErr)
		}
	})
	t.Run("IsNilValue variation", func(t *testing.T) {
		expectation := mock.ExpectExec("UPDATE table SET order = \\$1 WHERE id = \\$2")
		expectation.WithArgs(2, 1)
		expectation.WillReturnResult(sqlmock.NewResult(1, 1))
		DynamicUpdate(context.Background(), db, &DynamicUpdateParams{
			TableName:   "table",
			WhereColumn: "id",
			WhereValue:  1,
			Changes: map[string]interface{}{
				"order":  2,
				"status": "Default",
			},
			IsNilValue: func(value interface{}) bool {
				switch (value).(type) {
				case string:
					return true
				default:
					return false
				}
			},
		})
		if expectErr := mock.ExpectationsWereMet(); expectErr != nil {
			t.Fatal(expectErr)
		}
	})
}
