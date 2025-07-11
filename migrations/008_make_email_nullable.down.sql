-- Revert email field back to not null
ALTER TABLE users ALTER COLUMN email SET NOT NULL;