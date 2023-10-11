package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	repository "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory"
	"github.com/gonozov0/go-musthave-devops/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricByBody(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)

	testFloat64 := 32.5
	testInt64 := int64(10)

	testCases := []struct {
		name         string
		metric       shared.Metric
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "TestGauge",
			metric:       shared.Metric{ID: "temperature", MType: shared.Gauge, Value: &testFloat64},
			expectedCode: http.StatusOK,
		},
		{
			name:         "TestCounter",
			metric:       shared.Metric{ID: "visits", MType: shared.Counter, Delta: &testInt64},
			expectedCode: http.StatusOK,
		},
		{
			name:         "TestInvalidValue",
			metric:       shared.Metric{ID: "temperature", MType: shared.Gauge, Value: nil},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid metric value for type Gauge\n",
		},
		{
			name:         "TestInvalidDelta",
			metric:       shared.Metric{ID: "visits", MType: shared.Counter, Delta: nil},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid metric delta for type Counter\n",
		},
		{
			name:         "TestUnknownType",
			metric:       shared.Metric{ID: "unknown", MType: "unknownType"},
			expectedCode: http.StatusNotImplemented,
			expectedErr:  "Unknown metric type\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.metric)
			assert.NoError(t, err, "Failed to marshal metric")

			req, err := http.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(body))
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
			if tc.metric.Value != nil {
				assert.Equal(t, *tc.metric.Value, *resultMetric.Value, "Unexpected metric value")
			}
			if tc.metric.Delta != nil {
				assert.Equal(t, *tc.metric.Delta, *resultMetric.Delta, "Unexpected metric delta")
			}
		})
	}
}
