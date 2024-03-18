CREATE TABLE IF NOT EXISTS
comments (
    id SERIAL PRIMARY KEY,
    user_id CHAR(16) NOT NULL,
    post_id CHAR(16) NOT NULL,
		content VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE comments DROP CONSTRAINT IF EXISTS fk_user_id;
ALTER TABLE comments
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE comments DROP CONSTRAINT IF EXISTS fk_post_id;
ALTER TABLE comments
	ADD CONSTRAINT fk_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;
