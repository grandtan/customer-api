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

var testDB *gorm.DB

func SetupRouter(db *gorm.DB) *gin.Engine {
	return setupRouter(db)
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

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}

	// Test with invalid JSON
	req, _ = http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, but got %d", w.Code)
	}
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

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}

	// Test with non-existent customer
	req, _ = http.NewRequest("GET", "/customers/999", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status code 404, but got %d", w.Code)
	}

	// Test with invalid ID
	req, _ = http.NewRequest("GET", "/customers/invalid", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, but got %d", w.Code)
	}
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

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}

	// Test with non-existent customer
	req, _ = http.NewRequest("PUT", "/customers/999", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status code 404, but got %d", w.Code)
	}

	// Test with invalid JSON
	req, _ = http.NewRequest("PUT", "/customers/1", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, but got %d", w.Code)
	}

	// Test with invalid ID
	req, _ = http.NewRequest("PUT", "/customers/invalid", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, but got %d", w.Code)
	}
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

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}

	// Test with non-existent customer
	req, _ = http.NewRequest("DELETE", "/customers/999", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status code 404, but got %d", w.Code)
	}

	// Test with invalid ID
	req, _ = http.NewRequest("DELETE", "/customers/invalid", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code 400, but got %d", w.Code)
	}
}
