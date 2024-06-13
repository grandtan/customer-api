package handlers

import (
	"customer-api/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) DeleteCustomer(c *gin.Context) {
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
