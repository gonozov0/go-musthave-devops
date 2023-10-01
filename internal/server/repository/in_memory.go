package repository

import (
	"sync"
)

// InMemoryRepository is an in-memory implementation of the Repository interface.
type InMemoryRepository struct {
	gaugeMu   sync.RWMutex
	counterMu sync.RWMutex
	gauges    map[string]float64
	counters  map[string]int64
}

// NewInMemoryRepository creates a new InMemoryRepository and returns it as a Repository interface.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// UpdateGauge updates or sets a new gauge metric with the given name and value.
func (repo *InMemoryRepository) UpdateGauge(metricName string, value float64) (float64, error) {
	repo.gaugeMu.Lock()
	defer repo.gaugeMu.Unlock()
	repo.gauges[metricName] = value
	return value, nil
}

// UpdateCounter updates or sets a new counter metric with the given name and value.
func (repo *InMemoryRepository) UpdateCounter(metricName string, value int64) (int64, error) {
	repo.counterMu.Lock()
	defer repo.counterMu.Unlock()
	newValue := repo.counters[metricName] + value
	repo.counters[metricName] = newValue
	return newValue, nil
}

// GetGauge return gauge metric by name.
func (repo *InMemoryRepository) GetGauge(name string) (float64, error) {
	repo.gaugeMu.RLock()
	defer repo.gaugeMu.RUnlock()
	gauge, ok := repo.gauges[name]
	if !ok {
		return 0, ErrMetricNotFound
	}
	return gauge, nil
}

// GetCounter return counter metric by name.
func (repo *InMemoryRepository) GetCounter(name string) (int64, error) {
	repo.counterMu.RLock()
	defer repo.counterMu.RUnlock()
	counter, ok := repo.counters[name]
	if !ok {
		return 0, ErrMetricNotFound
	}
	return counter, nil
}

// GetAllGauges returns all gauge metrics.
func (repo *InMemoryRepository) GetAllGauges() ([]GaugeMetric, error) {
	repo.gaugeMu.RLock()
	defer repo.gaugeMu.RUnlock()

	gauges := make([]GaugeMetric, 0, len(repo.gauges))
	for name, value := range repo.gauges {
		gauges = append(gauges, GaugeMetric{Name: name, Value: value})
	}

	return gauges, nil
}

// GetAllCounters returns all counter metrics.
func (repo *InMemoryRepository) GetAllCounters() ([]CounterMetric, error) {
	repo.counterMu.RLock()
	defer repo.counterMu.RUnlock()

	counters := make([]CounterMetric, 0, len(repo.counters))
	for name, value := range repo.counters {
		counters = append(counters, CounterMetric{Name: name, Value: value})
	}

	return counters, nil
}
