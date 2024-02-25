package repository

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/belarte/metrix/model"
	_ "modernc.org/sqlite"
)

//go:embed sql/0001__create_tables.sql
var schema string

type Repository struct {
	db *sql.DB
}

func New(file string) (*Repository, error) {
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	return &Repository{db: db}, nil
}

func (d *Repository) Migrate() error {
	if _, err := d.db.Exec(schema); err != nil {
		return fmt.Errorf("error creating schema: %w", err)
	}
	return nil
}

func (d *Repository) Close() error {
	return d.db.Close()
}

func (d *Repository) GetMetrics() (model.Metrics, error) {
	rows, err := d.db.Query("SELECT id, title, unit, description FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("error querying metrics: %w", err)
	}
	defer rows.Close()

	var metrics []model.Metric
	for rows.Next() {
		var metric model.Metric
		if err = rows.Scan(&metric.ID, &metric.Title, &metric.Unit, &metric.Description); err != nil {
			return nil, fmt.Errorf("error scanning metric: %w", err)
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (d *Repository) GetMetric(id int) (model.Metric, error) {
	var metric model.Metric
	err := d.db.QueryRow("SELECT id, title, unit, description FROM metrics WHERE id = ?", id).
		Scan(&metric.ID, &metric.Title, &metric.Unit, &metric.Description)

	if err != nil {
		return model.Metric{}, fmt.Errorf("error querying metric: %w", err)
	}
	return metric, nil
}

func (d *Repository) UpsertMetric(metric model.Metric) (model.Metric, error) {
	if metric.ID == 0 {
		res, err := d.db.Exec(
			"INSERT INTO metrics (title, unit, description) VALUES (?, ?, ?) RETURNING id",
			&metric.Title, &metric.Unit, &metric.Description,
		)
		if err != nil {
			return model.Metric{}, fmt.Errorf("error inserting metric: %w", err)
		}
		var id int64
		if id, err = res.LastInsertId(); err != nil {
			return model.Metric{}, fmt.Errorf("error getting last insert ID: %w", err)
		}
		metric.ID = int(id)
	} else {
		_, err := d.db.Exec(
			"UPDATE metrics SET title = ?, unit = ?, description = ? WHERE id = ?",
			&metric.Title, &metric.Unit, &metric.Description, &metric.ID,
		)
		if err != nil {
			return model.Metric{}, fmt.Errorf("error updating metric: %w", err)
		}
	}

	return metric, nil
}

func (d *Repository) DeleteMetric(id int) error {
	res, err := d.db.Exec("DELETE FROM metrics WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error querying the database to delete the metric: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("error deleting metric, rows affected = %d: %w", count, err)
	}
	return nil
}

func (d *Repository) GetEntry(metricId int, date string) (model.Entry, error) {
	row := d.db.QueryRow("SELECT metric_id, value, date FROM entries WHERE metric_id = ? AND date = ?", metricId, date)

	var entry model.Entry
	if err := row.Scan(&entry.MetricID, &entry.Value, &entry.Date); err != nil {
		return model.Entry{}, fmt.Errorf("error scanning entry: %w", err)
	}

	return entry, nil
}

func (d *Repository) GetEntries() (model.Entries, error) {
	rows, err := d.db.Query("SELECT metric_id, value, date FROM entries")
	if err != nil {
		return nil, fmt.Errorf("error querying entries: %w", err)
	}
	defer rows.Close()

	var entries []model.Entry
	for rows.Next() {
		var entry model.Entry
		if err = rows.Scan(&entry.MetricID, &entry.Value, &entry.Date); err != nil {
			return nil, fmt.Errorf("error scanning entry: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (d *Repository) UpsertEntry(entry model.Entry) (model.Entry, error) {
	_, err := d.db.Exec(
		"INSERT OR REPLACE INTO entries (metric_id, value, date) VALUES (?, ?, ?)",
		&entry.MetricID, &entry.Value, &entry.Date,
	)

	if err != nil {
		return model.Entry{}, fmt.Errorf("error inserting entry: %w", err)
	}

	return entry, nil
}

func (d *Repository) DeleteEntry(metricId int, date string) error {
	res, err := d.db.Exec("DELETE FROM entries WHERE metric_id = ? AND date = ?", metricId, date)
	if err != nil {
		return fmt.Errorf("error querying the database to delete the entry: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return fmt.Errorf("error deleting entry, rows affected = %d: %w", count, err)
	}

	return nil
}

func (d *Repository) GetSortedEntriesForMetric(metricId int) (model.Entries, error) {
	rows, err := d.db.Query("SELECT metric_id, value, date FROM entries WHERE metric_id = ? ORDER BY date", metricId)
	if err != nil {
		return nil, fmt.Errorf("error querying entries: %w", err)
	}
	defer rows.Close()

	var entries model.Entries
	for rows.Next() {
		var entry model.Entry
		if err = rows.Scan(&entry.MetricID, &entry.Value, &entry.Date); err != nil {
			return nil, fmt.Errorf("error scanning entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
