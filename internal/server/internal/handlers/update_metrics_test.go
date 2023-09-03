package handlers

import (
	"fmt"
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
			method:     "gauge",
			value:      42.0,
			expectCode: http.StatusOK,
		},
		{
			name:       "Test counter",
			method:     "counter",
			value:      int64(42),
			expectCode: http.StatusOK,
		},
	}

	repo := storage.NewInMemoryRepository()
	handler := NewUpdateMetricsHandler(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/TestMetric/%v", tt.method, tt.value)
			req, err := http.NewRequest("POST", url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.UpdateMetrics(rr, req)

			assert.Equal(t, tt.expectCode, rr.Code)

			switch tt.method {
			case "gauge":
				gauges, err := repo.GetGauges()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(gauges))
				assert.Equal(t, tt.value, gauges["TestMetric"])
			case "counter":
				counters, err := repo.GetCounters()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(counters))
				assert.Equal(t, tt.value, counters["TestMetric"])
			}
		})
	}
}
