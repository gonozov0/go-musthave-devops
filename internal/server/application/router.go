package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonozov0/go-musthave-devops/internal/server/application/internal/handlers"
	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

func NewRouter(repo repository.Repository) *chi.Mux {
	handler := handlers.NewHandler(repo)
	router := chi.NewRouter()

	//router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)

	router.Get("/", handler.GetAllMetrics)
	router.Get("/value/{metricType}/{metricName}", handler.GetMetric)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handler.CreateMetric)

	return router
}
