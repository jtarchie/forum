CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY,
  email TEXT NOT NULL,
  username TEXT NOT NULL,
  provider TEXT CHECK(provider IN ('github')) NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp,
  updated_at datetime NOT NULL DEFAULT current_timestamp,
  UNIQUE(email, provider),
  UNIQUE(username)
);
