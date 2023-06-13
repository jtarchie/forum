CREATE TABLE IF NOT EXISTS migrations (
  id INTEGER NOT NULL PRIMARY KEY,
  version TEXT NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS migrations_versions ON migrations(version);