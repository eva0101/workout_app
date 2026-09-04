CREATE TABLE workoutapp.workout (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES workoutapp.users(id),
    training_day_id INT NOT NULL REFERENCES workoutapp.training_days(id),

    status VARCHAR(20) NOT NULL DEFAULT 'planned'
        CHECK(status IN ('planned', 'in_progress', 'completed')),

    fatigue_score INT CHECK(fatigue_score BETWEEN 0 AND 10),

    begin_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    total_time INTERVAL,

    CONSTRAINT unique_training_day UNIQUE (training_day_id),

    CHECK (
        completed_at IS NULL
        OR begin_at IS NULL
        OR completed_at >= begin_at
    )
);

CREATE UNIQUE INDEX one_active_workout_per_user
ON workoutapp.workout (user_id)
WHERE status = 'in_progress';

CREATE TABLE workoutapp.workout_set (
    id SERIAL PRIMARY KEY,
    workout_id INT NOT NULL REFERENCES workoutapp.workout(id),
    training_day_exercises_id INT NOT NULL REFERENCES workoutapp.training_day_exercises(id),

    set_number INT NOT NULL,
    reps_done INT CHECK(reps_done >= 0),
    weight NUMERIC(5,2) CHECK(weight >= 0),

    UNIQUE(workout_id, training_day_exercises_id, set_number)
);
