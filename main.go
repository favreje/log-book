package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "project_log.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logData := LogData{}

	// Testing getUserData func
	getUserData(&logData)
	displayUserInput(db, &logData)
}
