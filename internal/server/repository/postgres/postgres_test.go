package repository

import (
	"errors"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

var (
	countersNames = []string{"test_counter", "test_counter_all"}
	gaugesNames   = []string{"test_gauge", "test_gauge_all"}
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
	assert.NoError(s.T(), err)
}

func (s *PGRepositorySuite) TestUpdateAndGetGauge() {
	metricName := gaugesNames[0]
	value := 42.5
	updatedValue, err := s.repo.UpdateGauge(metricName, value)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), value, updatedValue)

	fetchedValue, err := s.repo.GetGauge(metricName)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), value, fetchedValue)
}

func (s *PGRepositorySuite) TestUpdateAndGetCounter() {
	metricName := countersNames[0]
	value := int64(5)
	updatedValue, err := s.repo.UpdateCounter(metricName, value)
	assert.NoError(s.T(), err)
	assert.Condition(s.T(), func() bool {
		return updatedValue >= value
	})

	fetchedValue, err := s.repo.GetCounter(metricName)
	assert.NoError(s.T(), err)
	assert.Condition(s.T(), func() bool {
		return fetchedValue >= value
	})
}

func (s *PGRepositorySuite) TestGetGaugeMetricNotFound() {
	_, err := s.repo.GetGauge("non_existing_gauge")
	assert.True(s.T(), errors.Is(err, repository.ErrMetricNotFound))
}

func (s *PGRepositorySuite) TestGetAllGauges() {
	value := 42.5
	for _, gauge := range gaugesNames {
		_, err := s.repo.UpdateGauge(gauge, value)
		assert.NoError(s.T(), err)
	}

	gauges, err := s.repo.GetAllGauges()
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), gauges)
	assert.True(s.T(), len(gauges) > 0)
	for _, gaugeName := range gaugesNames {
		assert.Contains(s.T(), gauges, repository.GaugeMetric{Name: gaugeName, Value: value})
	}
}

func (s *PGRepositorySuite) TestGetCounterMetricNotFound() {
	_, err := s.repo.GetCounter("non_existing_counter")
	assert.True(s.T(), errors.Is(err, repository.ErrMetricNotFound))
}

func (s *PGRepositorySuite) TestGetAllCounters() {
	value := int64(5)
	for _, counter := range countersNames {
		_, err := s.repo.UpdateCounter(counter, value)
		assert.NoError(s.T(), err)
	}

	counters, err := s.repo.GetAllCounters()
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), counters)
	assert.True(s.T(), len(counters) > 0)
	for _, counterName := range countersNames {
		for _, counter := range counters {
			if counter.Name == counterName {
				assert.Condition(s.T(), func() bool {
					return counter.Value >= value
				})
				break
			}
		}
	}
}

func (s *PGRepositorySuite) TearDownSuite() {
	for _, gaugeName := range gaugesNames {
		err := s.repo.DeleteGauge(gaugeName)
		assert.NoError(s.T(), err)
	}

	for _, counterName := range countersNames {
		err := s.repo.DeleteCounter(counterName)
		assert.NoError(s.T(), err)
	}
}
