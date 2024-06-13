package main

import (
	"bytes"
	"customer-api/handlers"
	"customer-api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func SetupRouter(db *gorm.DB) *gin.Engine {
	handlers.InitDatabase(db)
	r := gin.Default()
	r.POST("/customers", handlers.CreateCustomer)
	r.PUT("/customers/:id", handlers.UpdateCustomer)
	r.DELETE("/customers/:id", handlers.DeleteCustomer)
	r.GET("/customers/:id", handlers.GetCustomer)

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

func initTestDatabase() error {
	var err error
	testDB, err = gorm.Open(sqlite.Open("test_customers.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return testDB.AutoMigrate(&models.Customer{})
}

func resetTestDatabase() {
	testDB.Exec("DROP TABLE IF EXISTS customers")
	initTestDatabase()
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	if err := initTestDatabase(); err != nil {
		panic("failed to connect to test database")
	}
	code := m.Run()
	os.Remove("test_customers.db")
	os.Exit(code)
}

func TestCreateCustomer(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)
	customer := models.Customer{Name: "Test User", Age: 20}
	jsonValue, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test with invalid JSON
	req, _ = http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test with missing fields
	req, _ = http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte(`{"Name": ""}`)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateCustomer_ValidData(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)
	customer := models.Customer{Name: "Test User", Age: 20}
	jsonValue, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateCustomer_InvalidJSON(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateCustomer_MissingFields(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)

	// Test with missing Name field
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte(`{"Age": 20}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test with missing Age field
	req, _ = http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte(`{"Name": "Test User"}`)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test with missing both Name and Age fields
	req, _ = http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCustomer(t *testing.T) {
	resetTestDatabase()
	// Insert a customer to get
	customer := models.Customer{Name: "Test User", Age: 20}
	testDB.Create(&customer)
	r := SetupRouter(testDB)
	req, _ := http.NewRequest("GET", "/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test with non-existent customer
	req, _ = http.NewRequest("GET", "/customers/999", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test with invalid ID
	req, _ = http.NewRequest("GET", "/customers/invalid", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateCustomer(t *testing.T) {
	resetTestDatabase()
	// Insert a customer to update
	customer := models.Customer{Name: "Test User", Age: 20}
	testDB.Create(&customer)
	r := SetupRouter(testDB)
	updatedCustomer := models.Customer{Name: "Updated User", Age: 25}
	jsonValue, _ := json.Marshal(updatedCustomer)
	req, _ := http.NewRequest("PUT", "/customers/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test with non-existent customer
	req, _ = http.NewRequest("PUT", "/customers/999", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test with invalid JSON
	req, _ = http.NewRequest("PUT", "/customers/1", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test with invalid ID
	req, _ = http.NewRequest("PUT", "/customers/invalid", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteCustomer(t *testing.T) {
	resetTestDatabase()
	// Insert a customer to delete
	customer := models.Customer{Name: "Test User", Age: 20}
	testDB.Create(&customer)
	r := SetupRouter(testDB)
	req, _ := http.NewRequest("DELETE", "/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test with non-existent customer
	req, _ = http.NewRequest("DELETE", "/customers/999", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test with invalid ID
	req, _ = http.NewRequest("DELETE", "/customers/invalid", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInitDatabase(t *testing.T) {
	resetTestDatabase()
	err := initTestDatabase()
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err)
	}

	// Simulate error scenario by using an invalid driver name
	origDB := testDB
	testDB, err = gorm.Open(sqlite.Open("invalid/db/path"), &gorm.Config{})
	if err == nil {
		t.Fatalf("Expected error when opening invalid database, but got none")
	}
	testDB = origDB
}

func TestCustomerAgeBoundary(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)

	// Test with minimum age
	customer := models.Customer{Name: "Young User", Age: 0}
	testDB.Create(&customer)
	req, _ := http.NewRequest("GET", "/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test with maximum age
	customer = models.Customer{Name: "Old User", Age: 150}
	testDB.Create(&customer)
	req, _ = http.NewRequest("GET", "/customers/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInvalidRoute(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)

	// Test invalid route
	req, _ := http.NewRequest("GET", "/invalidroute", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)

	// Test method not allowed
	req, _ := http.NewRequest("POST", "/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Route not found")
}

func TestEmptyDatabase(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)

	// Test GET on empty database
	req, _ := http.NewRequest("GET", "/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCustomerAgeNegative(t *testing.T) {
	resetTestDatabase()
	r := SetupRouter(testDB)

	// Test with negative age
	customer := models.Customer{Name: "Negative Age User", Age: -1}
	jsonValue, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
