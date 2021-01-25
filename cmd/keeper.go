package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	kp "github.com/Quik95/book-keeper"
	"github.com/boltdb/bolt"
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
	db_location, err := filepath.Abs(location)
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(db_location, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	kp.DumpDBContents(db)
}
