package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
	filestorage "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory/internal/file_storage"
)

// inMemoryRepository is an in-memory implementation of the Repository interface.
type inMemoryRepository struct {
	gaugeMu     sync.RWMutex
	counterMu   sync.RWMutex
	gauges      map[string]float64
	counters    map[string]int64
	fileStorage *filestorage.FileStorage
	saveTicker  *time.Ticker
}

// NewInMemoryRepository creates a new inMemoryRepository and returns it as a Repository interface.
func NewInMemoryRepository() repository.Repository {
	return &inMemoryRepository{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// NewInMemoryRepositoryWithFileStorage creates a new in-memory repository with file storage capability.
//
// Parameters:
// ctx: The context for canceling background operations.
// wg: A wait group to track the completion of background operations.
// filePath: The path to the file where the metrics will be stored.
// intervalInSeconds: The interval in seconds to periodically save metrics to file storage.
//
//	If it is set to 0, the default minimal interval will be used (practically immediate).
//
// restore: If true, the repository will attempt to restore metrics from the specified file storage at initialization.
//
// Returns:
// - An instance of repository.Repository populated with metrics from the file storage (if restore is true and no errors occurred).
// - An error if there's a failure in loading metrics from file storage.
//
// Note:
// If the restore operation is enabled and fails, the function will return an error, and no repository will be created.
func NewInMemoryRepositoryWithFileStorage(
	ctx context.Context,
	wg *sync.WaitGroup,
	filePath string,
	intervalInSeconds uint64,
	restore bool,
) (repository.Repository, error) {
	var (
		metrics filestorage.Metrics
		err     error
	)
	fileStorage := filestorage.NewFileStorage(filePath)
	if restore {
		metrics, err = fileStorage.LoadMetrics()
		if err != nil {
			return nil, fmt.Errorf("failed to load metrics from file storage: %w", err)
		}
	}

	var duration time.Duration
	if intervalInSeconds > 0 {
		duration = time.Duration(intervalInSeconds) * time.Second
	} else {
		duration = time.Duration(1) * time.Nanosecond
	}
	saveTicker := time.NewTicker(duration)

	repo := inMemoryRepository{
		gauges:      gaugesToMap(metrics.Gauges),
		counters:    countersToMap(metrics.Counters),
		fileStorage: fileStorage,
		saveTicker:  saveTicker,
	}
	repo.startSaveMetricsInBackground(ctx, wg)

	return &repo, nil
}

// Ping checks the connection to the repository.
func (repo *inMemoryRepository) Ping() error {
	return nil
}

// UpdateGauge updates or sets a new gauge metric with the given name and value.
func (repo *inMemoryRepository) UpdateGauge(metricName string, value float64) (float64, error) {
	repo.gaugeMu.Lock()
	repo.gauges[metricName] = value
	repo.gaugeMu.Unlock()
	return value, nil
}

// UpdateCounter updates or sets a new counter metric with the given name and value.
func (repo *inMemoryRepository) UpdateCounter(metricName string, value int64) (int64, error) {
	repo.counterMu.Lock()
	newValue := repo.counters[metricName] + value
	repo.counters[metricName] = newValue
	repo.counterMu.Unlock()
	return newValue, nil
}

// UpdateGauges updates or sets a new gauge metrics with the given name and value.
func (repo *inMemoryRepository) UpdateGauges(metrics []repository.GaugeMetric) ([]repository.GaugeMetric, error) {
	repo.gaugeMu.Lock()
	for _, metric := range metrics {
		repo.gauges[metric.Name] = metric.Value
	}
	repo.gaugeMu.Unlock()
	return metrics, nil
}

// UpdateCounters updates or sets a new counter metrics with the given name and value.
func (repo *inMemoryRepository) UpdateCounters(metrics []repository.CounterMetric) ([]repository.CounterMetric, error) {
	newMetrics := make([]repository.CounterMetric, 0, len(metrics))
	repo.counterMu.Lock()
	for _, metric := range metrics {
		newValue := repo.counters[metric.Name] + metric.Value
		repo.counters[metric.Name] = metric.Value
		newMetrics = append(newMetrics, repository.CounterMetric{Name: metric.Name, Value: newValue})
	}
	repo.counterMu.Unlock()
	return newMetrics, nil
}

// GetGauge return gauge metric by name.
func (repo *inMemoryRepository) GetGauge(name string) (float64, error) {
	repo.gaugeMu.RLock()
	defer repo.gaugeMu.RUnlock()
	gauge, ok := repo.gauges[name]
	if !ok {
		return 0, repository.ErrMetricNotFound
	}
	return gauge, nil
}

// GetCounter return counter metric by name.
func (repo *inMemoryRepository) GetCounter(name string) (int64, error) {
	repo.counterMu.RLock()
	defer repo.counterMu.RUnlock()
	counter, ok := repo.counters[name]
	if !ok {
		return 0, repository.ErrMetricNotFound
	}
	return counter, nil
}

// GetAllGauges returns all gauge metrics.
func (repo *inMemoryRepository) GetAllGauges() ([]repository.GaugeMetric, error) {
	repo.gaugeMu.RLock()
	defer repo.gaugeMu.RUnlock()

	gauges := make([]repository.GaugeMetric, 0, len(repo.gauges))
	for name, value := range repo.gauges {
		gauges = append(gauges, repository.GaugeMetric{Name: name, Value: value})
	}

	return gauges, nil
}

// GetAllCounters returns all counter metrics.
func (repo *inMemoryRepository) GetAllCounters() ([]repository.CounterMetric, error) {
	repo.counterMu.RLock()
	defer repo.counterMu.RUnlock()

	counters := make([]repository.CounterMetric, 0, len(repo.counters))
	for name, value := range repo.counters {
		counters = append(counters, repository.CounterMetric{Name: name, Value: value})
	}

	return counters, nil
}

// DeleteGauge deletes gauge metric by name.
func (repo *inMemoryRepository) DeleteGauge(name string) error {
	repo.gaugeMu.Lock()
	delete(repo.gauges, name)
	repo.gaugeMu.Unlock()
	return nil
}

// DeleteCounter deletes counter metric by name.
func (repo *inMemoryRepository) DeleteCounter(name string) error {
	repo.counterMu.Lock()
	delete(repo.counters, name)
	repo.counterMu.Unlock()
	return nil
}

func (repo *inMemoryRepository) startSaveMetricsInBackground(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-repo.saveTicker.C:
				err := repo.fileStorage.SaveMetrics(repo.dumpMetrics())
				if err != nil {
					log.Fatalf("failed to save metrics through FileStorage: %v", err)
				}
			}
		}
	}()
}

func (repo *inMemoryRepository) dumpMetrics() filestorage.Metrics {
	repo.gaugeMu.RLock()
	gauges := make([]filestorage.GaugeMetric, 0, len(repo.gauges))
	for name, value := range repo.gauges {
		gauges = append(gauges, filestorage.GaugeMetric{Name: name, Value: value})
	}
	repo.gaugeMu.RUnlock()

	repo.counterMu.RLock()
	counters := make([]filestorage.CounterMetric, 0, len(repo.counters))
	for name, value := range repo.counters {
		counters = append(counters, filestorage.CounterMetric{Name: name, Value: value})
	}
	repo.counterMu.RUnlock()

	return filestorage.Metrics{
		Gauges:   gauges,
		Counters: counters,
	}
}

func gaugesToMap(gauges []filestorage.GaugeMetric) map[string]float64 {
	result := make(map[string]float64, len(gauges))
	for _, g := range gauges {
		result[g.Name] = g.Value
	}
	return result
}

func countersToMap(counters []filestorage.CounterMetric) map[string]int64 {
	result := make(map[string]int64, len(counters))
	for _, g := range counters {
		result[g.Name] = g.Value
	}
	return result
}
