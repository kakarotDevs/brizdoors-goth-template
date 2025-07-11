-- Add phone number column to users table
ALTER TABLE users ADD COLUMN phone VARCHAR(20) UNIQUE;

-- Add index for phone lookups
CREATE INDEX idx_users_phone ON users(phone);

-- Add constraint to ensure phone numbers are properly formatted (optional)
-- This ensures phone numbers are at least 10 digits and can include spaces, dashes, and parentheses
ALTER TABLE users ADD CONSTRAINT check_phone_format CHECK (phone IS NULL OR phone ~ '^[\+]?[0-9\s\-\(\)]{10,}$');