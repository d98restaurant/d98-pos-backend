package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pos-backend/config"
	"pos-backend/internal/handlers"
	"pos-backend/internal/middleware"
	"pos-backend/internal/repository"
	"pos-backend/internal/services"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize BadgerDB
	_, err := repository.InitBadgerDB(cfg.BadgerDBPath)
	if err != nil {
		log.Fatalf("Failed to initialize BadgerDB: %v", err)
	}
	defer repository.CloseDB()
	log.Println("✅ BadgerDB connected")

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	orderRepo := repository.NewOrderRepository()
	menuRepo := repository.NewMenuRepository()
	categoryRepo := repository.NewCategoryRepository()
	tableRepo := repository.NewTableRepository()
	cartRepo := repository.NewCartRepository()
	businessRepo := repository.NewBusinessRepository()
	settingsRepo := repository.NewSettingsRepository()

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	paymentService := services.NewPaymentService(cfg)
	notificationService := services.NewNotificationService(cfg)
	orderService := services.NewOrderService(orderRepo, tableRepo, menuRepo, notificationService)

	// WebSocket hub
	wsHub := utils.NewHub()
	go wsHub.Run()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	orderHandler := handlers.NewOrderHandler(orderService, wsHub)
	menuHandler := handlers.NewMenuHandler(menuRepo, categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	tableHandler := handlers.NewTableHandler(tableRepo)
	cartHandler := handlers.NewCartHandler(cartRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService, orderService)
	businessHandler := handlers.NewBusinessHandler(businessRepo)
	settingsHandler := handlers.NewSettingsHandler(settingsRepo)
	resetHandler := handlers.NewResetHandler(userRepo)
	adminHandler := handlers.NewAdminHandler(userRepo)

	// Router
	router := gin.Default()
	router.Use(middleware.CORS(cfg.FrontendURL))
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		utils.ServeWebSocket(wsHub, c.Writer, c.Request)
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"timestamp":   time.Now().Unix(),
			"environment": cfg.Environment,
			"database":    "badgerdb",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/reset-password", resetHandler.ForceResetPassword)
		api.GET("/admin/reset", resetHandler.ClearAndReset)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Auth routes
			protected.POST("/auth/change-password", authHandler.ChangePassword)
			protected.POST("/auth/change-password-admin", adminHandler.ChangeUserPassword)
			protected.GET("/auth/users", authHandler.GetUsers)
			protected.PUT("/auth/users/:id", authHandler.UpdateUser)
			protected.DELETE("/auth/users/:id", authHandler.DeleteUser)

			// Business routes
			protected.GET("/business", businessHandler.GetBusiness)
			protected.POST("/business", businessHandler.UpdateBusiness)

			// Settings routes
			protected.GET("/settings", settingsHandler.GetSettings)
			protected.POST("/settings", settingsHandler.UpdateSettings)

			// Order routes
			protected.GET("/orders", orderHandler.GetOrders)
			protected.GET("/orders/:id", orderHandler.GetOrderByID)
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
			protected.PATCH("/orders/:id/complete-payment", orderHandler.CompletePayment)
			protected.POST("/orders/:id/items", orderHandler.AddItemToOrder)
			protected.PATCH("/orders/:id/items/:itemId", orderHandler.UpdateItemQuantity)
			protected.DELETE("/orders/:id/items/:itemId", orderHandler.RemoveItemFromOrder)
			protected.PATCH("/orders/:id/items/:itemId/status", orderHandler.UpdateItemStatus)
			
			// Menu routes
			protected.GET("/menu", menuHandler.GetMenu)
			protected.POST("/menu", menuHandler.CreateMenuItem)
			protected.PUT("/menu/:id", menuHandler.UpdateMenuItem)
			protected.DELETE("/menu/:id", menuHandler.DeleteMenuItem)
			
			// Category routes
			protected.GET("/categories", categoryHandler.GetCategories)
			protected.POST("/categories", categoryHandler.CreateCategory)
			protected.PUT("/categories/:id", categoryHandler.UpdateCategory)
			protected.DELETE("/categories/:id", categoryHandler.DeleteCategory)
			protected.POST("/categories/reorder", categoryHandler.ReorderCategories)
			
			// Table routes
			protected.GET("/tables", tableHandler.GetTables)
			protected.POST("/tables", tableHandler.CreateTable)
			protected.PATCH("/tables/:tableNumber", tableHandler.UpdateTable)
			protected.DELETE("/tables/:tableNumber", tableHandler.DeleteTable)
			protected.GET("/tables/:tableNumber", tableHandler.GetTableByNumber)
			
			// Cart routes
			protected.GET("/cart", cartHandler.GetCart)
			protected.POST("/cart", cartHandler.SaveCart)
			protected.POST("/cart/items", cartHandler.AddItem)
			protected.PATCH("/cart/items/:itemId", cartHandler.UpdateItemQuantity)
			protected.DELETE("/cart/items/:itemId", cartHandler.RemoveItem)
			protected.DELETE("/cart/clear", cartHandler.ClearCart)
			
			// Payment routes
			protected.POST("/payments/create-order", paymentHandler.CreateRazorpayOrder)
			protected.POST("/payments/verify", paymentHandler.VerifyPayment)
			protected.POST("/payments/credit-sale", paymentHandler.ProcessCreditSale)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Port)
		log.Printf("📍 Environment: %s", cfg.Environment)
		log.Printf("🗄️  Database: BadgerDB at %s", cfg.BadgerDBPath)
		log.Printf("🌐 Health check: http://localhost:%s/health", cfg.Port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
