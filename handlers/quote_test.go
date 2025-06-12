package handlers

import (
	"Qoute-backend/config"
	"Qoute-backend/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupQuoteTestDB() {
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
}

func TestQuoteCRUDAndList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupQuoteTestDB()
	db := config.DB

	// Create user for auth context
	user := models.User{Username: "testuser", Password: "hashed"}
	db.Create(&user)

	r := gin.Default()
	// Mock auth middleware
	r.POST("/quotes", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		CreateQuote(c)
	})
	r.GET("/quotes", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		GetQuotes(c)
	})

	// Create quote
	payload := map[string]string{"content": "inspire", "author": "me"}
	body, _ := json.Marshal(payload)
	req1, _ := http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// List quotes
	req2, _ := http.NewRequest("GET", "/quotes", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "inspire")
	assert.Contains(t, w2.Body.String(), "me")
}
