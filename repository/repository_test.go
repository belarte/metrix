package repository_test

import (
	"testing"

	"github.com/belarte/metrix/model"
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

func (s *RepositoryTestSuite) TestAddNewEntry() {
	id := 1
	value := 1.0
	date := "2018-02-01"
	expectedEntry := model.NewEntry(1, 1.0, "2018-02-01")
	expectedSize := 4

	entry, err := s.db.UpsertEntry(id, value, date)
	s.NoError(err)

	entries, err := s.db.GetEntries()
	s.NoError(err)

	s.Equal(expectedEntry, entry)
	s.Equal(expectedSize, len(entries))
}

func (s *RepositoryTestSuite) TestUpdateEntry() {
	id := 1
	value := 7.0
	date := "2018-01-01"
	expectedEntry := model.NewEntry(1, 7.0, "2018-01-01")
	expectedSize := 3

	entry, err := s.db.UpsertEntry(id, value, date)
	s.NoError(err)

	entries, err := s.db.GetEntries()
	s.NoError(err)

	s.Equal(expectedEntry, entry)
	s.Equal(expectedSize, len(entries))
}

func (s *RepositoryTestSuite) TestAddEntryWhenMetricDoesNotExist() {
	id := -1
	value := 1.0
	date := "2018-02-01"
	expectedEntry := model.Entry{}

	entry, err := s.db.UpsertEntry(id, value, date)
	s.Error(err)
	s.Equal(expectedEntry, entry)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
