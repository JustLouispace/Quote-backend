package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "Qoute-backend/models"
)

type VoteHandler struct {
    db *gorm.DB
}

func NewVoteHandler(db *gorm.DB) *VoteHandler {
    return &VoteHandler{db: db}
}

// CreateVote handles the creation of a new vote
func (h *VoteHandler) CreateVote(c *gin.Context) {
    // Get user ID from context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get quote ID from URL parameter
    quoteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
        return
    }

    // Start a transaction
    tx := h.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Check if user has already voted for any quote
    var existingVote models.Vote
    if err := tx.Where("user_id = ?", userID).First(&existingVote).Error; err == nil {
        tx.Rollback()
        c.JSON(http.StatusConflict, gin.H{"error": "You have already voted for a quote"})
        return
    }

    // Check if quote exists
    var quote models.Quote
    if err := tx.Preload("Votes").First(&quote, quoteID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
        return
    }

    // Create new vote
    vote := models.Vote{
        UserID:  userID.(uint),
        QuoteID: uint(quoteID),
    }

    if err := tx.Create(&vote).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
        return
    }

    // Get updated vote count
    var voteCount int64
    if err := tx.Model(&models.Vote{}).Where("quote_id = ?", quoteID).Count(&voteCount).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vote count"})
        return
    }

    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process vote"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":   "Vote recorded successfully",
        "voteCount": voteCount,
    })
}

// DeleteVote handles the removal of a vote
func (h *VoteHandler) DeleteVote(c *gin.Context) {
    // Get user ID from context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get quote ID from URL parameter
    quoteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
        return
    }

    // Start a transaction
    tx := h.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Delete the vote
    result := tx.Where("user_id = ? AND quote_id = ?", userID, quoteID).Delete(&models.Vote{})
    if result.Error != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vote"})
        return
    }

    if result.RowsAffected == 0 {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"error": "Vote not found"})
        return
    }

    // Get updated vote count
    var voteCount int64
    if err := tx.Model(&models.Vote{}).Where("quote_id = ?", quoteID).Count(&voteCount).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vote count"})
        return
    }

    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process vote deletion"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":   "Vote removed successfully",
        "voteCount": voteCount,
    })
}

// GetVoteCount returns the number of votes for a quote
func (h *VoteHandler) GetVoteCount(c *gin.Context) {
    quoteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
        return
    }

    // Check if quote exists
    var quote models.Quote
    if err := h.db.First(&quote, quoteID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
        return
    }

    var count int64
    if err := h.db.Model(&models.Vote{}).Where("quote_id = ?", quoteID).Count(&count).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get vote count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"count": count})
}

// CheckUserVote checks if the current user has voted for a quote
func (h *VoteHandler) CheckUserVote(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    quoteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
        return
    }

    // Check if quote exists
    var quote models.Quote
    if err := h.db.First(&quote, quoteID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
        return
    }

    var vote models.Vote
    err = h.db.Where("user_id = ? AND quote_id = ?", userID, quoteID).First(&vote).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusOK, gin.H{"has_voted": false})
            return
        }
        // For other database errors, return an error response
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check vote status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"has_voted": true})
}