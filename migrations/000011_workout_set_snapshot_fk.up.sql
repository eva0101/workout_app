ALTER TABLE workoutapp.workout_exercise
    ADD CONSTRAINT workout_exercise_workout_id_id_key
        UNIQUE (workout_id, id);

ALTER TABLE workoutapp.workout_set
    ADD CONSTRAINT workout_set_workout_id_workout_exercise_id_fkey
        FOREIGN KEY (workout_id, workout_exercise_id)
        REFERENCES workoutapp.workout_exercise(workout_id, id)
        ON DELETE CASCADE;
