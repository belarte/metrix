package database

type Metric struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

type InMemory struct{}

func (db *InMemory) GetMetrics() ([]Metric, error) {
	return []Metric{
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
	}, nil
}
