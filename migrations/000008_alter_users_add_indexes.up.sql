CREATE INDEX IF NOT EXISTS users_id_hash
	ON users USING HASH (id);

CREATE INDEX IF NOT EXISTS users_friend_count
	ON users(friend_count);
