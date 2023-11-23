package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ElladanTasartir/buffyverse-api/internal/config"
	"github.com/ElladanTasartir/buffyverse-api/internal/http"
	"github.com/ElladanTasartir/buffyverse-api/internal/storage"
)

func main() {
	config, err := config.NewConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storage, err := storage.NewStorage(ctx, config.DB)
	if err != nil {
		log.Fatalf("failed to start storage. err = %v", err)
	}

	defer storage.Disconnect(ctx)

	server, err := http.NewServer(config, storage)
	if err != nil {
		log.Fatalf("failed to create server. err = %v", err)
	}

	go func() {
		err := server.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err = server.GracefulShutdown(); err != nil {
		log.Fatalf("failed graceful shutdown. err = %v", err)
	}

	os.Exit(0)
}
