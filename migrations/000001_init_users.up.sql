CREATE TABLE IF NOT EXISTS
users (
    id CHAR(16) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
	email VARCHAR NOT NULL UNIQUE,
	phone_number VARCHAR(15) NOT NULL UNIQUE,
    hashed_password BYTEA NOT NULL,
	friend_count INT NOT NULL DEFAULT 0,
	image_url VARCHAR NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);
