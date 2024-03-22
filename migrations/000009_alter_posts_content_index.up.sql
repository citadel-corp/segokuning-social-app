CREATE INDEX IF NOT EXISTS posts_content
	ON posts USING HASH (content);
