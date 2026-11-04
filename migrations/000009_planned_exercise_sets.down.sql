UPDATE workoutapp.training_day_exercises AS tde
SET sets = GREATEST(
    1,
    (
        SELECT count(*)
        FROM workoutapp.training_day_exercise_sets AS planned_set
        WHERE planned_set.training_day_exercise_id = tde.id
    )::INT
);

DROP TABLE workoutapp.training_day_exercise_sets;

ALTER TABLE workoutapp.training_day_exercises
    DROP CONSTRAINT training_day_exercises_sets_check,
    ADD CONSTRAINT training_day_exercises_sets_check
        CHECK (sets BETWEEN 1 AND 100);
