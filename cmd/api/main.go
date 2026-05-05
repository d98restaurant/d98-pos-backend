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
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connections
	mongoClient, err := repository.InitMongoDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())
	log.Println("✅ Connected to MongoDB Atlas")

	// Initialize repositories
	userRepo := repository.NewUserRepository(mongoClient, cfg.MongoDBName)
	orderRepo := repository.NewOrderRepository(mongoClient, cfg.MongoDBName)
	menuRepo := repository.NewMenuRepository(mongoClient, cfg.MongoDBName)
	categoryRepo := repository.NewCategoryRepository(mongoClient, cfg.MongoDBName)
	tableRepo := repository.NewTableRepository(mongoClient, cfg.MongoDBName)
	cartRepo := repository.NewCartRepository(mongoClient, cfg.MongoDBName)
	receiptRepo := repository.NewReceiptRepository(mongoClient, cfg.MongoDBName)
	businessRepo := repository.NewBusinessRepository(mongoClient, cfg.MongoDBName)
	settingsRepo := repository.NewSettingsRepository(mongoClient, cfg.MongoDBName)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	paymentService := services.NewPaymentService(cfg)
	notificationService := services.NewNotificationService(cfg)
	orderService := services.NewOrderService(orderRepo, tableRepo, menuRepo, notificationService)

	// Initialize WebSocket hub
	wsHub := utils.NewHub()
	go wsHub.Run()
	log.Println("✅ WebSocket hub started")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	orderHandler := handlers.NewOrderHandler(orderService, wsHub)
	menuHandler := handlers.NewMenuHandler(menuRepo, categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	tableHandler := handlers.NewTableHandler(tableRepo)
	cartHandler := handlers.NewCartHandler(cartRepo)
	receiptHandler := handlers.NewReceiptHandler(receiptRepo, businessRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService, orderService)
	businessHandler := handlers.NewBusinessHandler(businessRepo)
	settingsHandler := handlers.NewSettingsHandler(settingsRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Setup Gin router
	router := gin.Default()

	// Middleware (with proper CORS for production)
	router.Use(middleware.CORS(cfg.FrontendURL))
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		utils.ServeWebSocket(wsHub, c.Writer, c.Request)
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes (no auth required)
		router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"timestamp": time.Now().Unix(),
				"environment": cfg.Environment,
			})
		})

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/change-password", middleware.AuthMiddleware(cfg), authHandler.ChangePassword)
			auth.GET("/users", middleware.AuthMiddleware(cfg), authHandler.GetUsers)
			auth.PUT("/users/:id", middleware.AuthMiddleware(cfg), authHandler.UpdateUser)
			auth.DELETE("/users/:id", middleware.AuthMiddleware(cfg), authHandler.DeleteUser)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
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
			protected.GET("/orders/table/:tableNumber/active", orderHandler.GetActiveOrdersByTable)
			protected.POST("/orders/table/:tableNumber/complete-billing", orderHandler.CompleteTableBilling)
			protected.GET("/orders/cancellation-requests/pending", orderHandler.GetPendingCancellationRequests)
			protected.POST("/orders/:id/items/:itemId/request-cancellation", orderHandler.RequestItemCancellation)
			protected.POST("/orders/:id/items/:itemId/approve-cancellation", orderHandler.ApproveCancellation)
			protected.POST("/orders/:id/items/:itemId/reject-cancellation", orderHandler.RejectCancellation)
			protected.GET("/orders/credit-customers", orderHandler.GetCreditCustomers)
			protected.POST("/orders/credit-collection", orderHandler.ProcessCreditCollection)
			protected.PATCH("/orders/:id/change-payment", orderHandler.ChangePaymentMethod)

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

			// Receipt routes
			protected.GET("/receipts/:id", receiptHandler.GetReceipt)
			protected.POST("/receipts", receiptHandler.CreateReceipt)

			// Payment routes
			protected.POST("/payments/create-order", paymentHandler.CreateRazorpayOrder)
			protected.POST("/payments/verify", paymentHandler.VerifyPayment)
			protected.POST("/payments/credit-sale", paymentHandler.ProcessCreditSale)

			// Business routes
			protected.GET("/business", businessHandler.GetBusiness)
			protected.POST("/business", businessHandler.UpdateBusiness)

			// Settings routes
			protected.GET("/settings", settingsHandler.GetSettings)
			protected.POST("/settings", settingsHandler.UpdateSettings)

			// Notification routes
			protected.POST("/notifications/subscribe", notificationHandler.Subscribe)
			protected.POST("/notifications/unsubscribe", notificationHandler.Unsubscribe)
			protected.GET("/notifications/vapid-public-key", notificationHandler.GetVAPIDPublicKey)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Port)
		log.Printf("📍 Environment: %s", cfg.Environment)
		log.Printf("🌐 Health check: http://localhost:%s/health", cfg.Port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
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