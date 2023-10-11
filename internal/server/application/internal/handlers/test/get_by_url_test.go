package handlers_test

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

func TestGetMetricByURL(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	repo.UpdateGauge("temperature", 42.5)
	repo.UpdateCounter("visits", 100)

	router := application.NewRouter(repo)

	testCases := []struct {
		name         string
		metricType   string
		metricName   string
		expectedCode int
		expectedVal  string
	}{
		{"Valid Gauge", shared.Gauge, "temperature", http.StatusOK, "42.5"},
		{"Valid Counter", shared.Counter, "visits", http.StatusOK, "100"},
		{"Unknown Type", "unknownType", "unknown", http.StatusNotImplemented, ""},
		{"Nonexistent Gauge", shared.Gauge, "nonexistent", http.StatusNotFound, ""},
		{"Nonexistent Counter", shared.Counter, "nonexistent", http.StatusNotFound, ""},
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
