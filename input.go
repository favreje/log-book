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
	displayUserInput(logData, projectsMap)
	if !getProjId(logData, projectsMap) {
		return
	}
	inputDate, ok := getLogDate(logData, projectsMap)
	if !ok {
		return
	}
	if !getLogTime("Start", inputDate, logData, projectsMap) {
		return
	}
	if !getLogTime("End", inputDate, logData, projectsMap) {
		return
	}

	duration, err := logData.calculateDuration()
	if err != nil {
		return
	}
	logData.duration = duration
	displayUserInput(logData, projectsMap)

	if !getCategory(logData, projectsMap) {
		return
	}
	if !getDescription(logData, projectsMap) {
		return
	}
}

func getProjId(logData *LogData, projectsMap map[int]string) bool {
	for {
		userInput, ok := getUserInput("Project ID")
		if !ok {
			displayUserInput(logData, projectsMap)
			return false
		}

		line := strings.TrimSpace(userInput)
		projectId, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Invalid project ID: %v\n", err)
			continue
		}

		projectDesc, ok := projectsMap[projectId]
		if !ok {
			displayUserInput(logData, projectsMap)
			fmt.Println("No description matches Project Code:", projectId)
			continue
		}
		fmt.Printf("Project: %s\n", projectDesc)
		logData.projectId = projectId
		break
	}
	displayUserInput(logData, projectsMap)
	return true
}

func getLogDate(logData *LogData, projectsMap map[int]string) (time.Time, bool) {
	// Input date for start time and end time
	for {
		userInput, ok := getUserInput("Log Date (MM/DD/YY)")
		if !ok {
			displayUserInput(logData, projectsMap)
			return time.Time{}, false
		}

		line := strings.TrimSpace(userInput)
		parsedDate, err := time.Parse("01/02/06", line)
		if err != nil {
			displayUserInput(logData, projectsMap)
			fmt.Printf("Invalid Input Date: %v\n", err)
			continue
		}
		logData.startTime = parsedDate
		logData.endTime = parsedDate
		displayUserInput(logData, projectsMap)
		return parsedDate, true
	}
}

func getLogTime(
	boundaryType string,
	inputDate time.Time,
	logData *LogData,
	projectsMap map[int]string,
) bool {
	// Populates logData.startTime or logData.endTime based on boundaryType
	if !(boundaryType == "Start" || boundaryType == "End") {
		log.Fatalf("Invalid boundaryType: %v\n", boundaryType)
	}
	for {
		prompt := boundaryType + "Time (HH:MM)"
		userInput, ok := getUserInput(prompt)
		if !ok {
			logData.startTime = inputDate
			logData.endTime = inputDate
			displayUserInput(logData, projectsMap)
			return false
		}
		line := strings.TrimSpace(userInput)
		dateString := inputDate.Format("01/02/06") + " " + line
		layout := "01/02/06 15:04"
		inputTime, err := time.ParseInLocation(layout, dateString, time.Local)
		if err != nil {
			displayUserInput(logData, projectsMap)
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
	displayUserInput(logData, projectsMap)
	return true
}

func getCategory(logData *LogData, projectsMap map[int]string) bool {
	for {
		userInput, ok := getUserInput("Category")
		if !ok {
			displayUserInput(logData, projectsMap)
			return false
		}
		line := strings.TrimSpace(userInput)
		if line == "" {
			displayUserInput(logData, projectsMap)
			fmt.Println("Please provide a project category.")
			continue
		}
		logData.category = line
		break
	}
	displayUserInput(logData, projectsMap)
	return true
}

func getDescription(logData *LogData, projectsMap map[int]string) bool {
	for {
		userInput, ok := getUserInput("Description")
		if !ok {
			displayUserInput(logData, projectsMap)
			return false
		}
		line := strings.TrimSpace(userInput)
		if line == "" {
			displayUserInput(logData, projectsMap)
			fmt.Println("No description was entered")
			continue
		}
		logData.description = line
		break
	}
	displayUserInput(logData, projectsMap)
	return true
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
			*logData = LogData{}
			return
		default:
			fmt.Printf("%s is invalid. Please enter again.", line)
			time.Sleep(1750 * time.Millisecond)
		}
	}
}

func userEdit(scanner *bufio.Scanner, logData *LogData, projectsMap map[int]string) {
outerLoop:
	for {
		displayUserInput(logData, projectsMap)
		fmt.Println("EDIT MODE:")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Println(
			"(P)roject ID | (L)og date | (S)tart time | (E)nd time | (C)ategory | (D)escription | (R)eturn",
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
			displayUserInput(logData, projectsMap)
			getProjId(logData, projectsMap)
			continue outerLoop
		case 'l':
			displayUserInput(logData, projectsMap)
			revisedDate, ok := getLogDate(logData, projectsMap)
			if !ok {
				continue outerLoop
			}
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
			continue outerLoop
		case 's':
			displayUserInput(logData, projectsMap)
			getLogTime("Start", logData.startTime, logData, projectsMap)
			recalculatedDuration, err := logData.calculateDuration()
			if err == nil {
				logData.duration = recalculatedDuration
			}
			continue outerLoop
		case 'e':
			displayUserInput(logData, projectsMap)
			getLogTime("End", logData.endTime, logData, projectsMap)
			recalculatedDuration, err := logData.calculateDuration()
			if err == nil {
				logData.duration = recalculatedDuration
			}
			continue outerLoop
		case 'c':
			displayUserInput(logData, projectsMap)
			getCategory(logData, projectsMap)
			continue outerLoop
		case 'd':
			displayUserInput(logData, projectsMap)
			getDescription(logData, projectsMap)
			continue outerLoop
		case 'r':
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
