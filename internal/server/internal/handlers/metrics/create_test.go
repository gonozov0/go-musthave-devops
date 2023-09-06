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

func TestCreateMetric(t *testing.T) {
	repo := storage.NewInMemoryRepository()
	handler := NewHandler(repo)

	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handler.CreateMetric)

	testCases := []struct {
		name         string
		metricType   string
		metricName   string
		metricValue  interface{}
		expectedCode int
		expectedErr  string
	}{
		{"TestGauge", Gauge, "temperature", 32.5, http.StatusOK, ""},
		{"TestCounter", Counter, "visits", int64(10), http.StatusOK, ""},
		{"TestInvalidFloat", Gauge, "temperature", "not-a-float", http.StatusBadRequest, "Invalid float metricValue\n"},
		{"TestInvalidFloat", Counter, "visits", "not-an-int", http.StatusBadRequest, "Invalid integer metricValue\n"},
		{"TestUnknownType", "unknownType", "unknown", "0", http.StatusNotImplemented, "Unknown metric type\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("/update/%s/%s/%v", tc.metricType, tc.metricName, tc.metricValue), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)

			if tc.expectedCode == http.StatusOK {
				switch tc.metricType {
				case Gauge:
					gauge, err := repo.GetGauge(tc.metricName)
					assert.NoError(t, err)
					assert.Equal(t, tc.metricValue, gauge)
				case Counter:
					counter, err := repo.GetCounter(tc.metricName)
					assert.NoError(t, err)
					assert.Equal(t, tc.metricValue, counter)
				}
			} else {
				assert.Equal(t, tc.expectedErr, rr.Body.String())
			}
		})
	}
}
