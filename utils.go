package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func parseTimeFromDb(timeStr string) (time.Time, error) {
	format := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func confirmedSelection(msg string) bool {
	fmt.Printf("%s (Y/n): ", msg)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	line := strings.TrimSpace(scanner.Text())
	if line == "" {
		return true
	}
	lowerline := strings.ToLower(line)
	char, _ := utf8.DecodeRuneInString(lowerline)
	return char == 'y'
}

func getUserInput(prompt string) (string, bool) {
	fmt.Printf("%s ('Q' to Quit): ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatal()
		}
		return "", false
	}
	line := strings.TrimSpace(scanner.Text())
	lowerline := strings.ToLower(line)
	char, _ := utf8.DecodeRuneInString(lowerline)
	if char == 'q' {
		return "", false
	}
	return line, true
}
