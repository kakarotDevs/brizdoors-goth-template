-- Update password_resets table to support both email and phone
ALTER TABLE password_resets ADD COLUMN phone VARCHAR(20);
ALTER TABLE password_resets ADD COLUMN method VARCHAR(10) NOT NULL DEFAULT 'email';

-- Make email nullable since we now support phone
ALTER TABLE password_resets ALTER COLUMN email DROP NOT NULL;

-- Add index for phone lookups
CREATE INDEX idx_password_resets_phone ON password_resets(phone);

-- Add index for method lookups
CREATE INDEX idx_password_resets_method ON password_resets(method);

-- Add constraint to ensure method is valid
ALTER TABLE password_resets ADD CONSTRAINT check_method CHECK (method IN ('email', 'sms'));

-- Add constraint to ensure at least one contact method is provided
ALTER TABLE password_resets ADD CONSTRAINT check_contact_method CHECK (email IS NOT NULL OR phone IS NOT NULL);