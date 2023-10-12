package filestorage

// GaugeMetric is a struct that represents a gauge metric.
type GaugeMetric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// CounterMetric is a struct that represents a counter metric.
type CounterMetric struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// Metrics is a struct that represents a data to save metrics.
type Metrics struct {
	Gauges   []GaugeMetric
	Counters []CounterMetric
}
