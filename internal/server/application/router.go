package application

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/gonozov0/go-musthave-devops/internal/server/repository"

	"github.com/gonozov0/go-musthave-devops/internal/server/application/internal/handlers"
	"github.com/gonozov0/go-musthave-devops/internal/server/application/internal/middleware"
)

func NewRouter(repo repository.Repository) *chi.Mux {
	handler := handlers.NewHandler(repo)
	router := chi.NewRouter()

	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(chiMiddleware.StripSlashes)
	router.Use(middleware.GzipMiddleware)

	router.Get("/ping", handler.Ping)

	router.Get("/", handler.GetAllMetrics)

	router.Get("/value/{metricType}/{metricName}", handler.GetMetricByURL)
	router.Post("/value", handler.GetMetricByBody)

	router.Post("/update/{metricType}/{metricName}/{metricValue}", handler.UpdateMetricByURL)
	router.Post("/update", handler.UpdateMetricByBody)
	router.Post("/updates", handler.BatchUpdateMetrics)

	return router
}
