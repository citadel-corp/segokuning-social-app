DROP INDEX IF EXISTS users_username;

DROP TABLE IF EXISTS users;

ALTER TABLE users DROP CONSTRAINT users_username_unique UNIQUE (username);
ALTER TABLE users ALTER COLUMN created_at SET DEFAULT;
