package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	repository "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory"
	"github.com/gonozov0/go-musthave-devops/internal/shared"
)

func TestGetMetricByBody(t *testing.T) {
	testFloat64 := 32.5
	testInt64 := int64(10)

	repo := repository.NewInMemoryRepository()
	repo.UpdateGauge("temperature", testFloat64)
	repo.UpdateCounter("visits", testInt64)

	router := application.NewRouter(repo)

	testCases := []struct {
		name         string
		metric       shared.Metric
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "TestGauge",
			metric:       shared.Metric{ID: "temperature", MType: shared.Gauge},
			expectedCode: http.StatusOK,
			expectedErr:  "",
		},
		{
			name:         "TestCounter",
			metric:       shared.Metric{ID: "visits", MType: shared.Counter},
			expectedCode: http.StatusOK,
			expectedErr:  "",
		},
		{
			name:         "TestUnknownType",
			metric:       shared.Metric{ID: "unknown", MType: "unknownType"},
			expectedCode: http.StatusNotImplemented,
			expectedErr:  "Unknown metric type\n",
		},
		{
			name:         "TestNonexistentGauge",
			metric:       shared.Metric{ID: "nonexistent", MType: shared.Gauge},
			expectedCode: http.StatusNotFound,
			expectedErr:  "Metric not found\n",
		},
		{
			name:         "TestNonexistentCounter",
			metric:       shared.Metric{ID: "nonexistent", MType: shared.Counter},
			expectedCode: http.StatusNotFound,
			expectedErr:  "Metric not found\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.metric)
			assert.NoError(t, err, "Failed to marshal metric")

			req, err := http.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(body))
			assert.NoError(t, err, "Failed to create request")
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code, "Unexpected status code")

			if tc.expectedErr != "" {
				respBody := rr.Body.String()
				assert.Contains(t, respBody, tc.expectedErr, "Response body does not contain expected error message")
				return
			}

			resultMetric := shared.Metric{}
			err = json.Unmarshal(rr.Body.Bytes(), &resultMetric)
			assert.NoError(t, err, "Failed to unmarshal metric")
			assert.Equal(t, tc.metric.ID, resultMetric.ID, "Unexpected metric ID")
			assert.Equal(t, tc.metric.MType, resultMetric.MType, "Unexpected metric type")

			switch tc.metric.MType {
			case shared.Gauge:
				assert.Equal(t, testFloat64, *resultMetric.Value, "Unexpected metric value")
			case shared.Counter:
				assert.Equal(t, testInt64, *resultMetric.Delta, "Unexpected metric delta")
			}
		})
	}
}
