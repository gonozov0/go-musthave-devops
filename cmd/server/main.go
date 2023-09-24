package main

import (
	"log"
	"net/http"

	"github.com/gonozov0/go-musthave-devops/internal/server"
	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %s", err.Error())
	}

	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)
	log.Printf("Starting server on port %s\n", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
