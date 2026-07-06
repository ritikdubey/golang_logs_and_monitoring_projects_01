package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
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

	scanner := bufio.NewScanner(inputFile)

	ipMap := make(map[string]int)
	ipMapArr := []KeyVal{}

	for scanner.Scan() {

		line := scanner.Text()
		lineArr := strings.Split(line, " ")

		if len(lineArr) >= 6 {
			ipMap[lineArr[0]] = ipMap[lineArr[0]] + 1
		}

	}

	//fmt.Println(ipMap)

	for key, val := range ipMap {
		ipMapArr = append(ipMapArr, KeyVal{key, val})
	}

	slices.SortFunc(ipMapArr, func(a, b KeyVal) int {
		return cmp.Compare(b.Value, a.Value)
	})

	fmt.Println(ipMapArr)

	if err = scanner.Err(); err != nil {
		fmt.Printf("error in scan: %v", err)
	}

}
