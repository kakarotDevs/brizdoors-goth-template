-- Remove Australian phone format constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_australian_phone_format;