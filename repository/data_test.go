package repository_test

import "github.com/belarte/metrix/model"

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
