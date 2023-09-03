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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
