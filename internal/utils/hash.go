package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func HashPassword(password string) (string, error) {
	log.Printf("🔐 Hashing password with bcrypt (cost=%d)", bcryptCost)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		log.Printf("❌ Hash error: %v", err)
		return "", err
	}
	hash := string(bytes)
	log.Printf("✅ Password hashed successfully, hash length: %d", len(hash))
	return hash, nil
}

func CheckPasswordHash(password, hash string) bool {
	log.Printf("🔐 Verifying password with bcrypt")
	
	if len(hash) == 0 {
		log.Printf("❌ Empty hash provided")
		return false
	}
	
	// Check if it's a valid bcrypt hash
	if len(hash) < 3 || (hash[:3] != "$2a" && hash[:3] != "$2b" && hash[:3] != "$2y") {
		log.Printf("❌ Invalid hash format: %s...", hash[:min(20, len(hash))])
		return false
	}
	
	log.Printf("   Hash prefix: %s", hash[:3])
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("❌ Password verification failed: %v", err)
		return false
	}
	
	log.Printf("✅ Password verified successfully")
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
