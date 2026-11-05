ALTER TABLE workoutapp.training_day_exercises
    DROP CONSTRAINT training_day_exercises_sets_check,
    ADD CONSTRAINT training_day_exercises_sets_check
        CHECK (sets BETWEEN 0 AND 100);

CREATE TABLE workoutapp.training_day_exercise_sets (
    id SERIAL PRIMARY KEY,
    training_day_exercise_id INT NOT NULL
        REFERENCES workoutapp.training_day_exercises(id)
        ON DELETE CASCADE,
    set_number INT NOT NULL CHECK (set_number BETWEEN 1 AND 100),
    repetitions INT CHECK (repetitions >= 0),
    weight NUMERIC(5, 2) CHECK (weight >= 0),
    UNIQUE (training_day_exercise_id, set_number)
);

INSERT INTO workoutapp.training_day_exercise_sets (
    training_day_exercise_id,
    set_number
)
SELECT tde.id, generated.set_number
FROM workoutapp.training_day_exercises AS tde
CROSS JOIN LATERAL generate_series(1, tde.sets) AS generated(set_number);
