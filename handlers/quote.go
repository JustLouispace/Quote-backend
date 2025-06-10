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

// QuoteResponse represents a quote with its vote count
type QuoteResponse struct {
	models.Quote
	VoteCount int `json:"voteCount"`
}

// GetQuotes returns all quotes with their vote counts
func GetQuotes(c *gin.Context) {
	var quotes []models.Quote

	// Query params
	author := c.Query("author")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	order := c.DefaultQuery("order", "desc")

	db := config.DB.Preload("Votes") // Preload the Votes relationship

	// Filter by author
	if author != "" {
		db = db.Where("author = ?", author)
	}

	// Search in content or author
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("content LIKE ? OR author LIKE ?", like, like)
	}

	// Sorting
	if sortBy != "" && (order == "asc" || order == "desc") {
		db = db.Order(sortBy + " " + order)
	}

	if err := db.Find(&quotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format with vote counts
	response := make([]QuoteResponse, len(quotes))
	for i, quote := range quotes {
		response[i] = QuoteResponse{
			Quote:     quote,
			VoteCount: len(quote.Votes),
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetQuote returns a single quote by ID with its vote count
func GetQuote(c *gin.Context) {
	id := c.Param("id")
	var quote models.Quote

	if err := config.DB.Preload("Votes").First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	response := QuoteResponse{
		Quote:     quote,
		VoteCount: len(quote.Votes),
	}

	c.JSON(http.StatusOK, response)
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