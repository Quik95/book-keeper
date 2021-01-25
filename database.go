package keeper

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func DumpDBContents(db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Printf("Bucket name: %s\n", string(name))
			b.ForEach(func(key, value []byte) error {
				fmt.Printf("Key: %s\nValue: %s", string(key), string(value))
				return nil
			})
			fmt.Print("\n\n----------\n\n")
			return nil
		})
		return nil
	})
}
