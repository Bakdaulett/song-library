DROP INDEX IF EXISTS idx_songs_group_id;
DROP INDEX IF EXISTS idx_songs_song;

DROP TABLE IF EXISTS songs;
DROP TABLE IF EXISTS groups;

DELETE FROM schema_migrations WHERE version = 1;
