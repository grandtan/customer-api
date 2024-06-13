package handlers

import (
	"customer-api/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(db *gorm.DB) {
	DB = db
	DB.AutoMigrate(&models.Customer{})
}
