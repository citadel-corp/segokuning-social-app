ALTER TABLE user_friends DROP CONSTRAINT IF EXISTS user_friends_user_id_friend_id_unique;
ALTER TABLE user_friends ADD CONSTRAINT user_friends_user_id_friend_id_unique UNIQUE (user_id, friend_id);
