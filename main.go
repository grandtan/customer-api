package main

import (
	"customer-api/handlers"
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
	handlers.InitDatabase(DB)
}

func setupRouter(h handlers.CustomerHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/customers", h.CreateCustomer)
	r.PUT("/customers/:id", h.UpdateCustomer)
	r.DELETE("/customers/:id", h.DeleteCustomer)
	r.GET("/customers/:id", h.GetCustomer)

	// Handle method not allowed
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	})

	// Handle route not found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
	})

	return r
}

func main() {
	initDatabase()
	h := handlers.NewHandler()
	r := setupRouter(h)
	r.Run()
}
