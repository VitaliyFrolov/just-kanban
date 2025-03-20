package contextkeys

type ctxKey int

const (
	// KeyUserId is context key for user id, usually represents active (authenticated) user
	KeyUserId ctxKey = iota
)
