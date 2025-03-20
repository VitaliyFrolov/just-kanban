package models

// Board is project board in business logic layer.
type Board struct {
	Model
	// Name is the board's title.
	Name string `db:"name" json:"name"`
	// Description summarizes the board's purpose.
	Description string `db:"description" json:"description"`
}
