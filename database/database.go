package database

import (
	"errors"
	"fmt"
)

type Metric struct {
	ID          int    `form:"id"`
	Title       string `form:"title"`
	Unit        string `form:"unit"`
	Description string `form:"description"`
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

func (db *InMemory) GetMetric(id int) (Metric, error) {
    for _, m := range db.data {
        if m.ID == id {
            return m, nil
        }
    }

    errorMsg := fmt.Sprintf("Metric %d not found", id)
    return Metric{}, errors.New(errorMsg)
}

func (db *InMemory) AddMetric(m Metric) (Metric, error) {
    m.ID = nextId()
    db.data = append(db.data, m)
    return m, nil
}
