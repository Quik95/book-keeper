package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	kp "github.com/Quik95/book-keeper"
)

func main() {
	var location string
	if len(os.Args) == 2 {
		location = os.Args[1]
	} else {
		log.Fatal("Please provide the database location")
	}

	// when given a directory path automically set the database name
	file, err := os.Stat(location)
	switch {
	case err != nil:
		log.Fatalf("%s is not a valid path.", location)
	case file.IsDir():
		location = path.Join(location, "books.db")
	default:
		if !strings.HasSuffix(location, ".db") {
			log.Fatalf("%s uses an invalid extension", location)
		}
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

	kp.WaitForCommand(store)
}
