package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"bitaksi-go-driver/internal/config"
)

// ConnectMongo connects to MongoDB and returns a client instance.
func ConnectMongo(cfg *config.Config) (*mongo.Client, error) {
	uri := createMongoURI(cfg)

	clientOptions := options.Client().ApplyURI(uri)

	// Use a configurable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the MongoDB server to ensure connectivity
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB successfully!")
	return client, nil
}

// createMongoURI constructs the MongoDB URI from the configuration.
func createMongoURI(cfg *config.Config) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d",
		cfg.MongoDB.Username,
		cfg.MongoDB.Password,
		cfg.MongoDB.Host,
		cfg.MongoDB.Port,
	)
}
