package metrics

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetric(t *testing.T) {
	repo := storage.NewInMemoryRepository()
	repo.CreateGauge("temperature", 42.5)
	repo.CreateCounter("visits", 100)
	handler := NewHandler(repo)

	r := chi.NewRouter()
	r.Get("/value/{metricType}/{metricName}", handler.GetMetric)

	testCases := []struct {
		name         string
		metricType   string
		metricName   string
		expectedCode int
		expectedVal  interface{}
	}{
		{"Valid Gauge", Gauge, "temperature", http.StatusOK, float64(42.5)},
		{"Valid Counter", Counter, "visits", http.StatusOK, int64(100)},
		{"Unknown Type", "unknownType", "unknown", http.StatusNotImplemented, nil},
		{"Nonexistent Gauge", Gauge, "nonexistent", http.StatusNotFound, nil},
		{"Nonexistent Counter", Counter, "nonexistent", http.StatusNotFound, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET",
				fmt.Sprintf("/value/%s/%s", tc.metricType, tc.metricName), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, rr.Body.String())

			if tc.expectedCode == http.StatusOK {
				buf := bytes.NewBuffer(rr.Body.Bytes())
				var actualVal interface{}

				switch tc.metricType {
				case Gauge:
					var v float64
					binary.Read(buf, binary.LittleEndian, &v)
					actualVal = v
				case Counter:
					var v int64
					binary.Read(buf, binary.LittleEndian, &v)
					actualVal = v
				}

				assert.Equal(t, tc.expectedVal, actualVal)
			}
		})
	}
}
