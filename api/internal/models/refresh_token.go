package models

import "just-kanban/pkg/sqlddl"

// RefreshToken is refresh token in business logic layer
type RefreshToken struct {
	Model
	// Token is string that been generated as refresh token, will be compared to another string, which is possible tokens
	Token string `db:"token" json:"token"`
	// UserID is identifier of user who for refresh token been generated
	UserID sqlddl.ID `db:"user_id" json:"user_id"`
}
