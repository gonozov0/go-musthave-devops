package application

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server/application/internal/handlers"
	"github.com/gonozov0/go-musthave-devops/internal/server/repository"

	"github.com/stretchr/testify/assert"
)

func TestCreateMetric(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := NewRouter(repo)

	testCases := []struct {
		name         string
		metricType   string
		metricName   string
		metricValue  interface{}
		expectedCode int
		expectedErr  string
	}{
		{"TestGauge", handlers.Gauge, "temperature", 32.5, http.StatusOK, ""},
		{"TestCounter", handlers.Counter, "visits", int64(10), http.StatusOK, ""},
		{"TestInvalidFloat", handlers.Gauge, "temperature", "not-a-float", http.StatusBadRequest, "Invalid float metricValue\n"},
		{"TestInvalidFloat", handlers.Counter, "visits", "not-an-int", http.StatusBadRequest, "Invalid integer metricValue\n"},
		{"TestUnknownType", "unknownType", "unknown", "0", http.StatusNotImplemented, "Unknown metric type\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("/update/%s/%s/%v", tc.metricType, tc.metricName, tc.metricValue), nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)

			if tc.expectedCode == http.StatusOK {
				switch tc.metricType {
				case handlers.Gauge:
					gauge, err := repo.GetGauge(tc.metricName)
					assert.NoError(t, err)
					assert.Equal(t, tc.metricValue, gauge)
				case handlers.Counter:
					counter, err := repo.GetCounter(tc.metricName)
					assert.NoError(t, err)
					assert.Equal(t, tc.metricValue, counter)
				}
			} else {
				assert.Equal(t, tc.expectedErr, recorder.Body.String())
			}
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := NewRouter(repo)
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
		repo.CreateGauge("gauge1", 10.5)
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

func TestGetMetric(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	repo.CreateGauge("temperature", 42.5)
	repo.UpdateCounter("visits", 100)

	router := NewRouter(repo)

	testCases := []struct {
		name         string
		metricType   string
		metricName   string
		expectedCode int
		expectedVal  string
	}{
		{"Valid Gauge", handlers.Gauge, "temperature", http.StatusOK, "42.5"},
		{"Valid Counter", handlers.Counter, "visits", http.StatusOK, "100"},
		{"Unknown Type", "unknownType", "unknown", http.StatusNotImplemented, ""},
		{"Nonexistent Gauge", handlers.Gauge, "nonexistent", http.StatusNotFound, ""},
		{"Nonexistent Counter", handlers.Counter, "nonexistent", http.StatusNotFound, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET",
				fmt.Sprintf("/value/%s/%s", tc.metricType, tc.metricName), nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			strBody := recorder.Body.String()

			assert.Equal(t, tc.expectedCode, recorder.Code, strBody)

			if tc.expectedCode == http.StatusOK {

				assert.Equal(t, tc.expectedVal, strBody)
			}
		})
	}
}
