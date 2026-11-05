ALTER TABLE workoutapp.workout_set
    DROP CONSTRAINT workout_set_workout_id_workout_exercise_id_fkey;

ALTER TABLE workoutapp.workout_exercise
    DROP CONSTRAINT workout_exercise_workout_id_id_key;
