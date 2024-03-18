DROP INDEX IF EXISTS user_friends_user_id;

DROP TABLE IF EXISTS user_friends;

ALTER TABLE user_friends DROP CONSTRAINT fk_user_id;
ALTER TABLE user_friends DROP CONSTRAINT fk_friend_id;
ALTER TABLE user_friends ALTER COLUMN created_at SET DEFAULT;
