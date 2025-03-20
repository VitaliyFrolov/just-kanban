package database

import (
	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"log"
)

func NewPostgresConnection(user, password, host, port, dbName string) *sql.DB {
	conn, connErr := sql.Open("postgres", fmt.Sprintf(
		// todo: enable ssl?
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	))
	if connErr != nil {
		log.Fatal(connErr)
	}
	return conn
}
