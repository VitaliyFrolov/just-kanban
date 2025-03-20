package models

// User is app user in business logic layer
type User struct {
	Model
	// Email is user email. Must be unique and have appropriate format
	Email string `db:"name" json:"email"`
	// Password is hashed user password
	Password string `db:"password" json:"-"`
	// Avatar is url to user avatar
	Avatar string `db:"avatar" json:"avatar"`
	// Username of user must be unique across all users
	Username string `db:"username" json:"username"`
	// FirstName of user it is not unique string which displaying in ui near avatar image
	FirstName string `db:"first_name" json:"first_name"`
	// LastName of user it is not unique string which displaying in ui near avatar image
	LastName string `db:"last_name" json:"last_name"`
}
