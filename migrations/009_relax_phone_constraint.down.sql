-- Restore the phone format constraint
ALTER TABLE users ADD CONSTRAINT check_phone_format CHECK (phone IS NULL OR phone ~ '^[\+]?[0-9\s\-\(\)]{10,}$');