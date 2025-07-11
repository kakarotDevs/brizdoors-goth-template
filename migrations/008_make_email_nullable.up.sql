-- Make email field nullable to support phone-only accounts
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;