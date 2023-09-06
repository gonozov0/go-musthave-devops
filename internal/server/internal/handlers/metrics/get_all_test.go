package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"
)

func TestGetAllMetrics(t *testing.T) {
	t.Run("empty repository", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		handler := NewHandler(repo)

		req := httptest.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()

		handler.GetAllMetrics(w, req)

		res := w.Result()
		body, _ := io.ReadAll(res.Body)

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", res.StatusCode)
		}

		if !strings.Contains(string(body), "<h1>Metrics</h1>") {
			t.Errorf("HTML body does not contain header <h1>Metrics</h1>")
		}
	})

	t.Run("repository with metrics", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		repo.CreateGauge("gauge1", 10.5)
		repo.CreateCounter("counter1", 2)

		handler := NewHandler(repo)

		req := httptest.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()

		handler.GetAllMetrics(w, req)

		res := w.Result()
		body, _ := io.ReadAll(res.Body)

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", res.StatusCode)
		}

		if !strings.Contains(string(body), "<li>gauge1: 10.5</li>") {
			t.Errorf("HTML body does not contain gauge1 metrics")
		}

		if !strings.Contains(string(body), "<li>counter1: 2</li>") {
			t.Errorf("HTML body does not contain counter1 metrics")
		}
	})
}
