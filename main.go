package main

import (
	"log"
	"os"
	"time"

	"Qoute-backend/config"
	"Qoute-backend/handlers"
	"Qoute-backend/middleware"

	"github.com/gin-contrib/cors"
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

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173", "http://127.0.0.1:5173", "https://quote-frontend-zeta.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Auth routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Initialize vote handler
	voteHandler := handlers.NewVoteHandler(config.DB)

	// Protected routes
	quotes := router.Group("/quotes")
	quotes.Use(middleware.AuthMiddleware())
	{
		quotes.POST("/", handlers.CreateQuote)
		quotes.GET("/", handlers.GetQuotes)
		quotes.GET("/:id", handlers.GetQuote)
		quotes.PUT("/:id", handlers.UpdateQuote)
		quotes.DELETE("/:id", handlers.DeleteQuote)

		// Vote routes
		quotes.POST("/:id/vote", voteHandler.CreateVote)
		quotes.DELETE("/:id/vote", voteHandler.DeleteVote)
		quotes.GET("/:id/vote/count", voteHandler.GetVoteCount)
		quotes.GET("/:id/vote/check", voteHandler.CheckUserVote)
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 