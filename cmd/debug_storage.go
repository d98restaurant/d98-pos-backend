package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
)

func main() {
	opts := badger.DefaultOptions("./data/badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// Find all users
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		
		prefix := []byte("user:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := string(it.Item().Key())
			// Only get main user records (not indices)
			if len(key) > 5 && key[:5] == "user:" && len(key) < 20 {
				item := it.Item()
				err := item.Value(func(val []byte) error {
					var user map[string]interface{}
					if err := json.Unmarshal(val, &user); err != nil {
						return err
					}
					fmt.Printf("User: %v\n", user["username"])
					fmt.Printf("  ID: %v\n", user["_id"])
					fmt.Printf("  Password stored: '%v'\n", user["passwordHash"])
					fmt.Printf("  Password length: %d\n", len(user["passwordHash"].(string)))
					fmt.Println("---")
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	
	if err != nil {
		log.Fatal(err)
	}
}
