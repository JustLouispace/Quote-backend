package handlers

import (
	"Qoute-backend/config"
	"Qoute-backend/middleware"
	"Qoute-backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupInMemoryDB() {
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
}

func TestOneVotePerUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
	db := config.DB
	user := models.User{Username: "user2", Password: "hashed"}
	db.Create(&user)
	quote := models.Quote{Content: "unique", Author: "tester2"}
	db.Create(&quote)

	r := gin.Default()
	voteHandler := NewVoteHandler(db)
	r.POST("/quotes/:id/vote", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		voteHandler.CreateVote(c)
	})

	// First vote
	req1, _ := http.NewRequest("POST", "/quotes/1/vote", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Second vote by same user (should fail with 409)
	req2, _ := http.NewRequest("POST", "/quotes/1/vote", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)
	assert.Contains(t, w2.Body.String(), "already voted")
}

func TestVoteMetadataInResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
	db := config.DB
	user := models.User{Username: "metauser", Password: "hashed"}
	db.Create(&user)
	quote := models.Quote{Content: "meta", Author: "metaauthor"}
	db.Create(&quote)

	r := gin.Default()
	voteHandler := NewVoteHandler(db)
	r.POST("/quotes/:id/vote", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		voteHandler.CreateVote(c)
	})

	req, _ := http.NewRequest("POST", "/quotes/1/vote", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["vote"])
	vote := resp["vote"].(map[string]interface{})
	assert.EqualValues(t, float64(user.ID), vote["user_id"])
	assert.EqualValues(t, float64(quote.ID), vote["quote_id"])
	assert.NotEmpty(t, vote["created_at"])
}

func TestVoteRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
	db := config.DB
	voteHandler := NewVoteHandler(db)

	r := gin.Default()
	r.POST("/quotes/:id/vote", voteHandler.CreateVote)

	req, _ := http.NewRequest("POST", "/quotes/1/vote", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVoteRemoval(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
	db := config.DB
	user := models.User{Username: "deluser", Password: "hashed"}
	db.Create(&user)
	quote := models.Quote{Content: "del", Author: "delauthor"}
	db.Create(&quote)

	r := gin.Default()
	voteHandler := NewVoteHandler(db)
	r.POST("/quotes/:id/vote", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		voteHandler.CreateVote(c)
	})
	r.DELETE("/quotes/:id/vote", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		voteHandler.DeleteVote(c)
	})

	// Vote first
	req1, _ := http.NewRequest("POST", "/quotes/1/vote", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Remove vote
	req2, _ := http.NewRequest("DELETE", "/quotes/1/vote", nil)
	w2 := httptest.NewRecorder()
	assert.NotPanics(t, func() { r.ServeHTTP(w2, req2) })
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Vote removed successfully")
}

func TestVoteOnlyWhenZeroVotes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupInMemoryDB()
	db := config.DB

	// 1. Setup: Create two users, a quote, and a router with real auth middleware
	user1 := models.User{Username: "user1_vote_zero", Password: "password"}
	db.Create(&user1)
	user2 := models.User{Username: "user2_vote_zero", Password: "password"}
	db.Create(&user2)

	quote := models.Quote{Content: "A quote for zero vote test", Author: "Author"}
	db.Create(&quote)

	token1, err := middleware.GenerateJWT(user1.ID)
	assert.NoError(t, err)
	token2, err := middleware.GenerateJWT(user2.ID)
	assert.NoError(t, err)

	router := gin.Default()
	voteHandler := NewVoteHandler(db)
	// Use the actual auth middleware
	authedRoutes := router.Group("/")
	authedRoutes.Use(middleware.AuthMiddleware())
	{
		authedRoutes.POST("/quotes/:id/vote", voteHandler.CreateVote)
	}

	// 2. First vote from User 1 (should succeed)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", fmt.Sprintf("/quotes/%d/vote", quote.ID), nil)
	req1.Header.Set("Authorization", "Bearer "+token1)
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 3. Second vote from User 2 (should fail with conflict because quote already has a vote)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", fmt.Sprintf("/quotes/%d/vote", quote.ID), nil)
	req2.Header.Set("Authorization", "Bearer "+token2)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)
	var resp2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.Contains(t, resp2["error"], "Voting is only allowed when the quote has 0 votes")
}
