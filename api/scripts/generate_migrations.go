package main

import (
	"fmt"

	"just-kanban/internal/repositories"
	"just-kanban/pkg/migrations"
)

func main() {
	migrationCount := 1
	for _, tableSchema := range repositories.Tables {
		migrations.GenerateTable(migrationCount, tableSchema.Name, migrations.DirectionUp, tableSchema)
		migrations.GenerateTable(migrationCount, tableSchema.Name, migrations.DirectionDown, tableSchema)
		migrationCount++
	}
	fmt.Println("all migrations generated successfully")
}
