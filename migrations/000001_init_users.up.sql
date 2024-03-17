CREATE TABLE IF NOT EXISTS
    users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(15) NOT NULL,
        name VARCHAR(50) NOT NULL,
        hashed_password BYTEA NOT NULL,
        created_at TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS users_username
	ON users USING HASH (username);

ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);
ALTER TABLE users ALTER COLUMN created_at SET DEFAULT current_timestamp;
