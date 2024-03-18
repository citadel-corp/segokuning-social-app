CREATE TABLE IF NOT EXISTS
posts (
  id VARCHAR PRIMARY KEY,
  user_id INT NOT NULL,
  content VARCHAR(500) NOT NULL,
  tags TEXT[] NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE posts DROP CONSTRAINT IF EXISTS fk_user_id;
ALTER TABLE posts
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS posts_tags
	ON posts USING gin(tags);
