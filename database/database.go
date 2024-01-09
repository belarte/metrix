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

func nextIdBuilder(id int) func() int {
    return func() int {
        id = id + 1
        return id
    }
}

var nextId = nextIdBuilder(0)

func NewInMemory() *InMemory {
    nextId = nextIdBuilder(0)

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

func (db *InMemory) UpsertMetric(metric Metric) (Metric, error) {
    for i, m := range db.data {
        if metric.ID == m.ID {
            db.data[i] = metric
            return metric, nil
        }
    }
    metric.ID = nextId()
    db.data = append(db.data, metric)
    return metric, nil
}
