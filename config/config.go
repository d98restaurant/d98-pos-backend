package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port         string
	Environment  string
	JWTSecret    string
	JWTExpiry    time.Duration

	// Database
	MongoURI     string
	MongoDBName  string

	// Razorpay
	RazorpayKeyID     string
	RazorpayKeySecret string
	RazorpayWebhookSecret string

	// VAPID for Notifications
	VAPIDPublicKey  string
	VAPIDPrivateKey string

	// Frontend URL
	FrontendURL  string

	// WebSocket
	WebSocketPingInterval time.Duration
}

var AppConfig *Config

func LoadConfig() *Config {
	// Load .env file if exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "168"))
	wsPingIntervalSec, _ := strconv.Atoi(getEnv("WS_PING_INTERVAL_SEC", "30"))

	AppConfig = &Config{
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		JWTSecret:    getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpiry:    time.Duration(jwtExpiryHours) * time.Hour,
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:  getEnv("MONGO_DB_NAME", "pos_system"),
		RazorpayKeyID:     getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret: getEnv("RAZORPAY_KEY_SECRET", ""),
		RazorpayWebhookSecret: getEnv("RAZORPAY_WEBHOOK_SECRET", ""),
		VAPIDPublicKey:    getEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey:   getEnv("VAPID_PRIVATE_KEY", ""),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3000"),
		WebSocketPingInterval: time.Duration(wsPingIntervalSec) * time.Second,
	}

	// Log configuration status (without sensitive data)
	log.Printf("✅ Configuration loaded:")
	log.Printf("   Port: %s", AppConfig.Port)
	log.Printf("   Environment: %s", AppConfig.Environment)
	log.Printf("   MongoDB URI: %s (masked)", maskString(AppConfig.MongoURI, 20))
	log.Printf("   Razorpay Key ID: %s", maskString(AppConfig.RazorpayKeyID, 10))
	log.Printf("   Frontend URL: %s", AppConfig.FrontendURL)

	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskString(s string, showLen int) string {
	if len(s) <= showLen {
		return "***"
	}
	return s[:showLen] + "..."
}