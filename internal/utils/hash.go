package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func HashPassword(password string) (string, error) {
	log.Printf("🔐 Hashing password (length: %d)", len(password))
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		log.Printf("❌ Hash error: %v", err)
		return "", err
	}
	log.Printf("✅ Password hashed successfully")
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	log.Printf("🔐 Verifying password...")
	
	if len(hash) == 0 {
		log.Printf("❌ Empty hash provided")
		return false
	}
	
	// Check if the hash looks like a bcrypt hash (starts with $2a$ or $2b$)
	if len(hash) < 3 || (hash[:3] != "$2a" && hash[:3] != "$2b") {
		log.Printf("❌ Invalid hash format (not bcrypt): %s...", hash[:20])
		return false
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("❌ Hash verification failed: %v", err)
		return false
	}
	log.Printf("✅ Hash verification successful")
	return true
}
