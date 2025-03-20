package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"just-kanban/pkg/sqlddl"
)

type Direction string

const (
	DirectionUp   Direction = "up"
	DirectionDown           = "down"
	DirectionDrop           = "drop"
	// todo: incorrect dependency usage, pkg/migrations uses depends on specific project and we must define templates there
	dirMigrations string = "migrations"
	dirOutput            = "output"
	dirTemplates         = "templates"
)

func GenerateTable(n int, tableName string, direction Direction, data sqlddl.SchemaTable) {
	tmplDir := filepath.Join(dirMigrations, dirTemplates)
	tmplName := fmt.Sprintf("%s_table.%[1]s.sql.tmpl", direction)
	tmplPath := filepath.Join(tmplDir, tmplName)
	outputDir := filepath.Join(dirMigrations, dirOutput)
	outputName := fmt.Sprintf("%04d_%s_table.%s.sql", n, tableName, direction)
	outputPath := filepath.Join(outputDir, outputName)
	tmp := template.Must(
		template.New(tmplName).Funcs(template.FuncMap{"join": strings.Join}).ParseFiles(tmplPath),
	)
	mkDirsErr := os.MkdirAll(outputDir, os.ModePerm)
	if mkDirsErr != nil {
		log.Fatal(mkDirsErr)
	}
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	execErr := tmp.Execute(file, map[string]any{
		"TableName":       tableName,
		"ColumnID":        sqlddl.ColumnID,
		"Columns":         data.Columns,
		"ForeignKeys":     data.ForeignKeys,
		"ColumnCreatedAt": sqlddl.ColumnCreatedAt,
		"ColumnUpdatedAt": sqlddl.ColumnUpdatedAt,
	})
	if execErr != nil {
		log.Fatal(execErr)
	}
	fmt.Printf("%s migration for \"%s\" table is successfuly generated!\n", direction, tableName)
}
