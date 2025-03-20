package models

import (
	"just-kanban/internal/access"
	"just-kanban/pkg/sqlddl"
)

// BoardMember is member of project board in business logic layer
type BoardMember struct {
	Model
	// UserID is identifier of user that is project board member
	UserID sqlddl.ID `json:"user_id"`
	// BoardID is identifier of project board which member related to
	BoardID sqlddl.ID `json:"board_id"`
	// Role defines member accesses to desk, must be sync with access.Role constants
	Role access.Role `json:"role" validate:"oneof=owner manager regular"`
}
