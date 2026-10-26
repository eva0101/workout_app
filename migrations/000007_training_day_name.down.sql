ALTER TABLE workoutapp.training_days
    DROP CONSTRAINT training_days_name_length_check,
    DROP COLUMN name;
