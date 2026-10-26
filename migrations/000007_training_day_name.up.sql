ALTER TABLE workoutapp.training_days
    ADD COLUMN name VARCHAR(100);

UPDATE workoutapp.training_days
SET name = 'День ' || day_number::TEXT;

ALTER TABLE workoutapp.training_days
    ALTER COLUMN name SET NOT NULL,
    ADD CONSTRAINT training_days_name_length_check
        CHECK (char_length(name) BETWEEN 3 AND 100);
