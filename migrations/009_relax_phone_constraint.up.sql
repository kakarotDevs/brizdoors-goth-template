-- Remove the phone format constraint as it might be too restrictive
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_phone_format;