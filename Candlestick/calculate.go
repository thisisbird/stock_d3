package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	readCSV("new.csv")
}

func readCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE //開啟檔案的選項
	file2, err := os.OpenFile("calculate.csv", options, os.FileMode(0600))
	check(err)
	// i := 0
	for scanner.Scan() {
		_, err = fmt.Fprintln(file2, scanner.Text())
	}
	err = file.Close()
	check(err)
	err = file2.Close()
	check(err)

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
