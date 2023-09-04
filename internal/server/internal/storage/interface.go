package storage

// Repository describes the behavior for storing metrics.
type Repository interface {
	// UpdateGauge updates or adds a new gauge metric with the given name and value.
	UpdateGauge(metricName string, value float64) error
	// UpdateCounter updates or adds a new counter metric with the given name and value.
	UpdateCounter(metricName string, value int64) error
	// GetGauge return gauge metric by name.
	GetGauge(name string) (float64, error)
	// GetCounter return counter metric by name.
	GetCounter(name string) (int64, error)
}
