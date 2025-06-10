package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func displayUserInput(db *sql.DB, logData *LogData) {
	projectDesc, err := getProjectDesc(db, logData.projectId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No description matches Project Code:", logData.projectId)
			return
		}
		log.Fatal(err)
	}

	clearScreen()
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

	fmt.Println(strings.Repeat("-", 80))
	fmt.Print("(W)rite log entry\n\n")
	fmt.Println("Edit: (P)roject ID | (S)tart time | (E)nd time | (C)ategory | (D)escription")
}
