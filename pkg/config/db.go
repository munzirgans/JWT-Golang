package config

import (
	"database/sql"
	"log"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "munzir:munzirdev@tcp(localhost:3306)/jwt")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
