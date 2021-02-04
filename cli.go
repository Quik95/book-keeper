package keeper

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

		switch command {
		case string(Add):
			handleAdd(store, scanner)
		case string(List):
			handleList(store)
		case string(Delete):
			handleDelete(store, scanner)
		case string(Help):
			handleHelp()
		case string(Exit):
			os.Exit(0)
		default:
			fmt.Printf("%s is not a valid command\n", command)
		}
	}
}

type commandType string

const (
	Add    commandType = "add"
	List   commandType = "list"
	Delete commandType = "delete"
	Help   commandType = "help"
	Exit   commandType = "exit"
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
		startDate := printAndScan("Start Date (leave empty for the current day or ??? for and undefined date): ", scanner)
		// if date starts with ??? use the date's zero value
		if strings.HasPrefix(startDate, "???") {
			date = time.Time{}
			break
		} else if len(startDate) > 0 {
			// if user provided a date try to parse it
			if d, err := time.Parse("02-05-2006", startDate); err != nil {
				fmt.Printf("Couldn't parse the date: %s. Please try again.\n", startDate)
			} else {
				// date has been parsed successfully, break from the loop
				date = d
				break
			}
		} else {
			// user did not provide any date input, use todays date
			break
		}
	}
	book.DateStart = date

	date = time.Now()
	for {
		endDate := printAndScan("End Date (leave empty for the current day or ??? for an undefined date): ", scanner)
		if strings.HasPrefix(endDate, "???") {
			date = time.Time{}
			break
		} else if len(endDate) > 0 {
			if d, err := time.Parse("02-05-2006", endDate); err != nil {
				fmt.Printf("Couldn't parse the date: %s. Please try again.\n", endDate)
			} else {
				date = d
				break
			}
		} else {
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

func handleList(store Store) {
	if err := store.PrintBookEntries(); err != nil {
		fmt.Println(err)
	}
}

func handleDelete(store Store, scanner *bufio.Scanner) {
	var idx int
	for {
		input := printAndScan("Select a book number to delete: ", scanner)
		i, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s is not a valid book number. Please try again.\n", input)
		} else {
			idx = i
			break
		}
	}

	if err := store.DeleteBookEntry(idx); err != nil {
		fmt.Printf("Failed to delete this book.\n%s\n", err)
	}
}

func handleHelp() {
	fmt.Print(
		"Available commands:\n",
		"list: list books in the database\n",
		"add: adds a book to the database\n",
		"delete: removes a book from the database\n",
		"exit: exits from the program\n")
}
