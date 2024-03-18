CREATE TABLE IF NOT EXISTS
user_friends (
    id SERIAL PRIMARY KEY,
		user_id CHAR(16) NOT NULL,
    friend_id CHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE user_friends DROP CONSTRAINT IF EXISTS fk_user_id;
ALTER TABLE user_friends
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_friends DROP CONSTRAINT IF EXISTS fk_friend_id;
ALTER TABLE user_friends
	ADD CONSTRAINT fk_friend_id FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS user_friends_user_id
	ON user_friends USING HASH (user_id);
