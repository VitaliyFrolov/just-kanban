package access

// Role is a status of board member, which on board management depends
type Role string

const (
	// RoleOwner is the one who has privileges for management administration and board delete
	RoleOwner Role = "owner"
	// RoleManager is a one who has possibility for tasks creation and regular members access management
	RoleManager = "manager"
	// RoleRegular is a regular member of a board whose accesses defined by managers
	RoleRegular = "regular"
)
