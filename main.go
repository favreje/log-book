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

	_ "modernc.org/sqlite"
)

type LogData struct {
	projectId   int
	startTime   time.Time
	endTime     time.Time
	duration    time.Duration
	category    string
	description string
}

func (l LogData) calculateDuration() (time.Duration, error) {
	if l.startTime.IsZero() {
		return 0, fmt.Errorf("Start Time is not set")
	}
	if l.endTime.IsZero() {
		return 0, fmt.Errorf("End Time is not set")
	}
	if l.endTime.Before(l.startTime) {
		return 0, fmt.Errorf("End Time is before Start Time")
	}
	return l.endTime.Sub(l.startTime), nil
}

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

func getUserData(logData *LogData) {
	scanner := bufio.NewScanner(os.Stdin)

	// Project ID
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

	// Input date for start time and end time
	var inputDate time.Time
	for {
		fmt.Print("Log Date (MM/DD/YY): ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		parsedDate, err := time.Parse("01/02/06", line)
		if err != nil {
			fmt.Printf("Invalid Input Date: %v\n", err)
			continue
		}
		inputDate = parsedDate
		break
	}

	// Start time
	for {
		fmt.Print("Start Time (HH:MM): ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		dateString := inputDate.Format("01/02/06") + " " + line
		layout := "01/02/06 15:04"
		startTime, err := time.Parse(layout, dateString)
		if err != nil {
			fmt.Printf("Invalid Start Time: %v\n", err)
			continue
		}
		logData.startTime = startTime
		break
	}

	// End time
	for {
		fmt.Print("End Time (HH:MM): ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		dateString := inputDate.Format("01/02/06") + " " + line
		layout := "01/02/06 15:04"
		endTime, err := time.Parse(layout, dateString)
		if err != nil {
			fmt.Printf("Invalid End Time: %v\n", err)
			continue
		}
		logData.endTime = endTime
		break
	}

	// Duration
	duration, err := logData.calculateDuration()
	if err != nil {
		fmt.Printf("Duration could not be calculated: %v\n", err)
	}
	logData.duration = duration

	// Category
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

	// Description
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

func main() {
	db, err := sql.Open("sqlite", "project_log.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logData := LogData{}

	// Testing getUserData func

	getUserData(&logData)
	projectDesc, err := getProjectDesc(db, logData.projectId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No description matches Project Code:", logData.projectId)
			return
		}
		log.Fatal(err)
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(strings.Repeat(" ", 30), "Log Entry")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Printf("%-20s %02d \n", "Project ID:", logData.projectId)
	fmt.Printf("%-20s %s \n", "Project Description:", projectDesc)
	timeLayout := "01/02/06 15:04"
	fmt.Printf(
		"Start time: %23s \nEnd time: %25s \n",
		logData.startTime.Format(timeLayout),
		logData.endTime.Format(timeLayout),
	)

	fmt.Printf(
		"%-20s %2.2f hrs\n",
		"Duration:",
		logData.duration.Hours(),
	)

	fmt.Printf("%-20s %s\n", "Category:", logData.category)
	fmt.Printf("%-20s %s\n", "Description:", logData.description)
}
