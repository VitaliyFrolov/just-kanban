package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"flag"
	"fmt"
	"log"

	"just-kanban/internal/config"
	"just-kanban/pkg/database"
	"just-kanban/pkg/migrations"
)

const (
	flagDirection = "direction"
)

func main() {
	env := config.NewEnv()
	db := database.NewPostgresConnection(env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
	instance, instanceErr := postgres.WithInstance(db, &postgres.Config{})
	if instanceErr != nil {
		log.Fatalln(instanceErr)
	}
	scripts, scriptsErr := migrate.NewWithDatabaseInstance(
		"file://migrations/output",
		env.DBName,
		instance,
	)
	if scriptsErr != nil {
		log.Fatalln(scriptsErr)
	}
	direction := flag.String(flagDirection, "", "Direction of migration, if not provided then runs all")
	flag.Parse()
	switch migrations.Direction(*direction) {
	case migrations.DirectionUp:
		runErr := scripts.Up()
		if runErr != nil {
			log.Fatalln(runErr)
		}
	case migrations.DirectionDown:
		runErr := scripts.Down()
		if runErr != nil {
			log.Fatalln(runErr)
		}
	case migrations.DirectionDrop:
		runErr := scripts.Drop()
		if runErr != nil {
			log.Fatalln(runErr)
		}
	default:
		log.Fatalln(fmt.Errorf("invalid migration direction: %s", migrations.Direction(*direction)))
	}
	log.Println("All migrations finished")
}
