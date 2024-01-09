package database_test

import (
	"testing"

	"github.com/belarte/metrix/database"
	"github.com/stretchr/testify/assert"
)

var (
	initial = database.Metrics{
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
	afterInsertion = database.Metrics{
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
	afterUpdate = database.Metrics{
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
	metricToCreate = database.Metric{
		ID:          4,
		Title:       "Metric 4",
		Unit:        "Unit 4",
		Description: "description",
	}
	metricToUpdate = database.Metric{
		ID:          2,
		Title:       "New Title",
		Unit:        "Unit 4",
		Description: "new description",
	}
)

func TestDatabaseAddMetric(t *testing.T) {
	tests := map[string]struct {
		metric  database.Metric
		input   database.Metrics
		expected database.Metrics
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
