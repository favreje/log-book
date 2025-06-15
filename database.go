package main

import (
	"database/sql"
	"fmt"
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

func writeLogEntry(db *sql.DB, logData *LogData) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO log (project_id, start_time, end_time, category, description) VALUES (?, ?, ?, ?, ?)",
		logData.projectId,
		logData.startTime,
		logData.endTime,
		logData.category,
		logData.description,
	)
	if err != nil {
		return 0, fmt.Errorf("Failed to insert log entry %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Could not verify insert: %w", err)
	}
	if rowsAffected != 1 {
		return 0, fmt.Errorf("Expected one affected row, got %d rows", rowsAffected)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Could not retrieve insert ID: %w", err)
	}
	return insertId, nil
}

func readLogData(db *sql.DB, projectId *int) ([]LogData, error) {
	var rows *sql.Rows
	var err error
	if projectId != nil {
		query := `
      SELECT rowid, project_id, start_time, end_time, category, description
      FROM log
      WHERE project_id = ?
      ORDER BY project_id, start_time
    `
		rows, err = db.Query(query, *projectId)
	} else {
		query := `
      SELECT rowid, project_id, start_time, end_time, category, description
      FROM log
      ORDER BY project_id, start_time
    `
		rows, err = db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logRecords := []LogData{}
	for rows.Next() {
		var rowId int
		var logData LogData
		var startTimeStr string
		var endTimeStr string

		err := rows.Scan(
			&rowId,
			&logData.projectId,
			&startTimeStr,
			&endTimeStr,
			&logData.category,
			&logData.description,
		)
		if err != nil {
			return nil, err
		}

		startTime, err := parseTimeFromDb(startTimeStr)
		if err != nil {
			return nil, err
		}
		logData.startTime = startTime

		endTime, err := parseTimeFromDb(endTimeStr)
		if err != nil {
			return nil, err
		}
		logData.endTime = endTime
		duration, err := logData.calculateDuration()
		if err != nil {
			return nil, err
		}
		logData.duration = duration

		logRecords = append(logRecords, logData)
	}
	return logRecords, nil
}
