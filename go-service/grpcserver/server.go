package grpcserver

import (
	"context"

	"go-service/grpc/pb"
	"go-service/models"
	"go-service/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ItemServer struct {
	pb.UnimplementedItemServiceServer
	storage *storage.ItemStorage
}

func NewItemServer(storage *storage.ItemStorage) *ItemServer {
	return &ItemServer{
		storage: storage,
	}
}

func (s *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.ItemResponse, error) {
	item, err := s.storage.Get(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "item not found")
	}

	return &pb.ItemResponse{
		Id:    item.ID,
		Name:  item.Name,
		Price: item.Price,
	}, nil
}

func (s *ItemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.ItemResponse, error) {
	// Генерируем UUID
	itemID := uuid.New().String()

	item := models.Item{
		ID:    itemID,
		Name:  req.Name,
		Price: req.Price,
	}

	s.storage.Create(item)

	return &pb.ItemResponse{
		Id:    item.ID,
		Name:  item.Name,
		Price: item.Price,
	}, nil
}
