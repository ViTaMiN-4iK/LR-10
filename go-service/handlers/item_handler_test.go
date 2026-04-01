package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	
	"go-service/models"
	"go-service/storage"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	storage := storage.NewItemStorage()
	handler := NewItemHandler(storage)
	
	router := gin.Default()
	router.GET("/health", handler.Health)
	router.POST("/items", handler.CreateItem)
	router.GET("/items/:id", handler.GetItem)
	
	return router
}

func TestHealth(t *testing.T) {
	router := setupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestCreateItem(t *testing.T) {
	router := setupRouter()
	
	item := models.Item{
		Name:  "Test Item",
		Price: 99.99,
	}
	
	jsonValue, _ := json.Marshal(item)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var createdItem models.Item
	err := json.Unmarshal(w.Body.Bytes(), &createdItem)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdItem.ID)
	assert.Equal(t, "Test Item", createdItem.Name)
	assert.Equal(t, 99.99, createdItem.Price)
}

func TestGetItemNotFound(t *testing.T) {
	router := setupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items/non-existent-id", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetItemSuccess(t *testing.T) {
	router := setupRouter()
	
	// Сначала создаем элемент
	item := models.Item{
		Name:  "Get Test Item",
		Price: 49.99,
	}
	
	jsonValue, _ := json.Marshal(item)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createdItem models.Item
	json.Unmarshal(w.Body.Bytes(), &createdItem)
	
	// Теперь получаем его
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/items/"+createdItem.ID, nil)
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var retrievedItem models.Item
	err := json.Unmarshal(w2.Body.Bytes(), &retrievedItem)
	assert.NoError(t, err)
	assert.Equal(t, createdItem.ID, retrievedItem.ID)
	assert.Equal(t, "Get Test Item", retrievedItem.Name)
	assert.Equal(t, 49.99, retrievedItem.Price)
}