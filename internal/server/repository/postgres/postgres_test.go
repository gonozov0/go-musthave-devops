package repository

import (
	"errors"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

type PGRepositorySuite struct {
	suite.Suite
	repo             repository.Repository
	connectionString string
}

func TestPGRepositorySuite(t *testing.T) {
	cfg, err := server.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config: %s", err.Error())
	}
	if cfg.DatabaseDSN == "" {
		t.Skip("Skipping repository tests. Set DATABASE_DSN to run them.")
	}
	suite.Run(t, &PGRepositorySuite{
		connectionString: cfg.DatabaseDSN,
	})
}

func (s *PGRepositorySuite) SetupSuite() {
	var err error
	s.repo, err = NewPGRepository(s.connectionString)
	if err != nil {
		s.T().Fatalf("Could not init postgres repository: %s", err.Error())
	}
}

func (s *PGRepositorySuite) TestPing() {
	err := s.repo.Ping()
	require.NoError(s.T(), err)
}

func (s *PGRepositorySuite) TestUpdateAndGetGauge() {
	metricName := "test_gauge"
	value := 42.5

	updatedValue, err := s.repo.UpdateGauge(metricName, value)
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, updatedValue)
	// check idempotency
	updatedValue, err = s.repo.UpdateGauge(metricName, value)
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, updatedValue)

	fetchedValue, err := s.repo.GetGauge(metricName)
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, fetchedValue)

	err = s.repo.DeleteGauge(metricName)
	require.NoError(s.T(), err)
}

func (s *PGRepositorySuite) TestUpdateAndGetCounter() {
	metricName := "test_counter"
	value := int64(5)

	updatedValue, err := s.repo.UpdateCounter(metricName, value)
	require.NoError(s.T(), err)
	require.Equal(s.T(), value, updatedValue)

	updatedValue, err = s.repo.UpdateCounter(metricName, value)
	require.NoError(s.T(), err)
	require.Equal(s.T(), value*2, updatedValue)

	fetchedValue, err := s.repo.GetCounter(metricName)
	require.NoError(s.T(), err)
	require.Equal(s.T(), value*2, fetchedValue)

	err = s.repo.DeleteCounter(metricName)
	require.NoError(s.T(), err)
}

func (s *PGRepositorySuite) TestGetGaugeMetricNotFound() {
	_, err := s.repo.GetGauge("non_existing_gauge")
	require.True(s.T(), errors.Is(err, repository.ErrMetricNotFound))
}

func (s *PGRepositorySuite) TestUpdateAndGetAllGauges() {
	expectedGauges := []repository.GaugeMetric{
		{Name: "test_gauge_1", Value: 42.5},
		{Name: "test_gauge_2", Value: 111.},
	}
	actualGauges, err := s.repo.UpdateGauges(expectedGauges)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedGauges, actualGauges)

	actualGauges, err = s.repo.GetAllGauges()
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedGauges, actualGauges)

	for _, gauge := range expectedGauges {
		err = s.repo.DeleteGauge(gauge.Name)
		require.NoError(s.T(), err)
	}
}

func (s *PGRepositorySuite) TestGetCounterMetricNotFound() {
	_, err := s.repo.GetCounter("non_existing_counter")
	require.True(s.T(), errors.Is(err, repository.ErrMetricNotFound))
}

func (s *PGRepositorySuite) TestUpdateAndGetAllCounters() {
	expectedCounters := []repository.CounterMetric{
		{Name: "test_counter_1", Value: 42},
		{Name: "test_counter_2", Value: 111},
	}
	actualCounters, err := s.repo.UpdateCounters(expectedCounters)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedCounters, actualCounters)

	actualCounters, err = s.repo.GetAllCounters()
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedCounters, actualCounters)

	for _, counter := range expectedCounters {
		err = s.repo.DeleteCounter(counter.Name)
		require.NoError(s.T(), err)
	}
}
