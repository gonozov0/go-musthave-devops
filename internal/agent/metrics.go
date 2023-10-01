package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gonozov0/go-musthave-devops/internal/shared"

	"github.com/avast/retry-go"
)

// CollectMetrics collects metrics from the runtime and returns them as a slice
func CollectMetrics() []shared.Metric {
	var metrics []shared.Metric
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

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
	)

	return metrics
}

func newGaugeMetric(metricName string, metricValue interface{}) shared.Metric {
	var floatValue float64

	switch metricValueTyped := metricValue.(type) {
	case float64:
		floatValue = metricValueTyped
	case uint64:
		floatValue = float64(metricValueTyped)
	case uint32:
		floatValue = float64(metricValueTyped)
	default:
		panic(fmt.Sprintf("invalid metric value type: %T", metricValue))
	}

	return shared.Metric{ID: metricName, MType: shared.Gauge, Value: &floatValue}
}

// SendMetrics sends metrics to the server and returns a new metrics slice
func SendMetrics(metrics []shared.Metric, serverAddress string) ([]shared.Metric, error) {
	var url string
	if !strings.HasPrefix(serverAddress, "http") {
		url += "http://"
	}
	url += fmt.Sprintf("%s/update/", serverAddress)
	client := &http.Client{}
	var resp *http.Response

	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metric: %v", err)
		}
		err = retry.Do(
			func() error {
				r, err := client.Post(url, "application/json", bytes.NewBuffer(data))
				if err != nil {
					return err
				}
				defer r.Body.Close()
				resp = r
				return nil
			},
			retry.Attempts(2),
			retry.Delay(time.Second),
			retry.MaxJitter(500*time.Millisecond),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to send metrics: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-OK response while sending metrics: %s", resp.Status)
		}
	}
	return []shared.Metric{}, nil
}
