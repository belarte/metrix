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
