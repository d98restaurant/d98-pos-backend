package main

import (
	"encoding/json"
	"fmt"
	"log"

	"pos-backend/internal/repository"
)

func main() {
	db, err := repository.InitBadgerDB("./data/badger")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Try to find user by username
	var userID string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("user:username:demouser"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			userID = string(val)
			return nil
		})
	})

	if err != nil {
		fmt.Printf("User not found: %v\n", err)
		return
	}

	fmt.Printf("User ID found: %s\n", userID)

	// Get user data
	var user map[string]interface{}
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("user:" + userID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})

	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}

	fmt.Printf("User data:\n")
	fmt.Printf("  Username: %v\n", user["username"])
	fmt.Printf("  PasswordHash length: %d\n", len(user["passwordHash"].(string)))
	fmt.Printf("  PasswordHash preview: %s...\n", user["passwordHash"].(string)[:20])
}
