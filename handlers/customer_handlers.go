package handlers

import (
	"customer-api/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(db *gorm.DB) {
	DB = db
	DB.AutoMigrate(&models.Customer{})
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if customer.Name == "" || customer.Age <= 0 {
		log.Println("Invalid customer data:", customer)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and a positive Age are required"})
		return
	}
	if err := DB.Create(&customer).Error; err != nil {
		log.Println("Error creating customer:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if _, err := strconv.Atoi(id); err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := DB.First(&customer, id).Error; err != nil {
		log.Println("Customer not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	if err := c.ShouldBindJSON(&customer); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if customer.Name == "" || customer.Age < 0 {
		log.Println("Invalid customer data:", customer)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer data"})
		return
	}
	if err := DB.Save(&customer).Error; err != nil {
		log.Println("Error updating customer:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if _, err := strconv.Atoi(id); err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := DB.First(&customer, id).Error; err != nil {
		log.Println("Customer not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	if err := DB.Delete(&customer).Error; err != nil {
		log.Println("Error deleting customer:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}

func GetCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if _, err := strconv.Atoi(id); err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := DB.First(&customer, id).Error; err != nil {
		log.Println("Customer not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}
