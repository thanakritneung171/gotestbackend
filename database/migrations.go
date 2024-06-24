package database

import (
	"log"

	"gotestbackend/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Transaction{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	log.Println("Database migration completed.")
}
