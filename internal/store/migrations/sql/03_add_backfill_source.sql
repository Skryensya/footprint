-- Add backfill source for importing historical commits
INSERT OR IGNORE INTO event_source (id, name) VALUES (6, 'backfill');
