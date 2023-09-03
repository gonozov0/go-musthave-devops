package server

import (
	"log"
	"net/http"

	"github.com/gonozov0/go-musthave-devops/internal/server/internal/handlers"
	"github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"
)

func Start() {
	repo := storage.NewInMemoryRepository()
	handler := &handlers.UpdateMetricsHandler{Repo: repo}
	http.HandleFunc("/update/", handler.UpdateMetrics)

	log.Printf("Server started at http://127.0.0.1:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
