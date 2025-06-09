package main

import (
	"log"
	"os"

	"Qoute-backend/config"
	"Qoute-backend/handlers"
	"Qoute-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	config.InitDB()

	// Set default port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Auth routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Protected routes
	quotes := router.Group("/quotes")
	quotes.Use(middleware.AuthMiddleware())
	{
		quotes.POST("/", handlers.CreateQuote)
		quotes.GET("/", handlers.GetQuotes)
		quotes.GET("/:id", handlers.GetQuote)
		quotes.PUT("/:id", handlers.UpdateQuote)
		quotes.DELETE("/:id", handlers.DeleteQuote)
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 