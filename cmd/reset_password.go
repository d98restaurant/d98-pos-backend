package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Open BadgerDB
	opts := badger.DefaultOptions("./data/badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// New password
	newPassword := "admin123"
	
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	passwordHash := string(hash)
	
	fmt.Printf("New password hash: %s\n", passwordHash[:30]+"...")
	
	// Find and update user
	err = db.Update(func(txn *badger.Txn) error {
		// Find user by username
		userIDKey := []byte("user:username:demouser")
		item, err := txn.Get(userIDKey)
		if err != nil {
			return fmt.Errorf("user not found: %v", err)
		}
		
		var userID string
		err = item.Value(func(val []byte) error {
			userID = string(val)
			return nil
		})
		if err != nil {
			return err
		}
		
		fmt.Printf("Found user ID: %s\n", userID)
		
		// Get user data
		userKey := []byte("user:" + userID)
		item, err = txn.Get(userKey)
		if err != nil {
			return err
		}
		
		var user map[string]interface{}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
		if err != nil {
			return err
		}
		
		// Update password
		user["passwordHash"] = passwordHash
		
		// Save back
		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}
		
		return txn.Set(userKey, userData)
	})
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("✅ Password reset successfully!")
	fmt.Println("Username: demouser")
	fmt.Println("New password: admin123")
}
