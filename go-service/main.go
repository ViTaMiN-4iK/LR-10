package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"go-service/grpc/pb"
	"go-service/grpcserver"
	"go-service/handlers"
	"go-service/storage"
)

func main() {
	// Инициализация хранилища
	itemStorage := storage.NewItemStorage()

	// ========== REST (Gin) Сервер ==========
	itemHandler := handlers.NewItemHandler(itemStorage)

	// Настройка роутера
	router := gin.Default()

	// Регистрация эндпоинтов
	router.GET("/health", itemHandler.Health)
	router.POST("/items", itemHandler.CreateItem)
	router.GET("/items/:id", itemHandler.GetItem)

	// Настройка HTTP сервера
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Запуск HTTP сервера в горутине
	go func() {
		log.Println("Starting Go Gin REST server on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server error: %s\n", err)
		}
	}()

	// ========== gRPC Сервер ==========
	// Создаем gRPC сервер
	grpcSrv := grpc.NewServer()
	itemServer := grpcserver.NewItemServer(itemStorage)
	pb.RegisterItemServiceServer(grpcSrv, itemServer)

	// Слушаем порт 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	// Запуск gRPC сервера в горутине
	go func() {
		log.Println("Starting Go gRPC server on :50051")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// ========== Graceful Shutdown ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Shutdown HTTP сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("HTTP server forced to shutdown:", err)
	}

	// Shutdown gRPC сервера
	grpcSrv.GracefulStop()

	log.Println("All servers exited")
}
