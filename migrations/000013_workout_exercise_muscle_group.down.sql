ALTER TABLE workoutapp.workout_exercise
    DROP CONSTRAINT workout_exercise_muscle_group_name_check,
    DROP COLUMN muscle_group_name;
