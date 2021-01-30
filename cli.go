package keeper

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func WaitForCommand(store Store) {
	for {
		var command string
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("> ")
		if scanner.Scan() {
			command = scanner.Text()
		}

		fmt.Println(command)

		switch command {
		case string(Add):
			handleAdd(store, scanner)
		case string(Show):
			handleShow(store)
		case string(Exit):
			os.Exit(0)
		}
	}
}

type commandType string

const (
	Add  commandType = "add"
	Show             = "show"
	Exit             = "exit"
)

func printAndScan(msg string, scanner *bufio.Scanner) string {
	fmt.Print(msg)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func handleAdd(store Store, scanner *bufio.Scanner) {
	book := BookEntry{}

	title := printAndScan("Title: ", scanner)
	book.Title = title

	author := printAndScan("Author: ", scanner)
	book.Author = author

	date := time.Now()
	for {
		startDate := printAndScan("Start Date (leave empty for the current day): ", scanner)
		if startDate != "" {
			if d, err := time.Parse("02-05-2006", startDate); err != nil {
				fmt.Printf("Couldn't parse the date: %s. Please try again.\n", startDate)
			} else {
				date = d
				break
			}
		} else {
			break
		}
	}
	book.DateStart = date

	for {
		endDate := printAndScan("End Date: ", scanner)
		if d, err := time.Parse("02-05-2006", endDate); err != nil {
			fmt.Printf("Couldn't parse the date: %s. Please try again.\n", endDate)
		} else {
			date = d
			break
		}
	}
	book.DateEnd = date

	for {
		state := BookState(printAndScan("Reading State: ", scanner))
		if err := state.IsValid(); err == nil {
			book.State = BookState(state)
			break
		} else {
			fmt.Printf("%s. Please try again.\n", err)
		}
	}

	if err := store.AddBookEntry(book); err != nil {
		fmt.Printf("Failed to add a book to the collection\n%s\n", err)
	}
}

func handleShow(store Store) {
	if err := store.PrintBookEntries(); err != nil {
		fmt.Println(err)
	}
}
