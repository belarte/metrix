package repository_test

import "github.com/belarte/metrix/model"

var (
	initialMetrics = model.Metrics{
		{
			Title:       "Metric 1",
			Unit:        "unit",
			Description: "description",
		},
		{
			Title:       "Metric 2",
			Unit:        "unit",
			Description: "description",
		},
		{
			Title:       "Metric 3",
			Unit:        "unit",
			Description: "description",
		},
	}
	initialEntries = model.Entries{
		model.NewEntry(1, 5.0, "2018-01-01"),
		model.NewEntry(2, 2.1, "2018-01-11"),
		model.NewEntry(1, 1.0, "2018-01-15"),
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
