package handlers

import (
	"customer-api/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) UpdateCustomer(c *gin.Context) {
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
