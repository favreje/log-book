package main

import (
	"fmt"
	"sort"
	"strings"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func displayUserInput(logData *LogData, projectsMap map[int]string, inputState *InputState) {
	projectDesc, ok := projectsMap[logData.projectId]
	if !ok {
		projectDesc = strings.Repeat("*", 15)
	}

	clearScreen()
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(strings.Repeat(" ", 30), "Log Entry")
	fmt.Println(strings.Repeat("-", 80))

	if logData.projectId == 0 {
		fmt.Printf("%-20s %s \n", "Project ID:", "**")
	} else {
		fmt.Printf("%-20s %02d \n", "Project ID:", logData.projectId)
	}
	fmt.Printf("%-20s %s \n", "Project Description:", projectDesc)

	// Start time display
	timeLayout := "01/02/06 15:04"
	if inputState.dateEntered && inputState.startTimeEntered {
		fmt.Printf("Start time: %23s\n", logData.startTime.Format(timeLayout))
	} else if inputState.dateEntered {
		fmt.Printf("Start time: %17s --:--\n", inputState.baseDate.Format("01/02/06"))
	} else {
		fmt.Printf("Start time: %23s\n", "MM/DD/YY HH:MM")
	}

	// End time display
	if inputState.dateEntered && inputState.endTimeEntered {
		fmt.Printf("End time: %25s\n", logData.endTime.Format(timeLayout))
	} else if inputState.dateEntered {
		fmt.Printf("End time: %19s --:--\n", inputState.baseDate.Format("01/02/06"))
	} else {
		fmt.Printf("End time: %25s\n", "MM/DD/YY HH:MM")
	}

	// Display duration if we have start and end times
	if inputState.startTimeEntered && inputState.endTimeEntered {
		fmt.Printf("%-20s %2.2f hrs\n", "Duration:", logData.duration.Hours())
	} else {
		fmt.Printf("%-20s %s\n", "Duration:", "--.-- hrs")
	}

	if logData.category == "" {
		fmt.Printf("%-20s %s\n", "Category:", strings.Repeat("*", 15))
	} else {
		fmt.Printf("%-20s %s\n", "Category:", logData.category)
	}

	if logData.description == "" {
		fmt.Printf("%-20s %s\n", "Description:", strings.Repeat("*", 25))
	} else {
		fmt.Printf("%-20s %s\n", "Description:", logData.description)
	}

	fmt.Println(strings.Repeat("-", 80))
}

func displayMainMenu() {
	clearScreen()
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println(strings.Repeat(" ", 10), "MAIN MENU")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("(I)nput log entry")
	fmt.Println("(D)isplay log report")
	fmt.Println("(E)xit")
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

	// Calculate max length of category and description fields
	// Description field should be at least 8 characters to accomodate the header
	categoryLen := 0
	descLen := 0
	for _, record := range logRecords {
		thisCatLen := len(record.category)
		if thisCatLen > categoryLen {
			categoryLen = thisCatLen
		}
		if categoryLen < 8 {
			categoryLen = 8
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
			"%02d %*s %s %s %s %s %4.2f %*s %s\n",
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
	hrsColumnPadding := 26 + projectDescLen
	fmt.Printf("%*s%9.2f\n\n", -hrsColumnPadding, "Total Hours:", durationTotal)
}

func displayProjectList(projectsMap map[int]string) {
	clearScreen()
	columnLen := 65
	header := "Available Projects"
	padding := (columnLen - len(header)) / 2
	fmt.Printf("%*s%s\n", padding, "", header)
	fmt.Println(strings.Repeat("-", columnLen))
	fmt.Println("ID  Project")
	fmt.Println(strings.Repeat("-", columnLen))

	keys := []int{}
	for key := range projectsMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for key := range keys {
		fmt.Printf("%02d %s\n", key+1, projectsMap[key+1])
	}

	fmt.Println(strings.Repeat("-", columnLen))
}
