package main

import (
	"context"
	"fmt"
	"log"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storage, err := storage.NewStorage(ctx, config.DB)
	if err != nil {
		log.Fatalf("failed to start storage. err = %v", err)
	}

	defer storage.Disconnect(ctx)

	server, err := http.NewServer(config, storage)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Println(fmt.Errorf("failed to run server. err = %v", err))
	}
}
