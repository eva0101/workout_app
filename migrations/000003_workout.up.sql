CREATE TABLE workoutapp.workout (
    id SERIAL PRIMARY KEY,
    training_day_id INT NOT NULL REFERENCES workoutapp.training_days(id),

    status VARCHAR(20) NOT NULL DEFAULT 'planned'
        CHECK(status IN ('planned', 'in_progress', 'completed')),

    fatigue_score INT CHECK(fatigue_score BETWEEN 1 AND 10),

    begin_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    CHECK (
        completed_at IS NULL
        OR begin_at IS NULL
        OR completed_at >= begin_at
    )
);

CREATE TABLE workoutapp.workout_exercise (
    id SERIAL PRIMARY KEY,
    workout_id INT NOT NULL REFERENCES workoutapp.workout(id),
    exercise_id INT NOT NULL REFERENCES workoutapp.exercise(id),
    
    status VARCHAR(20) NOT NULL DEFAULT 'planned'
        CHECK(status IN ('planned', 'in_progress', 'completed')),

    sets_done INT CHECK(sets_done >= 0),
    reps_done INT CHECK(reps_done >= 0),
    duration_sec INT CHECK(duration_sec >= 0)
);
