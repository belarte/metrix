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
