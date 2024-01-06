package database

type Metric struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

type Metrics []Metric

type InMemory struct{
    data Metrics
}

func nextIdBuilder() func() int {
    id := 0
    return func() int {
        id = id + 1
        return id
    }
}

var nextId = nextIdBuilder()

func NewInMemory() *InMemory {
    return &InMemory{
        data: Metrics{
            {
                ID:          nextId(),
                Title:       "Metric 1",
                Unit:        "unit",
                Description: "description",
            },
            {
                ID:          nextId(),
                Title:       "Metric 2",
                Unit:        "unit",
                Description: "description",
            },
            {
                ID:          nextId(),
                Title:       "Metric 3",
                Unit:        "unit",
                Description: "description",
            },
        },
    }
}

func (db *InMemory) GetMetrics() ([]Metric, error) {
	return db.data, nil
}

func (db *InMemory) AddMetric(m Metric) error {
    m.ID = nextId()
    db.data = append(db.data, m)
    return nil
}
