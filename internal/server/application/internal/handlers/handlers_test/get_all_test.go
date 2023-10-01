package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

func TestGetAllMetrics(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)
	request := httptest.NewRequest("GET", "/", nil)

	t.Run("empty repository", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		strBody := recorder.Body.String()

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", recorder.Code)
		}

		if !strings.Contains(strBody, "<h1>Metrics</h1>") {
			t.Errorf("HTML body does not contain header <h1>Metrics</h1>")
		}
	})

	t.Run("repository with metrics", func(t *testing.T) {
		repo.UpdateGauge("gauge1", 10.5)
		repo.UpdateCounter("counter1", 2)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		strBody := recorder.Body.String()

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", recorder.Code)
		}

		if !strings.Contains(strBody, "<li>gauge1: 10.5</li>") {
			t.Errorf("HTML body does not contain gauge1 metrics")
		}

		if !strings.Contains(strBody, "<li>counter1: 2</li>") {
			t.Errorf("HTML body does not contain counter1 metrics")
		}
	})
}
