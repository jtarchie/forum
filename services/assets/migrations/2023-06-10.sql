CREATE TABLE IF NOT EXISTS forums (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  parent_id INTEGER DEFAULT 0,
  created_at datetime NOT NULL DEFAULT current_timestamp,
  updated_at datetime NOT NULL DEFAULT current_timestamp
);
CREATE VIEW IF NOT EXISTS ordered_forums AS WITH RECURSIVE paths(id, name, description, parent_id, path) AS (
  SELECT forums.id,
    forums.name,
    forums.description,
    forums.parent_id,
    forums.name
  FROM forums
  WHERE forums.parent_id = 0
  UNION ALL
  SELECT forums.id,
    forums.name,
    forums.description,
    forums.parent_id,
    paths.path || '/' || forums.name
  FROM forums
    JOIN paths ON paths.id = forums.parent_id
)
SELECT id,
  name,
  description,
  parent_id
FROM paths
ORDER BY path;