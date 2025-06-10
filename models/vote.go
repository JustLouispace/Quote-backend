package models

import "time"

type Vote struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex"`
    QuoteID   uint      `json:"quote_id" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    User      User      `json:"user" gorm:"foreignKey:UserID"`
    Quote     Quote     `json:"quote" gorm:"foreignKey:QuoteID"`
}