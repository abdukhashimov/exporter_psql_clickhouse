CREATE TABLE transactions (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id INTEGER,
    balls BIGINT,
    level_id INTEGER,
    step INTEGER,
    deleted_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_at TIMESTAMP
);