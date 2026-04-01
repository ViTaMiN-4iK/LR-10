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

func TestConcurrentCreateItem(t *testing.T) {
	router := setupRouter()

	concurrency := 10
	done := make(chan bool, concurrency)

	// Запускаем множество параллельных запросов
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			item := models.Item{
				Name:  "Concurrent Item",
				Price: float64(id),
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

			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < concurrency; i++ {
		<-done
	}

	// Проверяем, что все элементы созданы (можно добавить проверку через GET)
	// Это проверяет, что storage работает корректно при конкурентном доступе
}

func TestConcurrentGetItem(t *testing.T) {
	router := setupRouter()

	// Сначала создаем один элемент
	item := models.Item{
		Name:  "Concurrent Get Test",
		Price: 999.99,
	}
	jsonValue, _ := json.Marshal(item)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var createdItem models.Item
	json.Unmarshal(w.Body.Bytes(), &createdItem)
	itemID := createdItem.ID

	// Запускаем множество параллельных GET запросов
	concurrency := 50
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/items/"+itemID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var retrievedItem models.Item
			err := json.Unmarshal(w.Body.Bytes(), &retrievedItem)
			assert.NoError(t, err)
			assert.Equal(t, itemID, retrievedItem.ID)
			assert.Equal(t, "Concurrent Get Test", retrievedItem.Name)

			done <- true
		}()
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}
}

func TestConcurrentCreateAndGet(t *testing.T) {
	router := setupRouter()

	concurrency := 20
	done := make(chan bool, concurrency*2)
	itemIDs := make(chan string, concurrency)

	// Параллельно создаем элементы
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			item := models.Item{
				Name:  "Mixed Test",
				Price: float64(id),
			}
			jsonValue, _ := json.Marshal(item)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code == http.StatusCreated {
				var createdItem models.Item
				json.Unmarshal(w.Body.Bytes(), &createdItem)
				itemIDs <- createdItem.ID
			}
			done <- true
		}(i)
	}

	// Ждем создания всех элементов
	for i := 0; i < concurrency; i++ {
		<-done
	}
	close(itemIDs)

	// Параллельно получаем все созданные элементы
	for id := range itemIDs {
		go func(itemID string) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/items/"+itemID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}(id)
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}
}
