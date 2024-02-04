package database_test

import (
	"testing"

	"github.com/belarte/metrix/database"
	"github.com/belarte/metrix/model"
	"github.com/stretchr/testify/assert"
)

var (
	initial = model.Metrics{
		{
			ID:          1,
			Title:       "Metric 1",
			Unit:        "unit",
			Description: "description",
		},
		{
			ID:          2,
			Title:       "Metric 2",
			Unit:        "unit",
			Description: "description",
		},
		{
			ID:          3,
			Title:       "Metric 3",
			Unit:        "unit",
			Description: "description",
		},
	}
	afterInsertion = model.Metrics{
		{
			ID:          1,
			Title:       "Metric 1",
			Unit:        "unit",
			Description: "description",
		},
		{
			ID:          2,
			Title:       "Metric 2",
			Unit:        "unit",
			Description: "description",
		},
		{
			ID:          3,
			Title:       "Metric 3",
			Unit:        "unit",
			Description: "description",
		},
		{
			ID:          4,
			Title:       "Metric 4",
			Unit:        "Unit 4",
			Description: "description",
		},
	}
	afterUpdate = model.Metrics{
		{
			ID:          1,
			Title:       "Metric 1",
			Unit:        "unit",
			Description: "description",
		},
		{
			ID:          2,
			Title:       "New Title",
			Unit:        "Unit 4",
			Description: "new description",
		},
		{
			ID:          3,
			Title:       "Metric 3",
			Unit:        "unit",
			Description: "description",
		},
	}
	metricToCreate = model.Metric{
		ID:          4,
		Title:       "Metric 4",
		Unit:        "Unit 4",
		Description: "description",
	}
	metricToUpdate = model.Metric{
		ID:          2,
		Title:       "New Title",
		Unit:        "Unit 4",
		Description: "new description",
	}
)

func TestDatabaseAddMetric(t *testing.T) {
	tests := map[string]struct {
		metric   model.Metric
		input    model.Metrics
		expected model.Metrics
	}{
		"create a new metric": {
			metricToCreate,
			initial,
			afterInsertion,
		},
		"update a new metric": {
			metricToUpdate,
			initial,
			afterUpdate,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := database.NewInMemory()

			_, err := db.UpsertMetric(test.metric)
			assert.NoError(t, err)

			metrics, err := db.GetMetrics()
			assert.NoError(t, err)
			assert.ElementsMatch(t, metrics, test.expected)
		})
	}
}

func TestDatabaseAddEntry(t *testing.T) {
	tests := map[string]struct {
		metricId int
		value    float64
		date     string
		result   model.Entry
		size     int
		err      error
	}{
		"add a new entry": {
			1, 1.0, "2018-02-01", model.NewEntry(4, 1, 1.0, "2018-02-01"), 4, nil,
		},
		"metric does not exist": {
			-1, 1.0, "2018-02-01", model.Entry{}, 3, database.NewDatabaseError("metric not found"),
		},
		"metric already entered for that date": {
			1, 7.0, "2018-01-01", model.NewEntry(1, 1, 7.0, "2018-01-01"), 3, nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := database.NewInMemory()
			entry, err := db.UpsertEntry(test.metricId, test.value, test.date)
			entries, _ := db.GetEntries()
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.result, entry)
			assert.Equal(t, test.size, len(entries))
		})
	}
}
