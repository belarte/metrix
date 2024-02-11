package repository_test

import (
	"testing"

	"github.com/belarte/metrix/model"
	"github.com/belarte/metrix/repository"
	"github.com/stretchr/testify/suite"
)

const (
	inMemoryDatabase = ":memory:"
)

type RepositoryTestSuite struct {
	suite.Suite
	db *repository.Repository
}

func (s *RepositoryTestSuite) SetupTest() {
	db, err := repository.New(inMemoryDatabase)
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

func (s *RepositoryTestSuite) TestGetEntriesForMetricInOrder() {
	metric := model.Metric{Title: "Title", Unit: "Unit", Description: "Description"}

	m, err := s.db.UpsertMetric(metric)
	s.NoError(err)

	metricId := m.ID
	_, err = s.db.UpsertEntry(metricId, 1.0, "2020-01-05")
	s.NoError(err)
	_, err = s.db.UpsertEntry(metricId, 2.0, "2020-01-04")
	s.NoError(err)
	_, err = s.db.UpsertEntry(metricId, 3.0, "2020-01-07")
	s.NoError(err)
	_, err = s.db.UpsertEntry(metricId, 4.0, "2020-01-05")
	s.NoError(err)
	_, err = s.db.UpsertEntry(metricId, 5.0, "2020-01-09")
	s.NoError(err)

	metrics, err := s.db.GetSortedEntriesForMetric(metricId)
	s.NoError(err)
	s.Equal(4, len(metrics))
	s.Equal(2.0, metrics[0].Value)
	s.Equal(4.0, metrics[1].Value)
	s.Equal(3.0, metrics[2].Value)
	s.Equal(5.0, metrics[3].Value)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
