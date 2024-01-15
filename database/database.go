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

type Entry struct {
	ID       int     `form:"id"`
	MetricID int     `form:"metric"`
	Value    float64 `form:"value"`
	Date     string  `form:"date"`
}

type Entries []Entry

type InMemory struct {
	metric  Metrics
	entries Entries
}

func nextIdBuilder(id int) func() int {
	return func() int {
		id = id + 1
		return id
	}
}

var nextMetricId = nextIdBuilder(0)
var nextEntryId = nextIdBuilder(0)

func NewInMemory() *InMemory {
	nextMetricId = nextIdBuilder(0)
	nextEntryId = nextIdBuilder(0)

	return &InMemory{
		metric: Metrics{
			{ID: nextMetricId(), Title: "Metric 1", Unit: "unit", Description: "description"},
			{ID: nextMetricId(), Title: "Metric 2", Unit: "unit", Description: "description"},
			{ID: nextMetricId(), Title: "Metric 3", Unit: "unit", Description: "description"},
		},
		entries: Entries{
			{ID: nextEntryId(), MetricID: 1, Value: 5.0, Date: "2018-01-01"},
			{ID: nextEntryId(), MetricID: 2, Value: 2.1, Date: "2018-01-11"},
			{ID: nextEntryId(), MetricID: 1, Value: 1.0, Date: "2018-01-15"},
		},
	}
}

func (db *InMemory) GetMetrics() ([]Metric, error) {
	return db.metric, nil
}

func (db *InMemory) GetMetric(id int) (Metric, error) {
	for _, m := range db.metric {
		if m.ID == id {
			return m, nil
		}
	}

	errorMsg := fmt.Sprintf("Metric %d not found", id)
	return Metric{}, errors.New(errorMsg)
}

func (db *InMemory) UpsertMetric(metric Metric) (Metric, error) {
	for i, m := range db.metric {
		if metric.ID == m.ID {
			db.metric[i] = metric
			return metric, nil
		}
	}
	metric.ID = nextMetricId()
	db.metric = append(db.metric, metric)
	return metric, nil
}

type DatabaseError struct {
	message string
}

func NewDatabaseError(message string) DatabaseError {
	return DatabaseError{message}
}

func (e DatabaseError) Error() string {
	return e.message
}

func (db *InMemory) GetEntries() ([]Entry, error) {
	return db.entries, nil
}

func (db *InMemory) UpsertEntry(metricId int, value float64, date string) (Entry, error) {
	found := false
	for _, m := range db.metric {
		found = found || (m.ID == metricId)
	}
	if !found {
		return Entry{}, DatabaseError{"metric not found"}
	}

	for i, e := range db.entries {
		if e.MetricID == metricId && e.Date == date {
			entry := Entry{e.ID, metricId, value, date}
			db.entries[i] = entry
			return entry, nil
		}
	}

	entry := Entry{nextEntryId(), metricId, value, date}
	db.entries = append(db.entries, entry)
	return entry, nil
}

func (db *InMemory) GetSortedEntriesForMetric(metricId int) ([]Entry, error) {
    entries := []Entry{}
    for _, e := range db.entries {
        if e.MetricID == metricId {
            entries = append(entries, e)
        }
    }
    return entries, nil
}
