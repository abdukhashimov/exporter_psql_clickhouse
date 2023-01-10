INSERT INTO transactions(
        id,
        user_id,
        balls,
        level_id,
        step,
        updated_at,
        created_at
    )
SELECT i,
    i,
    300,
    1,
    1,
    now(),
    now()
FROM generate_series(1, 1000355) s(i);