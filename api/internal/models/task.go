package models

import "just-kanban/pkg/sqlddl"

type TaskStatus uint

const (
	TaskStatusBacklog TaskStatus = iota + 1
	TaskStatusProcess
	TaskStatusDone
)

// Task is project board's task in business logic layer
type Task struct {
	Model
	// BoardID is identifier of project board that in task lives
	BoardID sqlddl.ID `db:"board_id" json:"board_id"`
	// CreatorID is identifier of user that initial created task
	CreatorID sqlddl.ID `db:"creator_id" json:"creator_id"`
	// AssigneeID is identifier of user who is task assignee
	AssigneeID sqlddl.ID `db:"assignee_id" json:"assignee_id"`
	// Order is task position on its project board (BoardID), incremental
	Order int `db:"order" json:"order"`
	// Name is task title
	Name string `db:"name" json:"name"`
	// Description is task description which contains info about subject of task
	Description string `db:"description" json:"description"`
	// Status is task status that it's on at this moment
	Status TaskStatus `db:"status" json:"status"`
}
