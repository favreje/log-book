package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func getUserData(logData *LogData, projectsMap map[int]string, inputState *InputState) {
	displayUserInput(logData, projectsMap, inputState)
	fmt.Println("INPUT MODE")
	if !getProjId(logData, projectsMap, inputState) {
		return
	}
	if !getLogDate(logData, projectsMap, inputState) {
		return
	}
	if !getLogTime(Start, logData, projectsMap, inputState) {
		return
	}
	if !getLogTime(End, logData, projectsMap, inputState) {
		return
	}

	duration, err := logData.calculateDuration()
	if err != nil {
		return
	}
	logData.duration = duration
	displayUserInput(logData, projectsMap, inputState)
	fmt.Println("INPUT MODE")

	if !getCategory(logData, projectsMap, inputState) {
		return
	}
	if !getDescription(logData, projectsMap, inputState) {
		return
	}
}

func getProjId(logData *LogData, projectsMap map[int]string, inputState *InputState) bool {
	for {
		userInput, ok := getUserInput("Project ID")
		if !ok {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			return false
		}

		line := strings.TrimSpace(userInput)
		projectId, err := strconv.Atoi(line)
		if err != nil {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			fmt.Println("Invalid project ID")
			continue
		}

		projectDesc, ok := projectsMap[projectId]
		if !ok {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			fmt.Println("No description matches Project Code:", projectId)
			continue
		}
		fmt.Printf("Project: %s\n", projectDesc)
		logData.projectId = projectId
		break
	}
	displayUserInput(logData, projectsMap, inputState)
	fmt.Println("INPUT MODE")
	return true
}

func getLogDate(
	logData *LogData,
	projectsMap map[int]string,
	inputState *InputState,
) bool {
	// Input date for start time and end time
	for {
		userInput, ok := getUserInput("Log Date (MM/DD/YY)")
		if !ok {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			return false
		}

		line := strings.TrimSpace(userInput)
		parsedDate, err := time.Parse("01/02/06", line)
		if err != nil {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			fmt.Printf("Invalid Input Date: %v\n", err)
			continue
		}
		inputState.dateEntered = true
		inputState.baseDate = parsedDate
		displayUserInput(logData, projectsMap, inputState)
		fmt.Println("INPUT MODE")
		return true
	}
}

func getLogTime(
	boundaryType Boundary,
	logData *LogData,
	projectsMap map[int]string,
	inputState *InputState,
) bool {
	// Populates logData.startTime or logData.endTime based on boundaryType
	if !(boundaryType == Start || boundaryType == End) {
		log.Fatalf("Invalid boundaryType: %v\n", boundaryType)
	}

	// First ensure that a date has been entered, and if not, go do that first
	if !inputState.dateEntered {
		if !getLogDate(logData, projectsMap, inputState) {
			return false
		}
	}

	inputDate := inputState.baseDate

	for {
		prompt := string(boundaryType) + "Time (HH:MM)"
		userInput, ok := getUserInput(prompt)
		if !ok {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			return false
		}
		line := strings.TrimSpace(userInput)
		dateString := inputDate.Format("01/02/06") + " " + line
		layout := "01/02/06 15:04"
		inputTime, err := time.ParseInLocation(layout, dateString, time.Local)
		if err != nil {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			fmt.Printf("Invalid %s Time: %v\n", boundaryType, err)
			continue
		}
		if boundaryType == Start {
			inputState.startTimeEntered = true
			logData.startTime = inputTime
			break
		}
		inputState.endTimeEntered = true
		logData.endTime = inputTime
		break
	}
	displayUserInput(logData, projectsMap, inputState)
	fmt.Println("INPUT MODE")
	return true
}

func getCategory(logData *LogData, projectsMap map[int]string, inputState *InputState) bool {
	for {
		userInput, ok := getUserInput("Category")
		if !ok {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			return false
		}
		line := strings.TrimSpace(userInput)
		if line == "" {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			fmt.Println("Please provide a project category.")
			continue
		}
		logData.category = line
		break
	}
	displayUserInput(logData, projectsMap, inputState)
	fmt.Println("INPUT MODE")
	return true
}

func getDescription(logData *LogData, projectsMap map[int]string, inputState *InputState) bool {
	for {
		userInput, ok := getUserInput("Description")
		if !ok {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			return false
		}
		line := strings.TrimSpace(userInput)
		if line == "" {
			displayUserInput(logData, projectsMap, inputState)
			fmt.Println("INPUT MODE")
			fmt.Println("No description was entered")
			continue
		}
		logData.description = line
		break
	}
	displayUserInput(logData, projectsMap, inputState)
	fmt.Println("INPUT MODE")
	return true
}

func userConfirmation(
	db *sql.DB,
	logData *LogData,
	projectsMap map[int]string,
	inputState *InputState,
) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		displayUserInput(logData, projectsMap, inputState)
		fmt.Println("(W)rite | (E)dit | (C)ancel")
		fmt.Print("Selection: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Fatalf("Error reading input: %v", err)
			}
			scanner = bufio.NewScanner(os.Stdin)
			continue
		}
		line := strings.TrimSpace(scanner.Text())
		lowerLine := strings.ToLower(line)
		char, _ := utf8.DecodeRuneInString(lowerLine)

		switch char {
		case 'w':
			insertId, err := writeLogEntry(db, logData)
			if err != nil {
				log.Fatal(err)
			}
			statusMsg := fmt.Sprintf("Successfully saved log entry with ID: %d", insertId)
			*logData = LogData{}
			*inputState = InputState{}
			inputState.statusMsg = statusMsg
			return
		case 'e':
			userEdit(scanner, logData, projectsMap, inputState)
		case 'c':
			*logData = LogData{}
			*inputState = InputState{}
			return
		default:
			fmt.Printf("%s is invalid. Please enter again.", line)
			time.Sleep(1750 * time.Millisecond)
		}
	}
}

func userEdit(
	scanner *bufio.Scanner,
	logData *LogData,
	projectsMap map[int]string,
	inputState *InputState,
) {
outerLoop:
	for {
		displayUserInput(logData, projectsMap, inputState)
		fmt.Println("EDIT MODE")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Println(
			"(P)roject ID | (L)og date | (S)tart time | (E)nd time | (C)ategory | (D)escription | (R)eturn",
		)
		fmt.Print("Selection: ")
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
		case 'p':
			displayUserInput(logData, projectsMap, inputState)
			getProjId(logData, projectsMap, inputState)
			continue outerLoop
		case 'l':
			displayUserInput(logData, projectsMap, inputState)
			if !getLogDate(logData, projectsMap, inputState) {
				continue outerLoop
			}

			layout := "01/02/06 15:04"
			// Revise start time
			if inputState.startTimeEntered {
				timeString := logData.startTime.Format("15:04")
				dateString := inputState.baseDate.Format("01/02/06") + " " + timeString
				revisedTime, err := time.ParseInLocation(layout, dateString, time.Local)
				if err != nil {
					fmt.Printf("Invalid Start Time: %v\n", err)
					continue
				}
				logData.startTime = revisedTime
			}

			// Revise end time
			if inputState.endTimeEntered {
				timeString := logData.endTime.Format("15:04")
				dateString := inputState.baseDate.Format("01/02/06") + " " + timeString

				revisedTime, err := time.ParseInLocation(layout, dateString, time.Local)
				if err != nil {
					fmt.Printf("Invalid End Time: %v\n", err)
					continue
				}
				logData.endTime = revisedTime
			}
			continue outerLoop
		case 's':
			displayUserInput(logData, projectsMap, inputState)
			getLogTime(Start, logData, projectsMap, inputState)
			if inputState.startTimeEntered && inputState.endTimeEntered {
				recalculatedDuration, err := logData.calculateDuration()
				if err == nil {
					logData.duration = recalculatedDuration
				}
			}
			continue outerLoop
		case 'e':
			displayUserInput(logData, projectsMap, inputState)
			getLogTime(End, logData, projectsMap, inputState)
			if inputState.startTimeEntered && inputState.endTimeEntered {
				recalculatedDuration, err := logData.calculateDuration()
				if err == nil {
					logData.duration = recalculatedDuration
				}
			}
			continue outerLoop
		case 'c':
			displayUserInput(logData, projectsMap, inputState)
			getCategory(logData, projectsMap, inputState)
			continue outerLoop
		case 'd':
			displayUserInput(logData, projectsMap, inputState)
			getDescription(logData, projectsMap, inputState)
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
	fmt.Println("Enter report ID ('D' to display list or Ctrl-D to exit)")
	for {
		fmt.Print("Selection: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}
		line := strings.TrimSpace(scanner.Text())
		lowerline := strings.ToLower(line)
		char, _ := utf8.DecodeRuneInString(lowerline)

		if char == 'd' {
			clearScreen()
			displayProjectList(projectsMap)
			fmt.Println("Please select from the list above (or Ctrl-D to exit).")
			continue
		}
		id, err := strconv.Atoi(line)
		if err != nil {
			clearScreen()
			displayProjectList(projectsMap)
			fmt.Println(
				"Invalid project ID.\nPlease select from the list above (or Ctrl-D to exit).",
			)
			continue
		}
		projectDesc, ok := projectsMap[id]
		if !ok {
			clearScreen()
			displayProjectList(projectsMap)
			fmt.Println(
				"Invalid project ID.\nPlease select from the list above (or Ctrl-D to exit).",
			)
			continue
		}
		clearScreen()
		fmt.Printf("ID: %02d\nProject: %s\n", id, projectDesc)
		if confirmedSelection("Display report?") {
			return id, nil
		}
		clearScreen()
		displayProjectList(projectsMap)
		fmt.Println("Please select from the list above (or Ctrl-D to exit).")
	}
}
