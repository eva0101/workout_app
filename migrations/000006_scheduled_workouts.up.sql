CREATE TABLE workoutapp.scheduled_workouts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES workoutapp.users(id),
    training_day_id INTEGER NOT NULL
        REFERENCES workoutapp.training_days(id)
        ON DELETE CASCADE,
    scheduled_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, training_day_id, scheduled_date)
);

CREATE INDEX scheduled_workouts_user_date_idx
    ON workoutapp.scheduled_workouts (user_id, scheduled_date);
