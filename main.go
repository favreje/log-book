package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "project_log.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logData := &LogData{}
	projectsMap, err := loadProjectsMap(db)
	if err != nil {
		log.Fatal(err)
	}

	inputState := &InputState{}

	// Main menu
	scanner := bufio.NewScanner(os.Stdin)
	for {
		displayMainMenu()
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		lowerline := strings.ToLower(line)
		char, _ := utf8.DecodeRuneInString(lowerline)

		switch char {
		case 'i':
			getUserData(logData, projectsMap, inputState)
			userConfirmation(db, logData, projectsMap, inputState)
		case 'd':
			projectId, err := selectProject(projectsMap)
			if err != nil {
				log.Fatal(err)
			}
			logRecords, err := readLogData(db, &projectId)
			if err != nil {
				log.Fatal(err)
			}
			reportByProject(logRecords, projectsMap)
			return
		case 'e':
			fmt.Println("Exiting the application...")
			return
		}
	}
}
