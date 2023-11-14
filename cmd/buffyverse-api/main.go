package main

import (
	"log"

	"github.com/ElladanTasartir/buffyverse-api/internal/config"
	"github.com/ElladanTasartir/buffyverse-api/internal/http"
)

func main() {
	config, err := config.NewConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	server, err := http.NewServer(config.Port)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
