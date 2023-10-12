package repository

// Repository describes the behavior for storing metrics.
type Repository interface {
	// UpdateGauge updates or adds a new gauge metric with the given name and value.
	UpdateGauge(metricName string, value float64) (float64, error)
	// UpdateCounter updates or adds a new counter metric with the given name and value.
	UpdateCounter(metricName string, value int64) (int64, error)
	// GetGauge return gauge metric by name.
	GetGauge(name string) (float64, error)
	// GetCounter return counter metric by name.
	GetCounter(name string) (int64, error)
	// GetAllGauges returns all gauge metrics.
	GetAllGauges() ([]GaugeMetric, error)
	// GetAllCounters returns all counter metrics.
	GetAllCounters() ([]CounterMetric, error)
	// DeleteGauge deletes gauge metric by name.
	DeleteGauge(name string) error
	// DeleteCounter deletes counter metric by name.
	DeleteCounter(name string) error
	// Ping checks the connection to the repository.
	Ping() error
}
