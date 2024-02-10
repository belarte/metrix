package repository_test

import (
	"testing"

	"github.com/belarte/metrix/repository"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	db *repository.Repository
}

func (s *RepositoryTestSuite) SetupTest() {
	db, err := repository.New()
	s.NoError(err)

	err = db.Migrate()
	s.NoError(err)

	s.db = db
}

func (s *RepositoryTestSuite) TearDownTest() {
	err := s.db.Close()
	s.NoError(err)
}

func (s *RepositoryTestSuite) TestAddMetric() {
	_, err := s.db.UpsertMetric(metricToCreate)
	s.NoError(err)

	metrics, err := s.db.GetMetrics()
	s.NoError(err)
	s.ElementsMatch(metrics, afterInsertion)
}

func (s *RepositoryTestSuite) TestUpdateMetric() {
	_, err := s.db.UpsertMetric(metricToUpdate)
	s.NoError(err)

	metrics, err := s.db.GetMetrics()
	s.NoError(err)
	s.ElementsMatch(metrics, afterUpdate)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
