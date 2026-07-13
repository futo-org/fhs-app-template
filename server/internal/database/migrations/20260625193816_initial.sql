-- +goose Up
-- create "media" table
CREATE TABLE `media` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `title` varchar NULL,
  `url` varchar NULL,
  `path` varchar NULL
);

-- +goose Down
-- reverse: create "media" table
DROP TABLE `media`;
