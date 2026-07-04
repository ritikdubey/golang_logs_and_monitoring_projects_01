package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	inputFile, err := os.Open("input.log")
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("result.txt")
	if err != nil {
		log.Fatalf("Error while creating file: %v", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	var loggingLevelArr []string
	var loggingLevel string
	fmt.Print("ENTER LOGGING LEVEL (INFO, DEBUG, WARN, ERROR): ")
	fmt.Scanln(&loggingLevel)
	if strings.Contains(loggingLevel, ",") {
		loggingLevelArr = strings.Split(loggingLevel, ",")
	} else {
		loggingLevelArr = []string{loggingLevel}
	}

	if !strings.Contains(strings.ToUpper(loggingLevel), "INFO") && !strings.Contains(strings.ToUpper(loggingLevel), "DEBUG") &&
		!strings.Contains(strings.ToUpper(loggingLevel), "WARN") && !strings.Contains(strings.ToUpper(loggingLevel), "ERROR") {
		fmt.Println("Please select correct logging level")
	}

	infoCount := 0
	warnCount := 0
	debugCount := 0
	errorCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		for _, loggingLevel1 := range loggingLevelArr {

			if strings.Contains(line, "INFO") {
				infoCount++
			} else if strings.Contains(line, "WARN") {
				warnCount++
			} else if strings.Contains(line, "DEBUG") {
				debugCount++
			} else if strings.Contains(line, "ERROR") {
				errorCount++
			} else {
				fmt.Printf("Invalid Log Line: %v", line)
			}

			if strings.Contains(line, strings.ToUpper(string(loggingLevel1))) {
				fmt.Fprintf(writer, "%v\n", line)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Error while reading")
	}

	if err = writer.Flush(); err != nil {
		fmt.Println("Error while flush write")
	}

	fmt.Printf("INFO: %v\nWARN: %v\nDEBUG: %v\nERROR: %v\n", infoCount, warnCount, debugCount, errorCount)

}
