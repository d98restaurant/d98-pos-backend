package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10 // Use cost 10 for good balance

func HashPassword(password string) (string, error) {
	log.Printf("🔐 Hashing password (length: %d)", len(password))
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		log.Printf("❌ Hash error: %v", err)
		return "", err
	}
	log.Printf("✅ Password hashed successfully, hash: %s", string(bytes)[:20]+"...")
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	log.Printf("🔐 Verifying password (length: %d) against hash (length: %d)", len(password), len(hash))
	
	// Check if hash is empty or too short
	if len(hash) == 0 {
		log.Printf("❌ Empty hash provided")
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
