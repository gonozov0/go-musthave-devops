package server

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonozov0/go-musthave-devops/internal/server/internal/handlers/metrics"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"
)

func Run() {
	repo := storage.NewInMemoryRepository()
	handler := metrics.NewHandler(repo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Get("/", handler.GetAllMetrics)
	r.Get("/value/{metricType}/{metricName}", handler.GetMetric)
	r.Post("/update/{metricType}/{metricName}", handler.CreateMetric)

	log.Println("Starting server on port :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}