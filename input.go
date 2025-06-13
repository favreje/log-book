package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func getUserData(logData *LogData) {
	scanner := bufio.NewScanner(os.Stdin)

	getProjId(logData, scanner)
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

func getProjId(logData *LogData, scanner *bufio.Scanner) {
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
		inputTime, err := time.Parse(layout, dateString)
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

func userConfirmation(logData *LogData, projectsMap map[int]string) {
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
			getProjId(logData, scanner)
		case 'l':
			revisedDate := getLogDate(scanner)
			// Revise start time
			timeString := logData.startTime.Format("15:04")
			dateString := revisedDate.Format("01/02/06") + " " + timeString
			revisedTime, err := time.Parse("01/02/06 15:04", dateString)
			if err != nil {
				fmt.Printf("Invalid Start Time: %v\n", err)
				continue
			}
			logData.startTime = revisedTime

			// Revise end time
			timeString = logData.endTime.Format("15:04")
			dateString = revisedDate.Format("01/02/06") + " " + timeString
			revisedTime, err = time.Parse("01/02/06 15:04", dateString)
			if err != nil {
				fmt.Printf("Invalid End Time: %v\n", err)
				continue
			}
			logData.endTime = revisedTime
		case 's':
			getLogTime("Start", logData.startTime, logData, scanner)
			recalculatedDuration, err := logData.calculateDuration()
			if err != nil {
				fmt.Printf("Duration could not be calculated: %v\n", err)
				time.Sleep(1750 * time.Millisecond)
				return
			}
			logData.duration = recalculatedDuration
		case 'e':
			getLogTime("End", logData.endTime, logData, scanner)
			recalculatedDuration, err := logData.calculateDuration()
			if err != nil {
				fmt.Printf("Duration could not be calculated: %v\n", err)
				time.Sleep(1750 * time.Millisecond)
				return
			}
			logData.duration = recalculatedDuration
		case 'c':
			getCategory(logData, scanner)
		case 'd':
			getDescription(logData, scanner)
		case 'q':
			return
		default:
			fmt.Printf("%s is invalid. Please enter again.", line)
			time.Sleep(1750 * time.Millisecond)
		}
	}
}
