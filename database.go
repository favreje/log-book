package main

import (
	"database/sql"
)

func loadProjectsMap(db *sql.DB) (map[int]string, error) {
	projectsMap := make(map[int]string)
	rows, err := db.Query("SELECT id, title FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		err := rows.Scan(&id, &title)
		if err != nil {
			return nil, err
		}
		projectsMap[id] = title
	}
	return projectsMap, nil
}
