package contextkeys

import (
	"context"
	"errors"

	"just-kanban/pkg/sqlddl"
)

var (
	noUserIdErr = errors.New("no user identifier into context")
)

// GetUserId extracts user id from context
func GetUserId(ctx context.Context) (sqlddl.ID, error) {
	userId, ok := ctx.Value(KeyUserId).(sqlddl.ID)
	if !ok || userId == "" {
		return userId, noUserIdErr
	}
	return userId, nil
}
