package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type ResultRecord struct {
	ProductId      string
	Quantity       int64
	Category       string
	Price          float64
	IsDiscontinued bool
}

type Record struct {
	Index     int    `json:"index"`
	Timestamp string `json:"timestamp"`
	ProductId string `json:"product_id"`
	EventType string `json:"event_type"`
	Data      string `json:"data"`
}

type Data struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	InitialPrice string  `json:"initial_price"`
	Quantity     int64   `json:"quantity"`
	Old          float64 `json:"old"`
	New          float64 `json:"new"`
}

func main() {

	resultArr := []ResultRecord{}

	inputFile, err := os.Open("logs.jsonl")
	if err != nil {
		fmt.Printf("error while opening file: %v", err)
	}

	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		var record Record
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			fmt.Printf("error while record unmarshal: %v", err)
		}
		// fmt.Println(record)

		var data Data
		err = json.Unmarshal([]byte(record.Data), &data)
		if err != nil {
			fmt.Printf("error while data unmarshal: %v", err)
		}

		switch record.EventType {
		case "ITEM_CREATED":
			initialPrice, _ := strconv.ParseFloat(data.InitialPrice, 64)
			resultArr = append(resultArr, ResultRecord{record.ProductId, 0, data.Category, initialPrice, false})
		case "DISCONTINUED":
			updateStatus(record.ProductId, &resultArr)
		case "STOCK_REDUCED":
			updateStock(record.ProductId, -data.Quantity, &resultArr)
		case "STOCK_ADDED":
			updateStock(record.ProductId, data.Quantity, &resultArr)
		case "PRICE_CHANGED":
			updatePrice(record.ProductId, data.New, &resultArr)
		}

	}

	fmt.Println(resultArr)

	slices.SortFunc(resultArr, func(a, b ResultRecord) int {
		return cmp.Compare(a.ProductId, b.ProductId)
	})

	outputFiledata, err := json.MarshalIndent(resultArr, "", "	")
	if err != nil {
		fmt.Printf("error in marshal indent: %v", err)
	}

	err = os.WriteFile("output.json", outputFiledata, 0644)
	if err != nil {
		fmt.Printf("error in creating output file: %v", err)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error while scan: %v", err)
	}

}

func updateStock(productId string, quantity int64, resultRecords *[]ResultRecord) {
	for index, val := range *resultRecords {
		if productId == val.ProductId {
			(*resultRecords)[index].Quantity += quantity
		}
	}
}

func updatePrice(productId string, newPrice float64, resultRecords *[]ResultRecord) {
	for index, val := range *resultRecords {
		if productId == val.ProductId {
			(*resultRecords)[index].Price = newPrice
		}
	}
}

func updateStatus(productId string, resultRecords *[]ResultRecord) {

	for index, val := range *resultRecords {
		if val.ProductId == productId {
			(*resultRecords)[index].IsDiscontinued = true
		}
	}

}
