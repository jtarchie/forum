CREATE TABLE IF NOT EXISTS posts (
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  value TEXT NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp,
  updated_at datetime NOT NULL DEFAULT current_timestamp
);