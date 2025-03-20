package models

// UpdateBoard is data to update Board
type UpdateBoard struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
