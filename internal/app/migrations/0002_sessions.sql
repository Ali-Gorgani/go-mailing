-- Migration Up
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    token VARCHAR(511) NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Migration Down
DROP TABLE sessions;