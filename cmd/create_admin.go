package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type User struct {
	ID           string    `json:"_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	Role         string    `json:"role"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func main() {
	// Open BadgerDB
	opts := badger.DefaultOptions("./data/badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create admin user
	adminUser := User{
		ID:           "admin",
		Username:     "admin",
		Email:        "admin@pos.com",
		PasswordHash: "admin123", // Plain text password
		Role:         "admin",
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save admin user
	err = db.Update(func(txn *badger.Txn) error {
		// Save by ID
		data, err := json.Marshal(adminUser)
		if err != nil {
			return err
		}
		if err := txn.Set([]byte("user:admin"), data); err != nil {
			return err
		}
		// Save by username for lookup
		if err := txn.Set([]byte("user:username:admin"), []byte("admin")); err != nil {
			return err
		}
		// Save by email for lookup
		if err := txn.Set([]byte("user:email:admin@pos.com"), []byte("admin")); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Admin user created successfully!")
	fmt.Println("   Username: admin")
	fmt.Println("   Password: admin123")
	fmt.Println("   Role: admin")
}
