package repository

import (
	"database/sql"
	"fmt"

	"github.com/belarte/metrix/model"
	_ "modernc.org/sqlite"
)

const file = ":memory:"

const schema = `
	PRAGMA foreign_keys = ON;

    CREATE TABLE IF NOT EXISTS metrics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        unit TEXT NOT NULL,
        description TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS entries (
        metric_id INTEGER NOT NULL,
        date TEXT NOT NULL,
        value REAL NOT NULL,
		PRIMARY KEY (metric_id, date),
		FOREIGN KEY (metric_id) REFERENCES metrics (id) ON DELETE CASCADE
    );

    INSERT INTO metrics (title, unit, description) VALUES ('Metric 1', 'unit', 'description');
    INSERT INTO metrics (title, unit, description) VALUES ('Metric 2', 'unit', 'description');
    INSERT INTO metrics (title, unit, description) VALUES ('Metric 3', 'unit', 'description');

    INSERT INTO entries (metric_id, value, date) VALUES (1, 5.0, '2018-01-01');
    INSERT INTO entries (metric_id, value, date) VALUES (2, 2.1, '2018-01-11');
    INSERT INTO entries (metric_id, value, date) VALUES (1, 1.0, '2018-01-15');
`

type Repository struct {
	db *sql.DB
}

func New() (*Repository, error) {
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

func (d *Repository) GetMetrics() ([]model.Metric, error) {
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
	res, err := d.db.Exec(
		"INSERT OR REPLACE INTO metrics (id, title, unit, description) VALUES (?, ?, ?, ?)",
		metric.ID, metric.Title, metric.Unit, metric.Description,
	)
	if err != nil {
		return model.Metric{}, fmt.Errorf("error inserting metric: %w", err)
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return model.Metric{}, fmt.Errorf("error getting last insert ID: %w", err)
	}
	metric.ID = int(id)
	return metric, nil
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

func (d *Repository) UpsertEntry(metricId int, value float64, date string) (model.Entry, error) {
	_, err := d.db.Exec(
		"INSERT OR REPLACE INTO entries (metric_id, value, date) VALUES (?, ?, ?)",
		metricId, value, date,
	)

	if err != nil {
		return model.Entry{}, fmt.Errorf("error inserting entry: %w", err)
	}

	return model.Entry{MetricID: metricId, Value: value, Date: date}, nil
}
