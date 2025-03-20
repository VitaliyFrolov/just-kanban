package sqlddl

import "testing"

func TestTypeVarchar(t *testing.T) {
	result := TypeVarchar(100)
	if result != "VARCHAR(100)" {
		t.Fatalf("expect VARCHAR(100) but got %s", result)
	}
}
