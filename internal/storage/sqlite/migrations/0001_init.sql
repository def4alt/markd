-- +goose Up
CREATE TABLE IF NOT EXISTS bookmarks (
  id TEXT PRIMARY KEY,
  url TEXT NOT NULL,
  title TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'unchecked',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS bookmark_tags (
  bookmark_id TEXT NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (bookmark_id, tag),
  FOREIGN KEY (bookmark_id) REFERENCES bookmarks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS bookmark_tags;
DROP TABLE IF EXISTS bookmarks;
