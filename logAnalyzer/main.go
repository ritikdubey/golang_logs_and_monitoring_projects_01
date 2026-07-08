package main

import (
	"bufio"
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type KeyVal struct {
	Key   string
	Value int
}

func main() {

	inputFile, err := os.Open("access.log")
	if err != nil {
		fmt.Printf("error while reading file: %v", err)
	}

	outputFile, err := os.Create("report.txt")
	if err != nil {
		fmt.Printf("error while creating file: %v", err)
	}

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	ipMap := make(map[string]int)
	ipMapArr := []KeyVal{}

	urlMap := make(map[string]int)
	urlMapArr := []KeyVal{}

	statusCodeMap := make(map[string]int)
	statusCodeMapArr := []KeyVal{}

	httpMethodMap := make(map[string]int)
	httpMethodMapArr := []KeyVal{}

	invalidLogsArr := []string{}

	responseSizeSum := 0
	requestCount := 0
	largestResp := 0
	smallestResp := math.MaxInt64

	for scanner.Scan() {

		line := scanner.Text()
		lineArr := strings.Split(line, " ")

		if len(lineArr) >= 6 {
			ipMap[lineArr[0]] = ipMap[lineArr[0]] + 1
			urlMap[lineArr[3]] = urlMap[lineArr[3]] + 1
			statusCodeMap[lineArr[4]] = statusCodeMap[lineArr[4]] + 1
			httpMethodMap[lineArr[2]] = httpMethodMap[lineArr[2]] + 1
			respSize, _ := strconv.Atoi(lineArr[5])
			responseSizeSum += respSize
			requestCount += 1
			if respSize > largestResp {
				largestResp = respSize
			}
			if respSize < smallestResp {
				smallestResp = respSize
			}

			layout := "2006-01-02T15:04:05"
			reqTime, err := time.Parse(layout, strings.TrimSpace(lineArr[1]))

			fmt.Println(reqTime)
			fmt.Println(err)

		} else {
			invalidLogsArr = append(invalidLogsArr, line)
		}
	}

	for key, val := range ipMap {
		ipMapArr = append(ipMapArr, KeyVal{key, val})
	}

	for key, val := range urlMap {
		urlMapArr = append(urlMapArr, KeyVal{key, val})
	}

	for key, val := range statusCodeMap {
		statusCodeMapArr = append(statusCodeMapArr, KeyVal{key, val})
	}

	for key, val := range httpMethodMap {
		httpMethodMapArr = append(httpMethodMapArr, KeyVal{key, val})
	}

	slices.SortFunc(ipMapArr, func(a, b KeyVal) int {
		return cmp.Compare(b.Value, a.Value)
	})

	slices.SortFunc(urlMapArr, func(a, b KeyVal) int {
		return cmp.Compare(b.Value, a.Value)
	})

	slices.SortFunc(statusCodeMapArr, func(a, b KeyVal) int {
		return cmp.Compare(a.Key, b.Key)
	})

	slices.SortFunc(httpMethodMapArr, func(a, b KeyVal) int {
		return cmp.Compare(b.Value, a.Value)
	})

	fmt.Fprintf(writer, "Top 5 IP Addresses\n")
	for i := 0; i < 5; i++ {
		fmt.Fprintf(writer, "%v: %v\n", ipMapArr[i].Key, ipMapArr[i].Value)
	}

	fmt.Fprintf(writer, "\n\n")

	fmt.Fprintf(writer, "Top 10 URL\n")
	for i := 0; i < 10; i++ {
		fmt.Fprintf(writer, "%v: %v\n", urlMapArr[i].Key, urlMapArr[i].Value)
	}

	fmt.Fprintf(writer, "\n\n")

	fmt.Fprintf(writer, "HTTP Request Type Distribution\n")
	for _, v := range httpMethodMapArr {
		fmt.Fprintf(writer, "%v: %v\n", v.Key, v.Value)
	}

	fmt.Fprintf(writer, "\n\n")

	fmt.Fprintf(writer, "HTTP Status Code Distribution\n")
	for _, v := range statusCodeMapArr {
		fmt.Fprintf(writer, "%v: %v\n", v.Key, v.Value)
	}

	fmt.Fprintf(writer, "\n\n")
	fmt.Fprintf(writer, "Average Response Size: %v\n", responseSizeSum/requestCount)
	fmt.Fprintf(writer, "Total Requests: %v\n", requestCount)
	fmt.Fprintf(writer, "Total Response Size: %v\n", responseSizeSum)
	fmt.Fprintf(writer, "Largest Response: %v\n", largestResp)
	fmt.Fprintf(writer, "Smallest Response: %v\n", smallestResp)

	fmt.Fprintf(writer, "\n\n")

	fmt.Fprintf(writer, "Invalid Log Lines\n")
	for _, v := range invalidLogsArr {
		fmt.Fprintf(writer, "%v\n", v)
	}

	if err = scanner.Err(); err != nil {
		fmt.Printf("error in scan: %v", err)
	}

	if err := writer.Flush(); err != nil {
		fmt.Printf("error while flush: %v", err)
	}

}
