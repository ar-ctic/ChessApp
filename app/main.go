package main

import (
	"ChessApp/api"
	"ChessApp/db"
	"database/sql"
	"log"
)

func main() {

	db, err := db.NewSQLiteStorage()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":5000", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")

}
