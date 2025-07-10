-- Inserting initial player data
INSERT INTO players (id, name, password, email, level, created_at, updated_at, match_history, state)
VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'Alice Smith', 'hashed_password_1', 'alice@example.com', 1.5, '2025-07-01T10:00:00Z', '2025-07-02T12:00:00Z', '{"match1": true, "match2": false}', 'Online'),
    ('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'Bob Johnson', 'hashed_password_2', 'bob@example.com', 2.0, '2025-07-01T11:00:00Z', '2025-07-02T13:00:00Z', '{"match3": true}', 'InGame'),
    ('7d793037-a076-4d85-b3a4-6b0b0d6f2b6f', 'Charlie Brown', 'hashed_password_3', 'charlie@example.com', 1.0, '2025-07-01T12:00:00Z', '2025-07-02T14:00:00Z', '{}', 'Offline');