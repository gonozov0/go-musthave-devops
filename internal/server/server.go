package server

import (
	"log"
	"net/http"

	"github.com/gonozov0/go-musthave-devops/internal/server/internal/application"

	"github.com/gonozov0/go-musthave-devops/internal/server/internal/repository"
)

func Run() {
	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)
	log.Println("Starting server on port :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
