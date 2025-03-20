package contextkeys

import (
	"context"
	"testing"

	"just-kanban/pkg/sqlddl"
)

func TestGetUserId(t *testing.T) {
	ctx := context.WithValue(context.Background(), KeyUserId, sqlddl.ID("test_user_id"))
	result, err := GetUserId(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if result != "test_user_id" {
		t.Fatal("Incorrect user id extracted")
	}
}
