CREATE TABLE transactions (
    id Int64,
    user_id Int64,
    level_id Int64,
    balls Int64,
    step Int64,
    updated_at DATE,
    deleted_at Nullable(Date),
    created_at DATE,
    primary key (id)
) ENGINE = MergeTree
ORDER BY id;