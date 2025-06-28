package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
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
		displayMainMenu(inputState)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Fatalf("Error reading input: %v", err)
			}
			scanner = bufio.NewScanner(os.Stdin)
			continue
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
				if errors.Is(err, io.EOF) {
					continue
				}
				log.Fatalf("Fatal input error: %v", err)
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
