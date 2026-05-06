package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(frontendURL string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		frontendURL,
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"https://pos-backend-h2c3.onrender.com",
		"https://pos-proxy.d98restaurant.workers.dev",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"X-Forwarded-For",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length", "Authorization"}
	
	return cors.New(config)
}
