CREATE TABLE IF NOT EXISTS commit_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id TEXT NOT NULL,
	repo_path TEXT,
	commit TEXT NOT NULL,
	branch TEXT,
	timestamp TEXT NOT NULL,
	status INTEGER NOT NULL,
	UNIQUE(repo_id, commit)
);