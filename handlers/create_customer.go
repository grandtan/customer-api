package handlers

import (
	"customer-api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateCustomer(c *gin.Context) {
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
