-- Remove phone number column from users table
ALTER TABLE users DROP COLUMN IF EXISTS phone;