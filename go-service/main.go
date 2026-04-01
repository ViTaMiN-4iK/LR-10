package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/gin-gonic/gin"
	
	"go-service/handlers"
	"go-service/storage"
)

func main() {
	// Инициализация хранилища
	itemStorage := storage.NewItemStorage()
	itemHandler := handlers.NewItemHandler(itemStorage)
	
	// Настройка роутера
	router := gin.Default()
	
	// Регистрация эндпоинтов
	router.GET("/health", itemHandler.Health)
	router.POST("/items", itemHandler.CreateItem)
	router.GET("/items/:id", itemHandler.GetItem)
	
	// Настройка HTTP сервера
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	
	// Запуск сервера в горутине
	go func() {
		log.Println("Starting Go Gin server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	
	// Таймаут для завершения текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	
	log.Println("Server exiting")
}