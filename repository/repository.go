package repository

import (
	"database/sql"
	"fmt"

	"github.com/belarte/metrix/model"
	_ "modernc.org/sqlite"
)

const file = ":memory:"

const schema = `
    CREATE TABLE IF NOT EXISTS metrics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        unit TEXT NOT NULL,
        description TEXT NOT NULL
    );

    INSERT INTO metrics (title, unit, description) VALUES ('Metric 1', 'unit', 'description');
    INSERT INTO metrics (title, unit, description) VALUES ('Metric 2', 'unit', 'description');
    INSERT INTO metrics (title, unit, description) VALUES ('Metric 3', 'unit', 'description');
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
