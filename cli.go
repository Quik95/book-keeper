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

	setBookTitle(&book, scanner)
	setBookAuthor(&book, scanner)
	setBookDate("start", &book, scanner)
	setBookDate("end", &book, scanner)
	setBookState(&book, scanner)

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
	book, err := store.GetBookWithIndex(idx)
	if err != nil {
		fmt.Printf("A book with the index %d doesn't exists.\n", idx)
	}

	if err := store.DeleteBookEntry(itob(book.ID)); err != nil {
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
	updatedBook := oldBook
	var _ = askUntilTheConditionHasBeenMet(msg, scanner, func(s string) (interface{}, error) {
		switch s {
		case "1":
			setBookTitle(&updatedBook, scanner)
			return nil, nil
		case "2":
			setBookAuthor(&updatedBook, scanner)
			return nil, nil
		case "3":
			setBookDate("start", &updatedBook, scanner)
			return nil, nil
		case "4":
			setBookDate("end", &updatedBook, scanner)
			return nil, nil
		case "5":
			setBookState(&updatedBook, scanner)
			return nil, nil
		default:
			return nil, fmt.Errorf("%s is not a valid book property. Please try again.\n", s)
		}
	})
	return updatedBook
}

func setBookTitle(book *BookEntry, scanner *bufio.Scanner) {
	newTitle := askUntilTheConditionHasBeenMet("Title: ", scanner, func(s string) (interface{}, error) {
		if len(s) > 0 {
			return s, nil
		} else {
			return "", fmt.Errorf("Cannot use an empty string as a book title. Please try again.\n")
		}
	}).(string)
	book.Title = newTitle
}

func setBookAuthor(book *BookEntry, scanner *bufio.Scanner) {
	setBookState := askUntilTheConditionHasBeenMet("Author: ", scanner, func(s string) (interface{}, error) {
		if len(s) > 0 {
			return s, nil
		} else {
			return "", fmt.Errorf("Cannot use an empty string as a book author. Please try again.\n")
		}
	}).(string)
	book.Author = setBookState
}

func getBookState(scanner *bufio.Scanner) BookState {
	msg := "Select a book state:\n1. Reading\n2. Finished\n3. Dropped\n4. Suspended\n5. Re reading\nChoice: "
	st := askUntilTheConditionHasBeenMet(msg, scanner, func(s string) (interface{}, error) {
		switch s {
		case "1":
			return Reading, nil
		case "2":
			return Finished, nil
		case "3":
			return Dropped, nil
		case "4":
			return Suspended, nil
		case "5":
			return ReRead, nil
		default:
			return nil, fmt.Errorf("%s is not a valid book state. Please select a number from 1 to 5.\n", s)
		}
	}).(BookState)
	return st
}

func setBookDate(dateType string, book *BookEntry, scanner *bufio.Scanner) {
	var msg string
	if dateType == "start" {
		msg = "Start Date: "
	} else {
		msg = "End Date: "
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
		book.DateStart = newDate
	} else {
		book.DateEnd = newDate
	}
}

func setBookState(book *BookEntry, scanner *bufio.Scanner) {
	newState := getBookState(scanner)
	book.State = newState
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
		"update: updates a selected book\n",
		"help: prints this help\n",
		"exit: exits from the program\n")
}
