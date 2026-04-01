package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"go-service/models"
	"go-service/storage"
)

type ItemHandler struct {
	storage *storage.ItemStorage
}

func NewItemHandler(storage *storage.ItemStorage) *ItemHandler {
	return &ItemHandler{
		storage: storage,
	}
}

func (h *ItemHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req models.Item
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Генерируем UUID
	req.ID = uuid.New().String()
	
	h.storage.Create(req)
	
	c.JSON(http.StatusCreated, req)
}

func (h *ItemHandler) GetItem(c *gin.Context) {
	id := c.Param("id")
	
	item, err := h.storage.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	
	c.JSON(http.StatusOK, item)
}