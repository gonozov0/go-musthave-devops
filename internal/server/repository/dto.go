package repository

// GaugeMetric is a struct that represents a gauge metric.
type GaugeMetric struct {
	Name  string  `db:"name"`
	Value float64 `db:"value"`
}

// CounterMetric is a struct that represents a counter metric.
type CounterMetric struct {
	Name  string `db:"name"`
	Value int64  `db:"value"`
}
