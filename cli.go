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
		case string(Update):
			handleUpdate(store, scanner)
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
	Update commandType = "update"
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

func parseDateInput(input string) (time.Time, error) {
	if strings.HasPrefix(input, "???") {
		//user want's unspecified date
		return time.Time{}, nil
	} else if len(input) > 0 {
		// user provided a date, try to parse it
		if d, err := time.Parse(DateFormat, input); err != nil {
			return time.Time{}, err
		} else {
			// date has been parsed successfully, break from the loop
			return d, nil
		}
	}
	// user did not provide any date input, use todays date
	return time.Now(), nil
}

func handleAdd(store Store, scanner *bufio.Scanner) {
	book := BookEntry{}

	title := askUntilTheConditionHasBeenMet("Title: ", scanner, func(s string) (interface{}, error) {
		if len(s) <= 0 {
			return "", fmt.Errorf("Cannot use an empty string as a book title. Please try again.\n")
		}
		return s, nil
	}).(string)
	book.Title = title

	author := askUntilTheConditionHasBeenMet("Author: ", scanner, func(s string) (interface{}, error) {
		if len(s) <= 0 {
			return "", fmt.Errorf("Cannot use an empty string as a book author. Please try again.\n")
		}
		return s, nil
	}).(string)
	book.Author = author

	startDate := askUntilTheConditionHasBeenMet("Start Date (leave empty for the current day or ??? for and undefined date): ", scanner, func(s string) (interface{}, error) {
		date, err := parseDateInput(s)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse the date: %s. Please try again.\n", s)
		} else {
			return date, nil
		}
	}).(time.Time)
	book.DateStart = startDate

	endDate := askUntilTheConditionHasBeenMet("End Date (leave empty for the current day or ??? for and undefined date): ", scanner, func(s string) (interface{}, error) {
		date, err := parseDateInput(s)
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse the date: %s. Please try again.\n", s)
		} else {
			return date, nil
		}
	}).(time.Time)
	book.DateEnd = endDate

	state := askUntilTheConditionHasBeenMet("Reading State: ", scanner, func(s string) (interface{}, error) {
		st := BookState(s)
		if err := st.IsValid(); err == nil {
			return st, nil
		} else {
			return nil, fmt.Errorf("%s is not a valid reading state. Please try again.\n", s)
		}
	}).(BookState)
	book.State = state

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
	idx := getBookIdx(store, scanner)

	if err := store.DeleteBookEntry(idx); err != nil {
		fmt.Printf("Failed to delete this book.\n%s\n", err)
	}
}

func handleUpdate(store Store, scanner *bufio.Scanner) {
	bookIdx := getBookIdx(store, scanner)
	oldBookEntry, err := store.GetBookWithIndex(bookIdx)
	if err != nil {
		fmt.Println("Failed to retrieve the book")
		return
	}
	newBookEntry := updateBookProperty(oldBookEntry, scanner)

	if err := store.UpdateBookEntry(itob(oldBookEntry.ID), newBookEntry); err != nil {
		fmt.Println("Failed to update this book.")
	}
}

func getBookIdx(store Store, scanner *bufio.Scanner) int {
	maxIdx := store.GetNumberOfBookEntries()

	return askUntilTheConditionHasBeenMet("Please select the book index: ", scanner, func(s string) (interface{}, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("Cannot parse %s into a number. Please try again.\n", s)
		} else {
			if i > 0 && i <= maxIdx {
				return i, nil
			} else {
				return 0, fmt.Errorf("%d is not a valid book index. Please try again.\n", i)
			}
		}
	}).(int)
}

func updateBookProperty(oldBook BookEntry, scanner *bufio.Scanner) BookEntry {
	msg := "Please select the property you want to update:1. Title\n2. Author\n3. Date Start\n4. Date End\n5. Book State\nChoice: "
	return askUntilTheConditionHasBeenMet(msg, scanner, func(s string) (interface{}, error) {
		switch s {
		case "1":
			return updateBookTitle(oldBook, scanner), nil
		case "2":
			return updateBookAuthor(oldBook, scanner), nil
		case "3":
			return updateBookDate("start", oldBook, scanner), nil
		case "4":
			return updateBookDate("end", oldBook, scanner), nil
		case "5":
			return updateBookState(oldBook, scanner), nil
		default:
			return BookEntry{}, fmt.Errorf("%s is not a valid book property. Please try again.\n", s)
		}
	}).(BookEntry)
}

func updateBookTitle(oldBook BookEntry, scanner *bufio.Scanner) BookEntry {
	newTitle := askUntilTheConditionHasBeenMet("Please provide a new title: ", scanner, func(s string) (interface{}, error) {
		if len(s) > 0 {
			return s, nil
		} else {
			return "", fmt.Errorf("Cannot use an empty string as a book title. Please try again.\n")
		}
	}).(string)
	oldBook.Title = newTitle
	return oldBook
}

func updateBookAuthor(oldBook BookEntry, scanner *bufio.Scanner) BookEntry {
	newAuthor := askUntilTheConditionHasBeenMet("Please provide a new book author: ", scanner, func(s string) (interface{}, error) {
		if len(s) > 0 {
			return s, nil
		} else {
			return "", fmt.Errorf("Cannot use an empty string as a book author. Please try again.\n")
		}
	}).(string)
	oldBook.Author = newAuthor
	return oldBook
}

func updateBookDate(dateType string, oldBook BookEntry, scanner *bufio.Scanner) BookEntry {
	var msg string
	if dateType == "start" {
		msg = "Please provide a new start date: "
	} else {
		msg = "Please provide a new end date: "
	}

	newDate := askUntilTheConditionHasBeenMet(msg, scanner, func(s string) (interface{}, error) {
		d, err := parseDateInput(s)
		if err != nil {
			return "", fmt.Errorf("Couldn't parse date input. Please try again.\n")
		} else {
			return d, nil
		}
	}).(time.Time)
	if dateType == "start" {
		oldBook.DateStart = newDate
	} else {
		oldBook.DateEnd = newDate
	}
	return oldBook
}

func updateBookState(oldBook BookEntry, scanner *bufio.Scanner) BookEntry {
	newState := askUntilTheConditionHasBeenMet("Please provide a new reading state: ", scanner, func(s string) (interface{}, error) {
		bookState := BookState(s)
		if err := bookState.IsValid(); err != nil {
			return nil, fmt.Errorf("%s is not a valid reading state. Please try again.\n", s)
		}
		return bookState, nil
	}).(BookState)
	oldBook.State = newState
	return oldBook
}

func askUntilTheConditionHasBeenMet(msg string, scanner *bufio.Scanner, callback func(s string) (interface{}, error)) interface{} {
	for {
		input := printAndScan(msg, scanner)
		if r, err := callback(input); err == nil {
			return r
		} else {
			fmt.Print(err)
		}
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
