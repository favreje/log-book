package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func getUserData(logData *LogData, projectsMap map[int]string) {
	scanner := bufio.NewScanner(os.Stdin)

	getProjId(logData, scanner, projectsMap)
	inputDate := getLogDate(scanner)
	getLogTime("Start", inputDate, logData, scanner)
	getLogTime("End", inputDate, logData, scanner)

	duration, err := logData.calculateDuration()
	if err != nil {
		fmt.Printf("Duration could not be calculated: %v\n", err)
		return
	}
	logData.duration = duration

	getCategory(logData, scanner)
	getDescription(logData, scanner)
}

func getProjId(logData *LogData, scanner *bufio.Scanner, projectsMap map[int]string) {
	for {
		fmt.Print("Project ID: ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		projectId, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Invalid project ID: %v\n", err)
			continue
		}

		projectDesc, ok := projectsMap[projectId]
		if !ok {
			fmt.Println("No description matches Project Code:", projectId)
			continue
		}
		fmt.Printf("Project: %s\n", projectDesc)
		logData.projectId = projectId
		break
	}
}

func getLogDate(scanner *bufio.Scanner) time.Time {
	// Input date for start time and end time
	for {
		fmt.Print("Log Date (MM/DD/YY): ")
		if !scanner.Scan() {
			return time.Time{}
		}

		line := strings.TrimSpace(scanner.Text())
		parsedDate, err := time.Parse("01/02/06", line)
		if err != nil {
			fmt.Printf("Invalid Input Date: %v\n", err)
			continue
		}
		return parsedDate
	}
}

func getLogTime(
	boundaryType string,
	inputDate time.Time,
	logData *LogData,
	scanner *bufio.Scanner,
) {
	// Populates logData.startTime or logData.endTime based on boundaryType
	if !(boundaryType == "Start" || boundaryType == "End") {
		fmt.Printf("Invalid boundaryType: %v\n", boundaryType)
		return
	}
	for {
		fmt.Printf("%s Time (HH:MM): ", boundaryType)
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		dateString := inputDate.Format("01/02/06") + " " + line
		layout := "01/02/06 15:04"
		// inputTime, err := time.Parse(layout, dateString)
		inputTime, err := time.ParseInLocation(layout, dateString, time.Local)
		if err != nil {
			fmt.Printf("Invalid %s Time: %v\n", boundaryType, err)
			continue
		}
		if boundaryType == "Start" {
			logData.startTime = inputTime
			break
		}
		logData.endTime = inputTime
		break
	}
}

func getCategory(logData *LogData, scanner *bufio.Scanner) {
	for {
		fmt.Print("Category: ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Println("No category was entered")
			continue
		}
		logData.category = line
		break
	}
}

func getDescription(logData *LogData, scanner *bufio.Scanner) {
	for {
		fmt.Print("Description: ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Println("No description was entered")
			continue
		}
		logData.description = line
		break
	}
}

func userConfirmation(db *sql.DB, logData *LogData, projectsMap map[int]string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		displayUserInput(logData, projectsMap)
		fmt.Println("(W)rite | (E)dit | (C)ancel")
		fmt.Print("Selection: ")
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		lowerLine := strings.ToLower(line)
		char, _ := utf8.DecodeRuneInString(lowerLine)

		switch char {
		case 'w':
			fmt.Println("Writing data to log...")
			insertId, err := writeLogEntry(db, logData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Successfully saved log entry with ID: %d\n", insertId)
			return
		case 'e':
			userEdit(scanner, logData, projectsMap)
		case 'c':
			fmt.Println("Cancelling the log entry...")
			return
		default:
			fmt.Printf("%s is invalid. Please enter again.", line)
			time.Sleep(1750 * time.Millisecond)
		}
	}
}

func userEdit(scanner *bufio.Scanner, logData *LogData, projectsMap map[int]string) {
	for {
		displayUserInput(logData, projectsMap)
		fmt.Println("EDIT MODE:")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Println(
			"(P)roject ID | (L)og date | (S)tart time | (E)nd time | (C)ategory | (D)escription | (Q)uit ",
		)
		fmt.Print("Selection: ")
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		lowerline := strings.ToLower(line)
		char, _ := utf8.DecodeRuneInString(lowerline)

		switch char {
		case 'p':
			getProjId(logData, scanner, projectsMap)
			return
		case 'l':
			revisedDate := getLogDate(scanner)
			layout := "01/02/06 15:04"
			// Revise start time
			timeString := logData.startTime.Format("15:04")
			dateString := revisedDate.Format("01/02/06") + " " + timeString
			revisedTime, err := time.ParseInLocation(layout, dateString, time.Local)
			if err != nil {
				fmt.Printf("Invalid Start Time: %v\n", err)
				continue
			}
			logData.startTime = revisedTime

			// Revise end time
			timeString = logData.endTime.Format("15:04")
			dateString = revisedDate.Format("01/02/06") + " " + timeString

			revisedTime, err = time.ParseInLocation(layout, dateString, time.Local)
			if err != nil {
				fmt.Printf("Invalid End Time: %v\n", err)
				continue
			}
			logData.endTime = revisedTime
			return
		case 's':
			getLogTime("Start", logData.startTime, logData, scanner)
			recalculatedDuration, err := logData.calculateDuration()
			if err != nil {
				fmt.Printf("Duration could not be calculated: %v\n", err)
				time.Sleep(1750 * time.Millisecond)
				return
			}
			logData.duration = recalculatedDuration
			return
		case 'e':
			getLogTime("End", logData.endTime, logData, scanner)
			recalculatedDuration, err := logData.calculateDuration()
			if err != nil {
				fmt.Printf("Duration could not be calculated: %v\n", err)
				time.Sleep(1750 * time.Millisecond)
				return
			}
			logData.duration = recalculatedDuration
			return
		case 'c':
			getCategory(logData, scanner)
			return
		case 'd':
			getDescription(logData, scanner)
			return
		case 'q':
			return
		default:
			fmt.Printf("%s is invalid. Please enter again.", line)
			time.Sleep(1750 * time.Millisecond)
		}
	}
}

func selectProject(projectsMap map[int]string) (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	clearScreen()
	for {
		fmt.Println("Enter report ID (or 'D' to display list)")
		fmt.Print("Selection: ")
		if !scanner.Scan() {
			return 0, scanner.Err()
		}
		line := strings.TrimSpace(scanner.Text())
		lowerline := strings.ToLower(line)
		char, _ := utf8.DecodeRuneInString(lowerline)

		if char == 'd' {
			clearScreen()
			displayProjectList(projectsMap)
			continue
		}
		id, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Invalid project ID: %v\n", err)
			continue
		}
		projectDesc, ok := projectsMap[id]
		if !ok {
			fmt.Printf("Invalid project ID: %v\n", err)
			continue
		}
		clearScreen()
		fmt.Printf("ID: %02d\nProject: %s\n", id, projectDesc)
		if confirmedSelection("Display report?") {
			return id, nil
		}
	}
}
