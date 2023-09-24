package handlers

import (
	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

// Handler is a struct that holds the repository to update metrics.
type Handler struct {
	Repo repository.Repository
}

// NewHandler constructs a new MetricsHandler.
func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}
