package config

import (
	"log"

	"Qoute-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	
	// Configure GORM with silent logger
	DB, err = gorm.Open(sqlite.Open("quotes.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schema
	err = DB.AutoMigrate(&models.Quote{}, &models.User{}, &models.Vote{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Enable foreign key constraints for SQLite
	db, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	db.SetMaxOpenConns(1) // SQLite only allows one write at a time
}