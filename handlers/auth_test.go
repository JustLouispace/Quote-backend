package handlers

import (
	"Qoute-backend/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthTestDB() {
	os.Setenv("DATABASE_DSN", ":memory:")
	config.InitDB()
}

func TestRegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAuthTestDB()
	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)

	// Register
	registerPayload := map[string]string{"username": "user1", "password": "pass1"}
	regBody, _ := json.Marshal(registerPayload)
	req1, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(regBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)
	assert.Contains(t, w1.Body.String(), "User registered successfully")

	// Login
	loginPayload := map[string]string{"username": "user1", "password": "pass1"}
	loginBody, _ := json.Marshal(loginPayload)
	req2, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "token")
	assert.Contains(t, w2.Body.String(), "user")
}
