package main

import (
	"database/sql"
)

func getProjectDesc(db *sql.DB, projectId int) (string, error) {
	var projectDesc string
	err := db.QueryRow("SELECT title FROM projects WHERE id = ?", projectId).Scan(&projectDesc)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
	}
	return projectDesc, nil
}
