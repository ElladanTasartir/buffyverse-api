package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/ElladanTasartir/buffyverse-api/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewStorage(ctx context.Context, config config.DBConfig) (*Storage, error) {
	storage := &Storage{}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		return storage, fmt.Errorf("failed to connect to database. err = %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return storage, fmt.Errorf("failed to ping database. err = %v", err)
	}

	storage.client = client
	storage.database = client.Database(config.Database)

	log.Println("successfully connected to DB")

	return storage, nil
}

func (s *Storage) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
