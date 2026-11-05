DO $$
BEGIN
    RAISE EXCEPTION
        'workout snapshot migration cannot be rolled back without discarding per-workout plan history';
END
$$;
