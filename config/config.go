package config

import (
	"log"
	"time"
)

type Config struct {
	Port                   string
	Environment            string
	JWTSecret              string
	JWTExpiry              time.Duration
	BadgerDBPath           string
	RazorpayKeyID          string
	RazorpayKeySecret      string
	RazorpayWebhookSecret  string
	VAPIDPublicKey         string
	VAPIDPrivateKey        string
	FrontendURL            string
	WebSocketPingInterval  time.Duration
}

var AppConfig *Config

func LoadConfig() *Config {
	// JWT expiry: 168 hours = 7 days
	jwtExpiryHours := 168
	// WebSocket ping interval: 30 seconds
	wsPingIntervalSec := 30

	AppConfig = &Config{
		// Server Configuration
		Port:        "8080",
		Environment: "production",
		
		// JWT Authentication - Your existing secret
		JWTSecret: "bb8dc9d18888be7a54a38e9742d5c596cb539ba3855ec1d981a282880f61b1c8",
		JWTExpiry: time.Duration(jwtExpiryHours) * time.Hour,
		
		// BadgerDB - Embedded database (faster than MongoDB!)
		BadgerDBPath: "./data/badger",
		
		// Razorpay Payment Gateway - Your existing live keys
		RazorpayKeyID:        "rzp_live_RlSrpAntWBQZ7c",
		RazorpayKeySecret:    "5CLMqfS55PX63GFgONpHD24y",
		RazorpayWebhookSecret: "posd98",
		
		// VAPID Keys for Push Notifications
		VAPIDPublicKey:  "BLiNlzBDszn00bfL0w25zye0nz725AzWDHT-9RM-pksxzBbMWiO9W8cVTZdL1W5ouLIAcTMKyFkN4eXq4vAb_Kg",
		VAPIDPrivateKey: "NAUVzYdtLXUT17EH3rneMyn9AeKmrcEEIz6tJIiw62E",
		
		// Frontend URL - Your Cloudflare Worker URL
		FrontendURL: "https://pos-proxy.d98restaurant.workers.dev",
		
		// WebSocket Configuration
		WebSocketPingInterval: time.Duration(wsPingIntervalSec) * time.Second,
	}

	log.Printf("✅ Configuration loaded successfully")
	log.Printf("   📍 Environment: %s", AppConfig.Environment)
	log.Printf("   🗄️  Database: BadgerDB at %s", AppConfig.BadgerDBPath)
	log.Printf("   💳 Razorpay Key: %s", AppConfig.RazorpayKeyID[:10]+"...")
	log.Printf("   🌐 Frontend URL: %s", AppConfig.FrontendURL)

	return AppConfig
}
