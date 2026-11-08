ALTER TABLE workoutapp.workout
    ADD COLUMN training_day_name VARCHAR(100);

UPDATE workoutapp.workout w
SET training_day_name = td.name
FROM workoutapp.training_days td
WHERE td.id = w.training_day_id;

ALTER TABLE workoutapp.workout_exercise
    ADD COLUMN completed_at TIMESTAMPTZ;
