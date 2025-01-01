package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"bitaksi-go-driver/internal/api"
	"bitaksi-go-driver/internal/config"
	"bitaksi-go-driver/internal/db"
	"bitaksi-go-driver/internal/repository"
	"bitaksi-go-driver/internal/service"
)

// @title Driver Service API
// @version 1.0
// @description This is the API documentation for the Driver Service.
// @host localhost:8080
// @BasePath /driver/api/v1

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize MongoDB
	mongoClient, err := db.ConnectMongo(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Set up database
	db := mongoClient.Database(cfg.MongoDB.Database)

	// Initialize repository and service
	driverRepo := repository.NewDriverRepository(db, "drivers")
	driverService := service.NewDriverService(&driverRepo)

	// Set up router
	router := api.SetupRouter(&driverService, cfg)

	// Add Swagger documentation
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start HTTP server with graceful shutdown
	startServer(router, cfg)
}

// startServer starts the HTTP server and handles graceful shutdown
func startServer(router http.Handler, cfg *config.Config) {
	// Create HTTP server
	address := fmt.Sprintf(":%d", cfg.Ports.HTTP)
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// Run the server in a separate goroutine
	go func() {
		log.Printf("Starting server on %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan

	// Gracefully shut down the server
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shut down: %v", err)
	}

	log.Println("Server exited gracefully")
}
