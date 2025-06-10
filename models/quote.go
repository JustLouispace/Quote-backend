package models

import (
	"time"

	"gorm.io/gorm"
)

type Quote struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"not null"`
	Author    string         `json:"author" gorm:"not null"`
	Votes     []Vote         `json:"votes,omitempty" gorm:"foreignKey:QuoteID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
