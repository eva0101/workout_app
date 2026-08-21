CREATE TABLE workoutapp.statistics (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES workoutapp.users(id),
    exercise_id INT NOT NULL REFERENCES workoutapp.exercise(id),

    total_sets INT DEFAULT 0,
    total_repeats INT DEFAULT 0,
    maximum_repeats INT DEFAULT 0,
    total_volume NUMERIC(12, 2) DEFAULT 0,
    maximum_weight NUMERIC(5, 2) DEFAULT 0,

    UNIQUE (user_id, exercise_id)
);
