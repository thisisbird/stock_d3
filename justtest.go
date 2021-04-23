package main

import (
	"fmt"
	"time"
)

func main() {

	b := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
	c := b.Weekday()
	final := 0//結算日
	if c == 0 {
		final = 18
	}
	if c == 1 {
		final = 17
	}
	if c == 2 {
		final = 16
	}
	if c == 3 {
		final = 15
	}
	if c == 4 {
		final = 21
	}
	if c == 5 {
		final = 20
	}
	if c == 6 {
		final = 19
	}
	fmt.Println(final)

}

func isNextMonth(t time.Time) string {

	return "123"
}
