package main

import (
	"fmt"

	"github.com/piquette/finance-go/quote"
)

func main() {

	q, err := quote.Get("AAPL")
	if err != nil {
		// Uh-oh.
		panic(err)
	}

	// Success!
	fmt.Println(q)
}
