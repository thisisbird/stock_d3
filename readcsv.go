package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("2020_fut.csv")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())

		sli := strings.Split(scanner.Text(), ",")
		fmt.Println(sli[2:5])
		panic(123)
		// writeTxt(scanner.Text())
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

}

func writeTxt(data string) {
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE //開啟檔案的選項
	file, err := os.OpenFile("signatures.txt", options, os.FileMode(0600))
	check(err)
	_, err = fmt.Fprintln(file, data)
	check(err)
	err = file.Close()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
