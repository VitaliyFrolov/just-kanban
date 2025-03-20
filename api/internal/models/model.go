package models

import (
	"time"

	"just-kanban/pkg/sqlddl"
)

// Model is abstract model in business logic layer, should be used as base for another models
type Model struct {
	sqlddl.ID `db:"id" json:"id"`
	// CreatedAt contains info as timestamp when this model been created
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	// UpdatedAt contains info as timestamp when this model been updated last time
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
