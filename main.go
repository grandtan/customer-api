package main

import (
	"customer-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("customers.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	DB.AutoMigrate(&models.Customer{})
}

func createCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Create(&customer)
	c.JSON(http.StatusOK, customer)
}

func updateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer
	if err := DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Save(&customer)
	c.JSON(http.StatusOK, customer)
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer
	if err := DB.Delete(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}

func getCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer
	if err := DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func main() {
	initDatabase()
	r := gin.Default()
	r.POST("/customers", createCustomer)
	r.PUT("/customers/:id", updateCustomer)
	r.DELETE("/customers/:id", deleteCustomer)
	r.GET("/customers/:id", getCustomer)
	r.Run()
}
