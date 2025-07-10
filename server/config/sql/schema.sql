-- Creating the players table
CREATE TABLE IF NOT EXISTS players (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    level DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    match_history JSONB NOT NULL DEFAULT '{}',
    state VARCHAR(20) NOT NULL DEFAULT 'Offline',
    CONSTRAINT valid_name CHECK (length(name) >= 3 AND length(name) <= 50),
    CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT valid_level CHECK (level >= 0),
    CONSTRAINT valid_state CHECK (state IN ('InGame', 'InQuery', 'Offline', 'Online'))
);