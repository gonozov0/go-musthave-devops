package middleware_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	repository "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory"
	"github.com/gonozov0/go-musthave-devops/internal/shared"
	"github.com/stretchr/testify/require"
)

func TestGZipMiddleware(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)

	testFloat := 32.5
	testMetric := shared.Metric{ID: "temperature", MType: shared.Gauge, Value: &testFloat}

	body, err := json.Marshal(testMetric)
	require.NoError(t, err, "Failed to marshal metric")

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, _ = gz.Write(body)
	_ = gz.Close()

	req, err := http.NewRequest(http.MethodPost, "/update/", &b)
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	gr, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)
	respBody, err := io.ReadAll(gr)
	require.NoError(t, err)

	resultMetric := shared.Metric{}
	err = json.Unmarshal(respBody, &resultMetric)
	require.NoError(t, err, "Failed to unmarshal metric")
	require.Equal(t, testMetric.ID, resultMetric.ID, "Unexpected metric ID")
	require.Equal(t, testMetric.MType, resultMetric.MType, "Unexpected metric type")
	require.Equal(t, *testMetric.Value, *resultMetric.Value, "Unexpected metric value")
}
