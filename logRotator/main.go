package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	counter := 0
	const oneMB int64 = 1000000

	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error while creating file: %v", err)
	}
	defer file.Close()

	for {

		if counter == 10 {
			file.Close()
			newFilename := fmt.Sprintf("logs_%s.log", time.Now().Format("2006-01-02_15-04-05"))
			err := os.Rename("logs.log", newFilename)
			if err != nil {
				fmt.Printf("Error while renaming: %v", err)
			}
			break
		}

		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}

		if fileInfo.Size() >= oneMB {
			file.Close()
			// err := os.Rename("logs.log", "logs"+time.Now().Format("2006-01-02_15:04:05")+".log")
			newFilename := fmt.Sprintf("logs_%s.log", time.Now().Format("2006-01-02_15-04-05.000000"))
			err := os.Rename("logs.log", newFilename)
			if err != nil {
				fmt.Printf("Error while renaming: %v", err)
			}
			file, err = os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("error while creating file: %v", err)
			}
			counter++
		}

		log.SetOutput(file)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("This is an info log message.")
		log.Printf("User %d successfully logged in.\n", 42)

	}

}
