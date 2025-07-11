CREATE TABLE password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster token lookups
CREATE INDEX idx_password_resets_token ON password_resets(token);

-- Index for cleanup queries
CREATE INDEX idx_password_resets_expires_at ON password_resets(expires_at);

-- Index for user lookups
CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);