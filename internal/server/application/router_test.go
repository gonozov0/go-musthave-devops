package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/shared"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricByURL(t *testing.T) {
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
		{"TestGauge", shared.Gauge, "temperature", 32.5, http.StatusOK, ""},
		{"TestCounter", shared.Counter, "visits", int64(10), http.StatusOK, ""},
		{"TestInvalidFloat", shared.Gauge, "temperature", "not-a-float", http.StatusBadRequest, "Invalid float metricValue\n"},
		{"TestInvalidFloat", shared.Counter, "visits", "not-an-int", http.StatusBadRequest, "Invalid integer metricValue\n"},
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
		repo.UpdateGauge("gauge1", 10.5)
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

func TestGetMetricByURL(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	repo.UpdateGauge("temperature", 42.5)
	repo.UpdateCounter("visits", 100)

	router := NewRouter(repo)

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

func TestUpdateMetricByBody(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := NewRouter(repo)

	testFloat64 := 32.5
	testInt64 := int64(10)

	testCases := []struct {
		name         string
		metric       shared.Metric
		expectedCode int
		expectedErr  string
	}{
		{name: "TestGauge", metric: shared.Metric{ID: "temperature", MType: shared.Gauge, Value: &testFloat64}, expectedCode: http.StatusOK},
		{name: "TestCounter", metric: shared.Metric{ID: "visits", MType: shared.Counter, Delta: &testInt64}, expectedCode: http.StatusOK},
		{name: "TestInvalidValue", metric: shared.Metric{ID: "temperature", MType: shared.Gauge, Value: nil}, expectedCode: http.StatusBadRequest, expectedErr: "Invalid metric value for type Gauge\n"},
		{name: "TestInvalidDelta", metric: shared.Metric{ID: "visits", MType: shared.Counter, Delta: nil}, expectedCode: http.StatusBadRequest, expectedErr: "Invalid metric delta for type Counter\n"},
		{name: "TestUnknownType", metric: shared.Metric{ID: "unknown", MType: "unknownType"}, expectedCode: http.StatusNotImplemented, expectedErr: "Unknown metric type\n"},
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

func TestGetMetricByBody(t *testing.T) {
	testFloat64 := 32.5
	testInt64 := int64(10)

	repo := repository.NewInMemoryRepository()
	repo.UpdateGauge("temperature", testFloat64)
	repo.UpdateCounter("visits", testInt64)

	router := NewRouter(repo)

	testCases := []struct {
		name         string
		metric       shared.Metric
		expectedCode int
		expectedErr  string
	}{
		{name: "TestGauge", metric: shared.Metric{ID: "temperature", MType: shared.Gauge}, expectedCode: http.StatusOK, expectedErr: ""},
		{name: "TestCounter", metric: shared.Metric{ID: "visits", MType: shared.Counter}, expectedCode: http.StatusOK, expectedErr: ""},
		{name: "TestUnknownType", metric: shared.Metric{ID: "unknown", MType: "unknownType"}, expectedCode: http.StatusNotImplemented, expectedErr: "Unknown metric type\n"},
		{name: "TestNonexistentGauge", metric: shared.Metric{ID: "nonexistent", MType: shared.Gauge}, expectedCode: http.StatusNotFound, expectedErr: "Metric not found\n"},
		{name: "TestNonexistentCounter", metric: shared.Metric{ID: "nonexistent", MType: shared.Counter}, expectedCode: http.StatusNotFound, expectedErr: "Metric not found\n"},
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
