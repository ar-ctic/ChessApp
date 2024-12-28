package db

import (
	"database/sql"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func init() {

	db, err := sql.Open("sqlite3", "file:demo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY NOT NULL,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
	`

	_, err = db.Exec(userTable)

	if err != nil {
		log.Fatal(err)
	}
}

func NewSQLiteStorage() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "demo.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
