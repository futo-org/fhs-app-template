-- +goose Up
-- create "albums" table
CREATE TABLE `albums` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` varchar NOT NULL
);
-- create "album_media" table
CREATE TABLE `album_media` (
  `album_id` integer NOT NULL,
  `media_id` integer NOT NULL,
  CONSTRAINT `0` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`album_id`) REFERENCES `albums` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- +goose Down
-- reverse: create "album_media" table
DROP TABLE `album_media`;
-- reverse: create "albums" table
DROP TABLE `albums`;
