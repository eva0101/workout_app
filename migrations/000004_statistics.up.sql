CREATE TABLE workoutapp.statistics (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES workoutapp.users(id),
    exercise_id INT NOT NULL,

    total_num_sets INT NOT NULL DEFAULT 0,
    total_num_repeats INT NOT NULL DEFAULT 0,
    maximum_repeats INT NOT NULL DEFAULT 0,

    UNIQUE (user_id, exercise_id)
);
