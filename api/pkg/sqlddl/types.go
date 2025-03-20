package sqlddl

import "fmt"

type ID string

const (
	TypeText = "TEXT"
	TypeInt  = "INT"
)

func TypeVarchar(n int) string {
	return fmt.Sprintf("VARCHAR(%d)", n)
}
