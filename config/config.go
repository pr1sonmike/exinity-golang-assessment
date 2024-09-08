package config

import (
	"exinity-golang-assessment/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("POSTGRES_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&models.Transaction{})
	return db
}
