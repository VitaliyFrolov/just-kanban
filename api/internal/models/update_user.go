package models

// UpdateUser is data to update User
type UpdateUser struct {
	Avatar    *string `json:"avatar"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
