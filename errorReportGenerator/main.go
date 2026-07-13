package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"
)

type KeyVal struct {
	Key    string
	ValArr []string
}

func main() {

	inputFile, err := os.Open("application.log")
	if err != nil {
		fmt.Printf("error while opening file: %v\n", err)
	}

	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Printf("error while create output file: %v", err)
	}

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	currentRequestId := ""

	errorArr := []KeyVal{}
	errorMap := make(map[string][]string)
	statusMap := make(map[string]int)
	uniqueErrorMap := make(map[string]int)
	nullPointerExceptionArr := []string{}
	ioExceptionArr := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lineArr := strings.Fields(line)
		line = strings.Join(lineArr, " ")
		requestId, isLineValid := validateLogLine(line)
		if isLineValid {
			currentRequestId = requestId
			statusMap[lineArr[2]] = statusMap[lineArr[2]] + 1
		} else {
			errorMap[currentRequestId] = append(errorMap[currentRequestId], line)
		}
	}

	for k, v := range errorMap {
		errorArr = append(errorArr, KeyVal{k, v})
	}

	slices.SortFunc(errorArr, func(a, b KeyVal) int {
		return cmp.Compare(a.Key, b.Key)
	})

	fmt.Fprintf(writer, "Status Code count: \n")
	for k, v := range statusMap {
		fmt.Fprintf(writer, "Response Code: %v, Count: %v\n", k, v)
	}

	for _, v := range errorArr {
		uniqueErrorMap[strings.TrimSpace(v.ValArr[0])] = uniqueErrorMap[v.ValArr[0]] + 1

		if strings.Contains(v.ValArr[0], "NullPointerException") {
			nullPointerExceptionArr = append(nullPointerExceptionArr, v.ValArr[1])
		}

		if strings.Contains(v.ValArr[0], "IOException") {
			ioExceptionArr = append(ioExceptionArr, v.ValArr[1])
		}

	}

	fmt.Fprintf(writer, "\n\nUnique Error Codes: \n")
	for key, value := range uniqueErrorMap {
		fmt.Fprintf(writer, "%v: %v\n", key, value)
	}

	fmt.Fprintf(writer, "\n\nNull Pointer Exceptions: \n")
	for _, value := range nullPointerExceptionArr {
		re := regexp.MustCompile(`\((.*?)\)`)
		match := re.FindStringSubmatch(value)
		fmt.Fprintf(writer, "%v\n", match[1])
	}

	fmt.Fprintf(writer, "\n\nIO Exceptions: \n")
	for _, value := range ioExceptionArr {
		re := regexp.MustCompile(`\((.*?)\)`)
		match := re.FindStringSubmatch(value)
		fmt.Fprintf(writer, "%v\n", match[1])
	}

	if err = scanner.Err(); err != nil {
		fmt.Printf("error while scanning: %v\n", err)
	}

	if err = writer.Flush(); err != nil {
		fmt.Printf("error while writer flush: %v", err)
	}

}

func validateLogLine(line string) (string, bool) {

	lineArr := strings.Split(line, " ")
	// fmt.Println(lineArr)

	layout := "2006-01-02 15:04:05"
	requestId := ""

	if len(lineArr) > 5 {
		_, err := time.Parse(layout, lineArr[0]+" "+lineArr[1])
		if err != nil {
			fmt.Println(err)
			return "", false
		}
		if !(lineArr[2] == "INFO" || lineArr[2] == "DEBUG" || lineArr[2] == "ERROR" || lineArr[2] == "WARN" || lineArr[2] == "FATAL") {
			return "", false
		}
		requestArr := strings.Split(lineArr[4], "=")
		requestId = requestArr[1]
	} else {
		return "", false
	}

	return requestId, true
}
