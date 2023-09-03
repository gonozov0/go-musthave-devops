package storage

import "sync"

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
func (repo *InMemoryRepository) UpdateGauge(metricName string, value float64) error {
	repo.gaugeMu.Lock()
	defer repo.gaugeMu.Unlock()
	repo.gauges[metricName] = value
	return nil
}

// UpdateCounter updates or sets a new counter metric with the given name and value.
func (repo *InMemoryRepository) UpdateCounter(metricName string, value int64) error {
	repo.counterMu.Lock()
	defer repo.counterMu.Unlock()
	repo.counters[metricName] += value
	return nil
}

// GetGauges returns all gauge metrics.
func (repo *InMemoryRepository) GetGauges() (map[string]float64, error) {
	repo.gaugeMu.RLock()
	defer repo.gaugeMu.RUnlock()
	return repo.gauges, nil
}

// GetCounters returns all counter metrics.
func (repo *InMemoryRepository) GetCounters() (map[string]int64, error) {
	repo.counterMu.RLock()
	defer repo.counterMu.RUnlock()
	return repo.counters, nil
}
