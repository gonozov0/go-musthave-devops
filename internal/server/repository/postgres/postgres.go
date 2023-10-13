package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

type pgRepository struct {
	db *sqlx.DB
}

func NewPGRepository(connectionString string) (repository.Repository, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(getMigrationDirPath(), connectionString)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &pgRepository{db: db}, nil
}

func (r *pgRepository) Ping() error {
	return r.db.Ping()
}

func (r *pgRepository) UpdateGauge(metricName string, value float64) (float64, error) {
	_, err := r.db.Exec(
		`INSERT INTO gauges(name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2`,
		metricName,
		value,
	)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (r *pgRepository) UpdateGauges(metrics []repository.GaugeMetric) ([]repository.GaugeMetric, error) {
	if len(metrics) == 0 {
		return make([]repository.GaugeMetric, 0), nil
	}

	// Name can be duplicated in slice, so we need to fix it
	uniqueMetrics := make(map[string]float64, len(metrics))
	for _, metric := range metrics {
		uniqueMetrics[metric.Name] = metric.Value
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	var valueStrings []string
	for metricName := range uniqueMetrics {
		valueStrings = append(valueStrings, fmt.Sprintf("(:name%s, :value%s)", metricName, metricName))
	}

	queryStr := fmt.Sprintf(
		"INSERT INTO gauges(name, value) VALUES %s ON CONFLICT(name) DO UPDATE SET value = EXCLUDED.value RETURNING name, value",
		strings.Join(valueStrings, ","),
	)

	queryArgs := make(map[string]interface{})
	for name, value := range uniqueMetrics {
		queryArgs[fmt.Sprintf("name%s", name)] = name
		queryArgs[fmt.Sprintf("value%s", name)] = value
	}

	rows, err := tx.NamedQuery(queryStr, queryArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute named query: %w", err)
	}
	defer rows.Close()

	var updatedMetrics []repository.GaugeMetric
	for rows.Next() {
		var metric repository.GaugeMetric
		if err := rows.StructScan(&metric); err != nil {
			return nil, fmt.Errorf("failed to scan struct: %w", err)
		}
		updatedMetrics = append(updatedMetrics, metric)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedMetrics, nil
}

func (r *pgRepository) GetGauge(name string) (float64, error) {
	var value float64
	err := r.db.Get(&value, `SELECT value FROM gauges WHERE name = $1`, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repository.ErrMetricNotFound
		}
		return 0, err
	}
	return value, nil
}

func (r *pgRepository) GetAllGauges() ([]repository.GaugeMetric, error) {
	var metrics []repository.GaugeMetric
	err := r.db.Select(&metrics, `SELECT name, value FROM gauges g`)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *pgRepository) UpdateCounter(metricName string, value int64) (int64, error) {
	query := `
		INSERT INTO counters(name, value) 
		VALUES ($1, $2) 
		ON CONFLICT(name) 
		DO UPDATE SET value = counters.value + $2 
		RETURNING value
	`
	var newValue int64
	err := r.db.QueryRow(query, metricName, value).Scan(&newValue)
	if err != nil {
		return 0, err
	}
	return newValue, nil
}

func (r *pgRepository) UpdateCounters(metrics []repository.CounterMetric) ([]repository.CounterMetric, error) {
	if len(metrics) == 0 {
		return make([]repository.CounterMetric, 0), nil
	}

	// Name can be duplicated in slice, so we need to fix it
	uniqueMetrics := make(map[string]int64, len(metrics))
	for _, metric := range metrics {
		uniqueMetrics[metric.Name] += metric.Value
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	var valueStrings []string
	for metricName := range uniqueMetrics {
		valueStrings = append(valueStrings, fmt.Sprintf("(:name%s, :value%s)", metricName, metricName))
	}

	queryStr := fmt.Sprintf(
		"INSERT INTO counters(name, value) VALUES %s ON CONFLICT(name) DO UPDATE SET value = counters.value + EXCLUDED.value RETURNING name, value",
		strings.Join(valueStrings, ","),
	)

	queryArgs := make(map[string]interface{})
	for name, value := range uniqueMetrics {
		queryArgs[fmt.Sprintf("name%s", name)] = name
		queryArgs[fmt.Sprintf("value%s", name)] = value
	}

	rows, err := tx.NamedQuery(queryStr, queryArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute named query: %w", err)
	}
	defer rows.Close()

	var updatedMetrics []repository.CounterMetric
	for rows.Next() {
		var metric repository.CounterMetric
		if err := rows.StructScan(&metric); err != nil {
			return nil, fmt.Errorf("failed to scan struct: %w", err)
		}
		updatedMetrics = append(updatedMetrics, metric)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedMetrics, nil
}

func (r *pgRepository) GetCounter(name string) (int64, error) {
	var value int64
	err := r.db.Get(&value, `SELECT value FROM counters WHERE name = $1`, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repository.ErrMetricNotFound
		}
		return 0, err
	}
	return value, nil
}

func (r *pgRepository) GetAllCounters() ([]repository.CounterMetric, error) {
	var metrics []repository.CounterMetric
	err := r.db.Select(&metrics, `SELECT name, value FROM counters`)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *pgRepository) DeleteGauge(name string) error {
	_, err := r.db.Exec(`DELETE FROM gauges WHERE name = $1`, name)
	return err
}

func (r *pgRepository) DeleteCounter(name string) error {
	_, err := r.db.Exec(`DELETE FROM counters WHERE name = $1`, name)
	return err
}

func getMigrationDirPath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath)
	migrationsPath := filepath.Join(currentDir, "internal/migrations")
	return "file://" + migrationsPath
}
