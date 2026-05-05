package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/dgraph-io/badger/v4"
)

func main() {
	// Initialize BadgerDB
	db, err := repository.InitBadgerDB("./data/badger")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// New password to set
	newPassword := "admin123"
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New hash: %s\n", hash)

	// Find the user
	var user map[string]interface{}
	err = db.View(func(txn *badger.Txn) error {
		// Try to find user by ID 1 first
		key := []byte("user:1")
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})

	if err != nil {
		fmt.Printf("Error finding user: %v\n", err)
		return
	}

	fmt.Printf("Found user: %v\n", user)

	// Update the password
	err = db.Update(func(txn *badger.Txn) error {
		user["passwordHash"] = hash
		user["updatedAt"] = time.Now()

		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return txn.Set([]byte("user:1"), data)
	})

	if err != nil {
		fmt.Printf("Error updating: %v\n", err)
		return
	}

	fmt.Println("✅ Password reset successfully!")
	fmt.Printf("Username: %s\n", user["username"])
	fmt.Printf("New password: %s\n", newPassword)
}
