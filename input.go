package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
