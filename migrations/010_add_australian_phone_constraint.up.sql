-- Add constraint for Australian phone numbers
-- Accepts 10-digit numbers starting with 04 (mobile) or 02,03,07,08 (landline)
-- Also accepts 11-digit international format starting with 61
ALTER TABLE users ADD CONSTRAINT check_australian_phone_format CHECK (
    phone IS NULL OR
    phone ~ '^04[0-9]{8}$' OR  -- Mobile: 04XXXXXXXX
    phone ~ '^0[2378][0-9]{8}$' OR  -- Landline: 0XXXXXXXXX
    phone ~ '^61[0-9]{9}$'  -- International: 61XXXXXXXXX
);