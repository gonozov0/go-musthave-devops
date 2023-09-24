package agent

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/avast/retry-go"
)

// Metric is a struct that represents a metric
type Metric struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// CollectMetrics collects metrics from the runtime and returns them as a slice
func CollectMetrics() []Metric {
	var metrics []Metric
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	metrics = append(metrics,
		Metric{"Alloc", "gauge", float64(memStats.Alloc)},
		Metric{"BuckHashSys", "gauge", float64(memStats.BuckHashSys)},
		Metric{"Frees", "gauge", float64(memStats.Frees)},
		Metric{"GCCPUFraction", "gauge", memStats.GCCPUFraction},
		Metric{"GCSys", "gauge", float64(memStats.GCSys)},
		Metric{"HeapAlloc", "gauge", float64(memStats.HeapAlloc)},
		Metric{"HeapIdle", "gauge", float64(memStats.HeapIdle)},
		Metric{"HeapInuse", "gauge", float64(memStats.HeapInuse)},
		Metric{"HeapObjects", "gauge", float64(memStats.HeapObjects)},
		Metric{"HeapReleased", "gauge", float64(memStats.HeapReleased)},
		Metric{"HeapSys", "gauge", float64(memStats.HeapSys)},
		Metric{"LastGC", "gauge", float64(memStats.LastGC)},
		Metric{"Lookups", "gauge", float64(memStats.Lookups)},
		Metric{"MCacheInuse", "gauge", float64(memStats.MCacheInuse)},
		Metric{"MCacheSys", "gauge", float64(memStats.MCacheSys)},
		Metric{"MSpanInuse", "gauge", float64(memStats.MSpanInuse)},
		Metric{"MSpanSys", "gauge", float64(memStats.MSpanSys)},
		Metric{"Mallocs", "gauge", float64(memStats.Mallocs)},
		Metric{"NextGC", "gauge", float64(memStats.NextGC)},
		Metric{"NumForcedGC", "gauge", float64(memStats.NumForcedGC)},
		Metric{"NumGC", "gauge", float64(memStats.NumGC)},
		Metric{"OtherSys", "gauge", float64(memStats.OtherSys)},
		Metric{"PauseTotalNs", "gauge", float64(memStats.PauseTotalNs)},
		Metric{"StackInuse", "gauge", float64(memStats.StackInuse)},
		Metric{"StackSys", "gauge", float64(memStats.StackSys)},
		Metric{"Sys", "gauge", float64(memStats.Sys)},
		Metric{"TotalAlloc", "gauge", float64(memStats.TotalAlloc)},
	)

	return metrics
}

// SendMetrics sends metrics to the server and returns a new metrics slice
func SendMetrics(metrics []Metric, serverAddress string) ([]Metric, error) {
	for _, metric := range metrics {
		url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, metric.Type, metric.Name, metric.Value)
		client := &http.Client{}
		var resp *http.Response

		err := retry.Do(
			func() error {
				r, err := client.Post(url, "text/plain", bytes.NewBuffer([]byte{}))
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
	return []Metric{}, nil
}
