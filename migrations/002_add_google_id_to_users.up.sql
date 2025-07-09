ALTER TABLE users ADD COLUMN google_id TEXT UNIQUE;
-- Optionally: update any existing Google accounts here, set google_id = id for those with numeric id, then update id to a UUID.
-- Example (for simple small DB): you may need some scripting for this; consult with a DBA for production DB!
