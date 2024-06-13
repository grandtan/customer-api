package main

import (
	"customer-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("customers.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.Customer{})

	customers := []models.Customer{
		{Name: "John Doe", Age: 30},
		{Name: "Jane Doe", Age: 25},
		{Name: "Mike Smith", Age: 35},
	}

	for _, customer := range customers {
		db.Create(&customer)
	}
}
