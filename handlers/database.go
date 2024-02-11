package handlers

import "github.com/belarte/metrix/model"

type Database interface {
	GetMetrics() (model.Metrics, error)
	GetMetric(id int) (model.Metric, error)
	UpsertMetric(metric model.Metric) (model.Metric, error)
	UpsertEntry(entry model.Entry) (model.Entry, error)
	GetSortedEntriesForMetric(id int) (model.Entries, error)
}
