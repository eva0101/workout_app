ALTER TABLE workoutapp.workout
    DROP CONSTRAINT unique_training_day,
    DROP CONSTRAINT workout_training_day_id_fkey;

CREATE TABLE workoutapp.workout_exercise (
    id SERIAL PRIMARY KEY,
    workout_id INT NOT NULL
        REFERENCES workoutapp.workout(id)
        ON DELETE CASCADE,
    source_training_day_exercise_id INT NOT NULL,
    exercise_id INT NOT NULL,
    exercise_name VARCHAR(100) NOT NULL,
    exercise_title TEXT NOT NULL,
    sort_order INT NOT NULL CHECK (sort_order >= 1),
    UNIQUE (workout_id, source_training_day_exercise_id)
);

CREATE TABLE workoutapp.workout_exercise_sets (
    id SERIAL PRIMARY KEY,
    workout_exercise_id INT NOT NULL
        REFERENCES workoutapp.workout_exercise(id)
        ON DELETE CASCADE,
    set_number INT NOT NULL CHECK (set_number BETWEEN 1 AND 100),
    repetitions INT CHECK (repetitions >= 0),
    weight NUMERIC(5, 2) CHECK (weight >= 0),
    UNIQUE (workout_exercise_id, set_number)
);

WITH workout_source_exercises AS (
    SELECT w.id AS workout_id, tde.id AS source_training_day_exercise_id
    FROM workoutapp.workout w
    JOIN workoutapp.training_day_exercises tde
        ON tde.training_day_id = w.training_day_id

    UNION

    SELECT ws.workout_id, ws.training_day_exercises_id
    FROM workoutapp.workout_set ws
),
ranked_workout_exercises AS (
    SELECT
        wse.workout_id,
        tde.id AS source_training_day_exercise_id,
        tde.exercise_id,
        e.name AS exercise_name,
        e.title AS exercise_title,
        ROW_NUMBER() OVER (
            PARTITION BY wse.workout_id
            ORDER BY tde.id
        )::INT AS sort_order
    FROM workout_source_exercises wse
    JOIN workoutapp.training_day_exercises tde
        ON tde.id = wse.source_training_day_exercise_id
    JOIN workoutapp.exercise e
        ON e.id = tde.exercise_id
)
INSERT INTO workoutapp.workout_exercise (
    workout_id,
    source_training_day_exercise_id,
    exercise_id,
    exercise_name,
    exercise_title,
    sort_order
)
SELECT
    workout_id,
    source_training_day_exercise_id,
    exercise_id,
    exercise_name,
    exercise_title,
    sort_order
FROM ranked_workout_exercises;

INSERT INTO workoutapp.workout_exercise_sets (
    workout_exercise_id,
    set_number,
    repetitions,
    weight
)
SELECT
    we.id,
    tes.set_number,
    tes.repetitions,
    tes.weight
FROM workoutapp.workout_exercise we
JOIN workoutapp.training_day_exercise_sets tes
    ON tes.training_day_exercise_id = we.source_training_day_exercise_id;

ALTER TABLE workoutapp.workout_set
    ADD COLUMN workout_exercise_id INT;

UPDATE workoutapp.workout_set ws
SET workout_exercise_id = we.id
FROM workoutapp.workout_exercise we
WHERE we.workout_id = ws.workout_id
    AND we.source_training_day_exercise_id = ws.training_day_exercises_id;

ALTER TABLE workoutapp.workout_set
    DROP CONSTRAINT workout_set_training_day_exercises_id_fkey,
    DROP CONSTRAINT workout_set_workout_id_training_day_exercises_id_set_number_key,
    DROP COLUMN training_day_exercises_id,
    ALTER COLUMN workout_exercise_id SET NOT NULL,
    ADD CONSTRAINT workout_set_workout_exercise_id_fkey
        FOREIGN KEY (workout_exercise_id)
        REFERENCES workoutapp.workout_exercise(id)
        ON DELETE CASCADE,
    ADD CONSTRAINT workout_set_workout_exercise_id_set_number_key
        UNIQUE (workout_exercise_id, set_number);
