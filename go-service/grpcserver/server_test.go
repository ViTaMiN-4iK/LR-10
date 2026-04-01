package grpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-service/grpc/pb"
	"go-service/models"
	"go-service/storage"
)

func setupTestServer() (*ItemServer, *storage.ItemStorage) {
	s := storage.NewItemStorage()
	server := NewItemServer(s)
	return server, s
}

func TestGrpcCreateItem(t *testing.T) {
	server, storage := setupTestServer()

	// Создаем элемент
	req := &pb.CreateItemRequest{
		Name:  "Test Item",
		Price: 99.99,
	}

	resp, err := server.CreateItem(context.Background(), req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Id)
	assert.Equal(t, "Test Item", resp.Name)
	assert.Equal(t, 99.99, resp.Price)

	// Проверяем, что элемент действительно сохранен
	item, err := storage.Get(resp.Id)
	require.NoError(t, err)
	assert.Equal(t, resp.Id, item.ID)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, 99.99, item.Price)
}

func TestGrpcGetItem(t *testing.T) {
	server, storage := setupTestServer()

	// Сначала сохраняем элемент в хранилище
	item := models.Item{
		ID:    "test-id-123",
		Name:  "Existing Item",
		Price: 49.99,
	}
	storage.Create(item)

	// Получаем элемент
	req := &pb.GetItemRequest{Id: "test-id-123"}
	resp, err := server.GetItem(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "test-id-123", resp.Id)
	assert.Equal(t, "Existing Item", resp.Name)
	assert.Equal(t, 49.99, resp.Price)
}

func TestGrpcGetItemNotFound(t *testing.T) {
	server, _ := setupTestServer()

	// Пытаемся получить несуществующий элемент
	req := &pb.GetItemRequest{Id: "non-existent-id"}
	resp, err := server.GetItem(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)

	// Проверяем, что ошибка имеет правильный код
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, statusErr.Code())
	assert.Contains(t, statusErr.Message(), "item not found")
}

func TestGrpcCreateAndGetIntegration(t *testing.T) {
	server, _ := setupTestServer()

	// Создаем элемент через gRPC
	createReq := &pb.CreateItemRequest{
		Name:  "Integration Test",
		Price: 123.45,
	}
	createResp, err := server.CreateItem(context.Background(), createReq)
	require.NoError(t, err)

	// Получаем его через gRPC
	getReq := &pb.GetItemRequest{Id: createResp.Id}
	getResp, err := server.GetItem(context.Background(), getReq)
	require.NoError(t, err)

	assert.Equal(t, createResp.Id, getResp.Id)
	assert.Equal(t, "Integration Test", getResp.Name)
	assert.Equal(t, 123.45, getResp.Price)
}
