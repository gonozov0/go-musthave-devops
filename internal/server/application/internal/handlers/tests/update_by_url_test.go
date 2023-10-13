package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	repository "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory"
	"github.com/gonozov0/go-musthave-devops/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricByURL(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)

	testCases := []struct {
		name         string
		metricType   string
		metricName   string
		metricValue  interface{}
		expectedCode int
		expectedErr  string
	}{
		{"TestGauge", shared.Gauge, "temperature", 32.5, http.StatusOK, ""},
		{"TestCounter", shared.Counter, "visits", int64(10), http.StatusOK, ""},
		{
			"TestInvalidFloat",
			shared.Gauge,
			"temperature",
			"not-a-float",
			http.StatusBadRequest,
			"Invalid float metricValue\n",
		},
		{
			"TestInvalidFloat",
			shared.Counter,
			"visits",
			"not-an-int",
			http.StatusBadRequest,
			"Invalid integer metricValue\n",
		},
		{"TestUnknownType", "unknownType", "unknown", "0", http.StatusNotImplemented, "Unknown metric type\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("/update/%s/%s/%v", tc.metricType, tc.metricName, tc.metricValue),
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)

			if tc.expectedCode == http.StatusOK {
				switch tc.metricType {
				case shared.Gauge:
					gauge, err := repo.GetGauge(tc.metricName)
					assert.NoError(t, err)
					assert.Equal(t, tc.metricValue, gauge)
				case shared.Counter:
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
