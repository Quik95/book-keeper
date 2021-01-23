package main

import (
	"fmt"
	"os"
)

func main() {
	var location string
	if len(os.Args) == 2 {
		location = os.Args[1]
	} else {
		fmt.Println("Please provide the database location")
		os.Exit(1)
	}
	fmt.Printf("Datebase location: %s", location)
}
