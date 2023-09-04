package metrics

import "github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"

// Handler is a struct that holds the repository to update metrics.
type Handler struct {
	Repo storage.Repository
}

// NewHandler constructs a new MetricsHandler.
func NewHandler(repo storage.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}
