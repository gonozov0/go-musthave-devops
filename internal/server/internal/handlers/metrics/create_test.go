package metrics

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"
)

func TestUpdateMetricsHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		value      interface{}
		expectCode int
	}{
		{
			name:       "Test gauge",
			method:     Gauge,
			value:      42.0,
			expectCode: http.StatusOK,
		},
		{
			name:       "Test counter",
			method:     Counter,
			value:      int64(42),
			expectCode: http.StatusOK,
		},
	}

	repo := storage.NewInMemoryRepository()
	handler := NewHandler(repo)

	r := chi.NewRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricName := "TestMetric"
			url := fmt.Sprintf("/update/%s/%s/%v", tt.method, metricName, tt.value)
			req, err := http.NewRequest("POST", url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			r.Post("/update/{metricType}/{metricName}/{value}", handler.CreateMetric)
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectCode, rr.Code)

			switch tt.method {
			case Gauge:
				gauge, err := repo.GetGauge(metricName)
				assert.NoError(t, err)
				assert.Equal(t, tt.value, gauge)
			case Counter:
				counter, err := repo.GetCounter(metricName)
				assert.NoError(t, err)
				assert.Equal(t, tt.value, counter)
			}
		})
	}
}
