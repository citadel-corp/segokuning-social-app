ALTER TABLE users
ADD CONSTRAINT users_unique_email UNIQUE NULLS DISTINCT (email);

ALTER TABLE users
ADD CONSTRAINT users_unique_phone_number UNIQUE NULLS DISTINCT (phone_number);