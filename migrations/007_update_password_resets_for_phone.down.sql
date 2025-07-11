-- Revert password_resets table changes
ALTER TABLE password_resets DROP CONSTRAINT IF EXISTS check_contact_method;
ALTER TABLE password_resets DROP CONSTRAINT IF EXISTS check_method;
DROP INDEX IF EXISTS idx_password_resets_method;
DROP INDEX IF EXISTS idx_password_resets_phone;
ALTER TABLE password_resets DROP COLUMN IF EXISTS method;
ALTER TABLE password_resets DROP COLUMN IF EXISTS phone;
ALTER TABLE password_resets ALTER COLUMN email SET NOT NULL;