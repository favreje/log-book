package main

import (
	"fmt"
	"strings"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func displayUserInput(logData *LogData, projectsMap map[int]string) {
	projectDesc, ok := projectsMap[logData.projectId]
	if !ok {
		fmt.Println("No description matches Project Code:", logData.projectId)
		return
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
}

func displayMainMenu() {
	clearScreen()
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println(strings.Repeat(" ", 10), "MAIN MENU")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("1. (I)nput log entry")
	fmt.Println("2. (D)isplay log report")
	fmt.Println("3. (E)xit")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Print("Selection: ")
}

func reportByProject(logRecords []LogData, projectsMap map[int]string) {
	clearScreen()

	if len(logRecords) == 0 {
		fmt.Println("No log records found")
		return
	}
	projectDesc, ok := projectsMap[logRecords[0].projectId]
	if !ok {
		projectDesc = "Undetermined Project"
	}
	projectDescLen := len(projectDesc)

	// Calculate max length of category field
	categoryLen := 0
	descLen := 0
	for _, record := range logRecords {
		thisCatLen := len(record.category)
		if thisCatLen > categoryLen {
			categoryLen = thisCatLen
		}
		thisDescLen := len(record.description)
		if thisDescLen > descLen {
			descLen = thisDescLen
		}
	}
	fixedFieldTotalLength := 29
	columnDelimLen := 8
	columnLen := fixedFieldTotalLength + columnDelimLen + projectDescLen + categoryLen + descLen
	header := "LOG REPORT BY PROJECT ID"
	padding := (columnLen - len(header)) / 2

	fmt.Printf("%*s%s\n", padding, "", header)
	fmt.Println(strings.Repeat("-", columnLen))
	fmt.Printf(
		"%s %*s %s %-10s %s %-5s %-2s  %*s %s\n",
		"ID",
		-projectDescLen,
		"Project",
		"Day",
		"Date",
		"Start",
		"End",
		"Hrs",
		-categoryLen,
		"Category",
		"Description",
	)
	fmt.Println(strings.Repeat("-", columnLen))

	durationTotal := 0.0
	for _, record := range logRecords {
		dateStr := record.startTime.Format("01/02/2006")
		startStr := record.startTime.Format("15:04")
		weekDayStr := record.startTime.Format("Mon")
		endStr := record.endTime.Format("15:04")
		durationTotal += record.duration.Hours()

		fmt.Printf(
			"%2d %*s %s %s %s %s %4.2f %*s %s\n",
			record.projectId,
			-projectDescLen,
			projectDesc,
			weekDayStr,
			dateStr,
			startStr,
			endStr,
			record.duration.Hours(),
			-categoryLen,
			record.category,
			record.description,
		)
	}
	fmt.Println(strings.Repeat("-", columnLen))
	hrsColumnPadding := 30 + projectDescLen
	fmt.Printf("%*s%4.2f\n\n", -hrsColumnPadding, "Total Hours:", durationTotal)
}
