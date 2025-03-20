package sqlddl

type (
	SchemaTable struct {
		Name        string
		Columns     []SchemaColumn
		ForeignKeys []SchemaForeignKey
	}
	SchemaColumn struct {
		Name        string
		Type        string
		Constraints []string
	}
	SchemaForeignKey struct {
		ColumnName      string
		ReferenceTable  string
		ReferenceColumn string
		OnDelete        string
	}
)
