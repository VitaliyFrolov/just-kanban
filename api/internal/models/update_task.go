package models

import "just-kanban/pkg/sqlddl"

// UpdateTask is data to update Task
type UpdateTask struct {
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	Status      *TaskStatus `json:"status"`
	AssigneeID  *sqlddl.ID  `json:"assignee_id"`
}
