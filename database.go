package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/olekukonko/tablewriter"
)

// Store represents application data storage
type Store struct {
	db *bolt.DB
}

// LoadStore loads the Bolt database from given file path
func LoadStore(filepath string) (Store, error) {
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return Store{}, err
	}

	if err := setupDefaultBucket(db); err != nil {
		return Store{}, err
	}

	return Store{db}, nil
}

func setupDefaultBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("store")); err != nil {
			return err
		}
		return nil
	})
}

// Close calls close for each member of Store that needs to be closed
func (store Store) Close() error {
	if err := store.db.Close(); err != nil {
		return err
	}
	return nil
}

// DumpDBContents dumps the entire database contents to the console
func (store Store) DumpDBContents() {
	store.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Printf("Bucket name: %s\n----------\n", string(name))
			b.ForEach(func(key, value []byte) error {
				fmt.Printf("Key: %s\nValue: %s\n", string(key), string(value))
				return nil
			})
			fmt.Print("\n\n~~~~~~~~~~\n\n")
			return nil
		})
		return nil
	})
}

// PrintBookEntries prints books stored in the database in a friendly format
func (store Store) PrintBookEntries() error {
	return store.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("store"))
		if bkt == nil {
			return fmt.Errorf("Failed to retrieve the default store")
		}

		bookList := []BookEntry{}
		err := bkt.ForEach(func(_, value []byte) error {
			var book BookEntry
			err := json.Unmarshal(value, &book)
			if err != nil {
				return err
			}
			bookList = append(bookList, book)

			return nil
		})
		if err != nil {
			return err
		}

		sort.Slice(bookList, func(i, j int) bool {
			return bookList[i].DateStart.Before(bookList[j].DateStart)
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Index", "Title", "Author", "Start Date", "End Date", "Reading State"})
		form := "02 January 2006"
		for i, book := range bookList {
			table.Append(
				[]string{strconv.Itoa(i + 1), book.Title, book.Author, book.DateStart.Format(form), book.DateEnd.Format(form), string(book.State)},
			)
		}
		table.Render()

		return nil
	})
}

// AddBookEntry adds a book entry to the database
// ID has been assigned manually and will be overridden
func (store Store) AddBookEntry(be BookEntry) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("store"))
		if bkt == nil {
			return fmt.Errorf("Failed to retrieve the default store")
		}

		id, err := bkt.NextSequence()
		if err != nil {
			return err
		}
		be.ID = int(id)

		rawBytes, err := json.Marshal(be)
		if err != nil {
			return err
		}

		if err := bkt.Put(itob(be.ID), rawBytes); err != nil {
			return err
		}
		return nil
	})
}

// DeleteBookEntry removes a book entry with a given id from the database
func (store Store) DeleteBookEntry(bookID int) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("store"))
		if bkt == nil {
			return fmt.Errorf("Failed to retrieve the default store")
		}

		if err := bkt.Delete(itob(bookID)); err != nil {
			return err
		}
		return nil
	})
}

// itob returns an 8-byte big endian representation of v.
// taken from boltdb docs
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// BookEntry represents a single book with all it's parameters
// DateStart and DateEnd store timestamps when the user started
// and ended reading the book
// The first timestamp in DateStart corresponds to the first timestamps
// in DateEnd
type BookEntry struct {
	Title, Author      string
	DateStart, DateEnd time.Time
	State              BookState
	ID                 int
}

// BookState represents the state in which a book is currently in
type BookState string

const (
	// Reading indicates that the book is in the reading state
	Reading BookState = "reading"
	// Finished indicates that the book was finished
	Finished = "finished"
	// Dropped indicates that the book is in the dropped state
	Dropped = "dropped"
	// Suspended indicates that the user is no longer reading this book
	// but might return to it
	Suspended = "suspended"
)

// IsValid checks if a given instance of a BookState is valid
// if not returns an error
func (bs BookState) IsValid() error {
	switch bs {
	case Reading, Finished, Dropped, Suspended:
		return nil
	default:
		return fmt.Errorf("%s is not a valid book state", bs)
	}
}
