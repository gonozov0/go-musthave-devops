package storage

// Metric types as string constants
const (
	Gauge   = "gauge"
	Counter = "counter"
)

// Repository describes the behavior for storing metrics.
type Repository interface {
	// UpdateGauge updates or adds a new gauge metric with the given name and value.
	UpdateGauge(metricName string, value float64) error
	// UpdateCounter updates or adds a new counter metric with the given name and value.
	UpdateCounter(metricName string, value int64) error
}
