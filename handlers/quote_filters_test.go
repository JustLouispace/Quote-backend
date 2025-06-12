package handlers

import (
	"Qoute-backend/config"
	"Qoute-backend/middleware"
	"Qoute-backend/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouterWithQuotes initializes the router and creates a set of quotes for testing.
func setupTestRouterWithQuotes(t *testing.T) (*gin.Engine, []models.Quote) {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
	db := config.DB

	// Clean up previous data
	db.Migrator().DropTable(&models.User{}, &models.Quote{}, &models.Vote{})
	db.AutoMigrate(&models.User{}, &models.Quote{}, &models.Vote{})

	quotes := []models.Quote{
		{Content: "The only true wisdom is in knowing you know nothing.", Author: "Socrates"},
		{Content: "An unexamined life is not worth living.", Author: "Socrates"},
		{Content: "Be the change that you wish to see in the world.", Author: "Mahatma Gandhi"},
		{Content: "I think, therefore I am.", Author: "René Descartes"},
	}

	for i := range quotes {
		db.Create(&quotes[i])
	}

	router := gin.Default()

	// Register routes directly since they are package-level functions
	router.POST("/register", Register)
	router.POST("/login", Login)

	authed := router.Group("/")
	authed.Use(middleware.AuthMiddleware())
	{
		authed.GET("/quotes", GetQuotes)
	}

	return router, quotes
}

func TestGetQuotesFiltering(t *testing.T) {
	router, _ := setupTestRouterWithQuotes(t)

	// Create a user and get a token to access the endpoint
	user := models.User{Username: "filter_user", Password: "password"}
	config.DB.Create(&user)
	token, _ := middleware.GenerateJWT(user.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/quotes?author=Socrates", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []QuoteResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Len(t, response, 2)
	for _, quote := range response {
		assert.Equal(t, "Socrates", quote.Author)
	}
}

func TestGetQuotesSearch(t *testing.T) {
	router, _ := setupTestRouterWithQuotes(t)

	// Create a user and get a token
	user := models.User{Username: "search_user", Password: "password"}
	config.DB.Create(&user)
	token, _ := middleware.GenerateJWT(user.ID)

	w := httptest.NewRecorder()
	// Search for "the" which is in two quotes
	req, _ := http.NewRequest("GET", "/quotes?search=the", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []QuoteResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Len(t, response, 2)
	// Check if the response contains the expected quotes
	foundSocrates := false
	foundGandhi := false
	for _, quote := range response {
		if quote.Author == "Socrates" {
			foundSocrates = true
		}
		if quote.Author == "Mahatma Gandhi" {
			foundGandhi = true
		}
	}
	assert.True(t, foundSocrates, "Expected to find quote by Socrates")
	assert.True(t, foundGandhi, "Expected to find quote by Mahatma Gandhi")
}

func TestGetQuotesSorting(t *testing.T) {
	router, _ := setupTestRouterWithQuotes(t)

	// Create a user and get a token
	user := models.User{Username: "sort_user", Password: "password"}
	config.DB.Create(&user)
	token, _ := middleware.GenerateJWT(user.ID)

	// Test sorting by author ascending
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/quotes?sortBy=author&order=asc", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	var resp1 []QuoteResponse
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.Len(t, resp1, 4)
	assert.Equal(t, "Mahatma Gandhi", resp1[0].Author)
	assert.Equal(t, "René Descartes", resp1[1].Author)
	assert.Equal(t, "Socrates", resp1[2].Author)
	assert.Equal(t, "Socrates", resp1[3].Author)

	// Test sorting by content descending
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/quotes?sortBy=content&order=desc", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp2 []QuoteResponse
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.Len(t, resp2, 4)
	assert.Equal(t, "The only true wisdom is in knowing you know nothing.", resp2[0].Content)
}
