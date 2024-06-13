package main

import (
	"bytes"
	"customer-api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/customers", createCustomer)
	r.PUT("/customers/:id", updateCustomer)
	r.DELETE("/customers/:id", deleteCustomer)
	r.GET("/customers/:id", getCustomer)
	return r
}

func initTestDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test_customers.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	DB.AutoMigrate(&models.Customer{})
}

func resetTestDatabase() {
	DB.Exec("DROP TABLE IF EXISTS customers")
	initTestDatabase()
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	initTestDatabase()
	code := m.Run()
	os.Remove("test_customers.db")
	os.Exit(code)
}

func TestCreateCustomer(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter()
	customer := models.Customer{Name: "Test User", Age: 20}
	jsonValue, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
}

func TestGetCustomer(t *testing.T) {
	resetTestDatabase()
	// Insert a customer to get
	customer := models.Customer{Name: "Test User", Age: 20}
	DB.Create(&customer)
	r := SetupRouter()
	req, _ := http.NewRequest("GET", "/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
}

func TestUpdateCustomer(t *testing.T) {
	resetTestDatabase()
	// Insert a customer to update
	customer := models.Customer{Name: "Test User", Age: 20}
	DB.Create(&customer)
	r := SetupRouter()
	updatedCustomer := models.Customer{Name: "Updated User", Age: 25}
	jsonValue, _ := json.Marshal(updatedCustomer)
	req, _ := http.NewRequest("PUT", "/customers/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
}

func TestDeleteCustomer(t *testing.T) {
	resetTestDatabase()
	// Insert a customer to delete
	customer := models.Customer{Name: "Test User", Age: 20}
	DB.Create(&customer)
	r := SetupRouter()
	req, _ := http.NewRequest("DELETE", "/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
}
