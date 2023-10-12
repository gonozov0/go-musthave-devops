package inmemory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	filestorage "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory/internal/file_storage"
)

func TestLoadMetrics(t *testing.T) {
	fileName := "/tmp/test_load_metrics.json"
	metrics := filestorage.Metrics{
		Gauges: []filestorage.GaugeMetric{
			{
				Name: "TestMetric1",
			},
		},
		Counters: []filestorage.CounterMetric{
			{
				Name: "TestMetric2",
			},
		},
	}
	bytes_, err := json.Marshal(metrics)
	assert.NoError(t, err)

	err = os.WriteFile(fileName, bytes_, 0644)
	assert.NoError(t, err)

	repo, err := NewInMemoryRepositoryWithFileStorage(context.Background(), &sync.WaitGroup{}, fileName, 0, true)
	assert.NoError(t, err)

	gaugeMetric := metrics.Gauges[0]
	gaugeValue, err := repo.GetGauge(gaugeMetric.Name)
	assert.NoError(t, err)
	assert.Equal(t, gaugeMetric.Value, gaugeValue)

	counterMetric := metrics.Counters[0]
	counterValue, err := repo.GetCounter(counterMetric.Name)
	assert.NoError(t, err)
	assert.Equal(t, counterMetric.Value, counterValue)
}

func TestSaveMetrics(t *testing.T) {
	fileName := "/tmp/test_trigger_save.json"
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	repo, err := NewInMemoryRepositoryWithFileStorage(ctx, wg, fileName, 0, false)
	assert.NoError(t, err)

	metricValue := int64(42)
	metricID := "TestMetric"
	_, err = repo.UpdateCounter(metricID, metricValue)
	assert.NoError(t, err)

	time.Sleep(1 * time.Nanosecond)
	cancel()
	wg.Wait()

	bytes_, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(
		t,
		fmt.Sprintf(
			"{\"Gauges\":[],\"Counters\":[{\"name\":\"%s\",\"value\":%d}]}",
			metricID,
			metricValue,
		),
		string(bytes_),
		"File content is not equal to expected",
	)
}
