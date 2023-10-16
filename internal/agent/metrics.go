package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/avast/retry-go"

	"github.com/gonozov0/go-musthave-devops/internal/shared"
)

var pollCount int64

// CollectMetrics collects metrics from the runtime and returns them as a slice
func CollectMetrics() []shared.Metric {
	var metrics []shared.Metric
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)
	pollCount++

	metrics = append(metrics,
		newGaugeMetric("Alloc", memStats.Alloc),
		newGaugeMetric("BuckHashSys", memStats.BuckHashSys),
		newGaugeMetric("Frees", memStats.Frees),
		newGaugeMetric("GCCPUFraction", memStats.GCCPUFraction),
		newGaugeMetric("GCSys", memStats.GCSys),
		newGaugeMetric("HeapAlloc", memStats.HeapAlloc),
		newGaugeMetric("HeapIdle", memStats.HeapIdle),
		newGaugeMetric("HeapInuse", memStats.HeapInuse),
		newGaugeMetric("HeapObjects", memStats.HeapObjects),
		newGaugeMetric("HeapReleased", memStats.HeapReleased),
		newGaugeMetric("HeapSys", memStats.HeapSys),
		newGaugeMetric("LastGC", memStats.LastGC),
		newGaugeMetric("Lookups", memStats.Lookups),
		newGaugeMetric("MCacheInuse", memStats.MCacheInuse),
		newGaugeMetric("MCacheSys", memStats.MCacheSys),
		newGaugeMetric("MSpanInuse", memStats.MSpanInuse),
		newGaugeMetric("MSpanSys", memStats.MSpanSys),
		newGaugeMetric("Mallocs", memStats.Mallocs),
		newGaugeMetric("NextGC", memStats.NextGC),
		newGaugeMetric("NumForcedGC", memStats.NumForcedGC),
		newGaugeMetric("NumGC", memStats.NumGC),
		newGaugeMetric("OtherSys", memStats.OtherSys),
		newGaugeMetric("PauseTotalNs", memStats.PauseTotalNs),
		newGaugeMetric("StackInuse", memStats.StackInuse),
		newGaugeMetric("StackSys", memStats.StackSys),
		newGaugeMetric("Sys", memStats.Sys),
		newGaugeMetric("TotalAlloc", memStats.TotalAlloc),
		newCounterMetric("PollCount", pollCount),
		newGaugeMetric("RandomValue", time.Now().UnixNano()),
	)

	return metrics
}

func newGaugeMetric(metricName string, metricValue interface{}) shared.Metric {
	var floatValue float64

	switch metricValueTyped := metricValue.(type) {
	case float64:
		floatValue = metricValueTyped
	case int64:
		floatValue = float64(metricValueTyped)
	case uint64:
		floatValue = float64(metricValueTyped)
	case uint32:
		floatValue = float64(metricValueTyped)
	default:
		panic(fmt.Sprintf("invalid metric value type: %T", metricValue))
	}

	return shared.Metric{ID: metricName, MType: shared.Gauge, Value: &floatValue}
}

func newCounterMetric(metricName string, metricValue int64) shared.Metric {
	return shared.Metric{ID: metricName, MType: shared.Counter, Delta: &metricValue}
}

// SendMetrics sends metrics to the server and returns a new metrics slice
func SendMetrics(metrics []shared.Metric, serverAddress string) ([]shared.Metric, error) {
	var url string
	if !strings.HasPrefix(serverAddress, "http") {
		url += "http://"
	}
	url += fmt.Sprintf("%s/updates/", serverAddress)

	client := &http.Client{}
	var (
		statusCode int
		body       []byte
		buffer     bytes.Buffer
	)
	writer := gzip.NewWriter(&buffer)
	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(metrics); err != nil {
		return nil, fmt.Errorf("failed to encode metrics to JSON: %v", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}

	err := retry.Do(
		func() error {
			req, err := http.NewRequest(http.MethodPost, url, &buffer)
			if err != nil {
				return retry.Unrecoverable(err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Accept-Encoding", "gzip")

			r, err := client.Do(req)
			if err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					return err // retry only on timeout network errors
				}
				return retry.Unrecoverable(err)
			}
			defer r.Body.Close()
			statusCode = r.StatusCode
			body, err = io.ReadAll(r.Body)
			if err != nil {
				return retry.Unrecoverable(err)
			}
			return nil
		},
		retry.Attempts(5),
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send metrics: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"received non-OK response while sending metrics: %d, error: %s",
			statusCode,
			string(body),
		)
	}

	return []shared.Metric{}, nil
}
