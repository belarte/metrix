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

func populateDatabase(db *repository.Repository) error {
	for _, metric := range initialMetrics {
		_, err := db.UpsertMetric(metric)
		if err != nil {
			return err
		}
	}

	for _, entry := range initialEntries {
		_, err := db.UpsertEntry(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RepositoryTestSuite) SetupTest() {
	db, err := repository.New(inMemoryDatabase)
	s.NoError(err)

	err = db.Migrate()
	s.NoError(err)

	err = populateDatabase(db)
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

func (s *RepositoryTestSuite) TestDeleteMetric() {
	_, err := s.db.GetMetric(1)
	s.NoError(err)

	entriesBefore, err := s.db.GetSortedEntriesForMetric(1)
	s.NoError(err)

	err = s.db.DeleteMetric(1)
	s.NoError(err)

	_, err = s.db.GetMetric(1)
	s.Error(err)

	entriesAfter, err := s.db.GetSortedEntriesForMetric(1)
	s.NoError(err)

	s.Equal(2, len(entriesBefore))
	s.Equal(0, len(entriesAfter))
}

func (s *RepositoryTestSuite) TestCannotDeleteNonExistingMetric() {
	err := s.db.DeleteMetric(999)
	s.Error(err)
}

func (s *RepositoryTestSuite) TestGetExistingEntry() {
	entry, err := s.db.GetEntry(1, "2018-01-01")
	s.NoError(err)
	s.Equal(initialEntries[0], entry)
}

func (s *RepositoryTestSuite) TestGetNonExistingEntry() {
	_, err := s.db.GetEntry(999, "2018-01-01")
	s.Error(err)
}

func (s *RepositoryTestSuite) TestAddNewEntry() {
	input := model.NewEntry(1, 1.0, "2018-02-01")
	expectedEntry := model.NewEntry(1, 1.0, "2018-02-01")
	expectedSize := 4

	entry, err := s.db.UpsertEntry(input)
	s.NoError(err)

	entries, err := s.db.GetEntries()
	s.NoError(err)

	s.Equal(expectedEntry, entry)
	s.Equal(expectedSize, len(entries))
}

func (s *RepositoryTestSuite) TestUpdateEntry() {
	input := model.NewEntry(1, 7.0, "2018-01-01")
	expectedEntry := model.NewEntry(1, 7.0, "2018-01-01")
	expectedSize := 3

	entry, err := s.db.UpsertEntry(input)
	s.NoError(err)

	entries, err := s.db.GetEntries()
	s.NoError(err)

	s.Equal(expectedEntry, entry)
	s.Equal(expectedSize, len(entries))
}

func (s *RepositoryTestSuite) TestAddEntryWhenMetricDoesNotExist() {
	input := model.NewEntry(-1, 1.0, "2018-02-01")
	expectedEntry := model.Entry{}

	entry, err := s.db.UpsertEntry(input)
	s.Error(err)
	s.Equal(expectedEntry, entry)
}

func (s *RepositoryTestSuite) TestDeleteExistingEntry() {
	entries, err := s.db.GetSortedEntriesForMetric(1)
	s.NoError(err)
	s.Equal(2, len(entries))

	err = s.db.DeleteEntry(1, "2018-01-01")
	s.NoError(err)

	entries, err = s.db.GetSortedEntriesForMetric(1)
	s.NoError(err)
	s.Equal(1, len(entries))
}

func (s *RepositoryTestSuite) TestCannotDeleteNonExistingEntry() {
	err := s.db.DeleteEntry(999, "2184-11-27")
	s.Error(err)
}

func (s *RepositoryTestSuite) TestGetEntriesForMetricInOrder() {
	metric := model.Metric{Title: "Title", Unit: "Unit", Description: "Description"}

	m, err := s.db.UpsertMetric(metric)
	s.NoError(err)

	metricId := m.ID
	_, err = s.db.UpsertEntry(model.NewEntry(metricId, 1.0, "2020-01-05"))
	s.NoError(err)
	_, err = s.db.UpsertEntry(model.NewEntry(metricId, 2.0, "2020-01-04"))
	s.NoError(err)
	_, err = s.db.UpsertEntry(model.NewEntry(metricId, 3.0, "2020-01-07"))
	s.NoError(err)
	_, err = s.db.UpsertEntry(model.NewEntry(metricId, 4.0, "2020-01-05"))
	s.NoError(err)
	_, err = s.db.UpsertEntry(model.NewEntry(metricId, 5.0, "2020-01-09"))
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
