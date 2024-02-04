package database

import (
	"errors"
	"fmt"

	"github.com/belarte/metrix/model"
)

type InMemory struct {
	metric  model.Metrics
	entries model.Entries
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
		metric: model.Metrics{
			{ID: nextMetricId(), Title: "Metric 1", Unit: "unit", Description: "description"},
			{ID: nextMetricId(), Title: "Metric 2", Unit: "unit", Description: "description"},
			{ID: nextMetricId(), Title: "Metric 3", Unit: "unit", Description: "description"},
		},
		entries: model.Entries{
			{ID: nextEntryId(), MetricID: 1, Value: 5.0, Date: "2018-01-01"},
			{ID: nextEntryId(), MetricID: 2, Value: 2.1, Date: "2018-01-11"},
			{ID: nextEntryId(), MetricID: 1, Value: 1.0, Date: "2018-01-15"},
		},
	}
}

func (db *InMemory) GetMetrics() ([]model.Metric, error) {
	return db.metric, nil
}

func (db *InMemory) GetMetric(id int) (model.Metric, error) {
	for _, m := range db.metric {
		if m.ID == id {
			return m, nil
		}
	}

	errorMsg := fmt.Sprintf("Metric %d not found", id)
	return model.Metric{}, errors.New(errorMsg)
}

func (db *InMemory) UpsertMetric(metric model.Metric) (model.Metric, error) {
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

func (db *InMemory) GetEntries() ([]model.Entry, error) {
	return db.entries, nil
}

func (db *InMemory) UpsertEntry(metricId int, value float64, date string) (model.Entry, error) {
	found := false
	for _, m := range db.metric {
		found = found || (m.ID == metricId)
	}
	if !found {
		return model.Entry{}, DatabaseError{"metric not found"}
	}

	for i, e := range db.entries {
		if e.MetricID == metricId && e.Date == date {
			entry := model.Entry{
				ID:       e.ID,
				MetricID: metricId,
				Value:    value,
				Date:     date,
			}
			db.entries[i] = entry
			return entry, nil
		}
	}

	entry := model.Entry{
		ID:       nextEntryId(),
		MetricID: metricId,
		Value:    value,
		Date:     date,
	}
	db.entries = append(db.entries, entry)
	return entry, nil
}

func (db *InMemory) GetSortedEntriesForMetric(metricId int) ([]model.Entry, error) {
	entries := []model.Entry{}
	for _, e := range db.entries {
		if e.MetricID == metricId {
			entries = append(entries, e)
		}
	}
	return entries, nil
}
