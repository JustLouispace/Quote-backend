package handlers

import (
	"net/http"

	"Qoute-backend/config"
	"Qoute-backend/models"

	"github.com/gin-gonic/gin"
)

// CreateQuote handles the creation of a new quote
func CreateQuote(c *gin.Context) {
	var quote models.Quote
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Create(&quote)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, quote)
}

// GetQuotes returns all quotes
func GetQuotes(c *gin.Context) {
	var quotes []models.Quote
	result := config.DB.Find(&quotes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, quotes)
}

// GetQuote returns a single quote by ID
func GetQuote(c *gin.Context) {
	id := c.Param("id")
	var quote models.Quote

	result := config.DB.First(&quote, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// UpdateQuote updates an existing quote
func UpdateQuote(c *gin.Context) {
	id := c.Param("id")
	var quote models.Quote

	// First, find the quote
	if err := config.DB.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	// Then update it
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&quote)
	c.JSON(http.StatusOK, quote)
}

// DeleteQuote deletes a quote
func DeleteQuote(c *gin.Context) {
	id := c.Param("id")
	var quote models.Quote

	if err := config.DB.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	config.DB.Delete(&quote)
	c.JSON(http.StatusOK, gin.H{"message": "Quote deleted successfully"})
} 