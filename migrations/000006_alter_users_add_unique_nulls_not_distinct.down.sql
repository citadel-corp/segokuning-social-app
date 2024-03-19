ALTER TABLE users DROP CONSTRAINT users_unique_email UNIQUE NULLS NOT DISTINCT (email);
ALTER TABLE users DROP CONSTRAINT users_unique_phone_number UNIQUE NULLS NOT DISTINCT (phone_number);
