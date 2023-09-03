package storage

import "sync"

// InMemoryRepository is an in-memory implementation of the Repository interface.
type InMemoryRepository struct {
	gaugeMu   sync.Mutex
	counterMu sync.Mutex
	gauges    map[string]float64
	counters  map[string]int64
}

// NewInMemoryRepository creates a new InMemoryRepository and returns it as a Repository interface.
func NewInMemoryRepository() Repository {
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
