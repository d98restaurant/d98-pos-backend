package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	opts := badger.DefaultOptions("./data/badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	username := "finaluser"
	var userID string
	
	// Find user by username
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("user:username:" + username))
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
	
	fmt.Printf("User ID: %s\n", userID)
	
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
	
	passwordHash, ok := user["passwordHash"].(string)
	if !ok {
		fmt.Printf("PasswordHash not found or wrong type\n")
		return
	}
	
	fmt.Printf("Username: %v\n", user["username"])
	fmt.Printf("PasswordHash length: %d\n", len(passwordHash))
	fmt.Printf("PasswordHash preview: %s\n", passwordHash[:min(30, len(passwordHash))])
	fmt.Printf("Is bcrypt hash? %v\n", len(passwordHash) > 3 && (passwordHash[:3] == "$2a" || passwordHash[:3] == "$2b"))
	
	// Test bcrypt verification
	testPassword := "final123"
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(testPassword))
	fmt.Printf("bcrypt verification result: %v\n", err == nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
