package sqlquery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	noUpdateClausesErr = errors.New("no update clauses")
)

// DynamicUpdateParams defines params of a relational database update
type DynamicUpdateParams struct {
	// TableName is target table for update
	TableName string
	// WhereColumn defines column of "where" update condition
	WhereColumn string
	// WhereValue defines value of "where" update condition
	WhereValue interface{}
	// Changes is map of pointers to primitives
	Changes map[string]interface{}
	// IsNilValue is function for checking if a field of Changes shouldn't be updated
	IsNilValue func(value interface{}) bool
}

// DynamicUpdate updates relational database with DynamicUpdateParams
func DynamicUpdate(ctx context.Context, db *sql.DB, dup *DynamicUpdateParams) error {
	const queryStart = "UPDATE %s SET %s WHERE %s = $%d"
	var clausesStrings []string
	argsIndex := 0
	var args []interface{}
	for key, value := range dup.Changes {
		if !dup.IsNilValue(value) {
			clausesStrings = append(clausesStrings, fmt.Sprintf("%s = $%d", key, argsIndex+1))
			args = append(args, value)
			argsIndex++
		}
	}
	if len(clausesStrings) == 0 {
		return noUpdateClausesErr
	}
	clauses := strings.Join(clausesStrings, ", ")
	formattedQuery := fmt.Sprintf(queryStart, dup.TableName, clauses, dup.WhereColumn, len(args)+1)
	args = append(args, dup.WhereValue)
	_, execErr := db.ExecContext(ctx, formattedQuery, args...)
	return execErr
}
