package repository

import (
	"database/sql"
	"errors"

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

	migrationsPath := "file://internal/server/repository/postgres/internal/migrations"
	m, err := migrate.New(migrationsPath, connectionString)
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
