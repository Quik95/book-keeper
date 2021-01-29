package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	kp "github.com/Quik95/book-keeper"
)

func main() {
	var location string
	if len(os.Args) == 2 {
		location = os.Args[1]
	} else {
		fmt.Println("Please provide the database location")
		os.Exit(1)
	}

	// when given a directory path automically set the database name
	if filepath.Ext(location) != ".db" {
		location = filepath.Join(location, "books.db")
	}
	// parse and validate the database path
	dbLocation, err := filepath.Abs(location)
	if err != nil {
		log.Fatal(err)
	}

	store, err := kp.LoadStore(dbLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if err := store.PrintBookEntries(); err != nil {
		log.Fatal(err)
	}
}
