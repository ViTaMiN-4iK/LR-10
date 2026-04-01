package storage

import (
	"errors"
	"sync"
	
	"go-service/models"
)

type ItemStorage struct {
	mu    sync.RWMutex
	items map[string]models.Item
}

func NewItemStorage() *ItemStorage {
	return &ItemStorage{
		items: make(map[string]models.Item),
	}
}

func (s *ItemStorage) Create(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item.ID] = item
}

func (s *ItemStorage) Get(id string) (models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	item, exists := s.items[id]
	if !exists {
		return models.Item{}, errors.New("item not found")
	}
	return item, nil
}