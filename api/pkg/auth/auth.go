package auth

import "errors"

var (
	UnauthorizedErr = errors.New("unauthorized")
	TokenHeader     = "Authorization"
)
