package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	repository "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory"
	"github.com/gonozov0/go-musthave-devops/internal/shared"
)

func TestBatchUpdateMetrics(t *testing.T) {
	testFloat64 := 32.5
	testInt64 := int64(10)
	testMetrics := []shared.Metric{
		{ID: "temperature", MType: shared.Gauge, Value: &testFloat64},
		{ID: "visits", MType: shared.Counter, Delta: &testInt64},
	}

	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)

	testCases := []struct {
		name            string
		metrics         []shared.Metric
		expectedCode    int
		expectedErr     string
		expectedMetrics []shared.Metric
	}{
		{
			name:         "TestEmptyBody",
			metrics:      nil,
			expectedCode: http.StatusBadRequest,
			expectedErr:  "empty metrics\n",
		},
		{
			name:         "TestEmptyMetrics",
			metrics:      []shared.Metric{},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "empty metrics\n",
		},
		{
			name:            "TestMultipleMetrics",
			metrics:         testMetrics,
			expectedCode:    http.StatusOK,
			expectedMetrics: testMetrics,
		},
		{
			name:         "TestInvalidMetricType",
			metrics:      []shared.Metric{{ID: "unknown", MType: "unknownType"}},
			expectedCode: http.StatusNotImplemented,
			expectedErr:  "Unknown metric type\n",
		},
		{
			name:         "TestInvalidGaugeValue",
			metrics:      []shared.Metric{{ID: "temperature", MType: shared.Gauge, Value: nil}},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "value is required for gauge metric\n",
		},
		{
			name:         "TestInvalidCounterDelta",
			metrics:      []shared.Metric{{ID: "visits", MType: shared.Counter, Delta: nil}},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "delta is required for counter metric\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.metrics)
			require.NoError(t, err)

			request := httptest.NewRequest("POST", "/updates/", bytes.NewReader(body))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			require.Equal(t, tc.expectedCode, recorder.Code)

			if tc.expectedCode == http.StatusOK {
				var metrics []shared.Metric
				err = json.NewDecoder(recorder.Body).Decode(&metrics)
				require.NoError(t, err)
				require.Equal(t, tc.expectedMetrics, metrics)
			} else {
				require.Equal(t, tc.expectedErr, recorder.Body.String())
			}
		})
	}
}
