ALTER TABLE workoutapp.workout_exercise
    ADD COLUMN muscle_group_name VARCHAR(50);

UPDATE workoutapp.workout_exercise we
SET muscle_group_name = mg.name
FROM workoutapp.exercise e
JOIN workoutapp.muscle_group mg ON mg.id = e.muscle_group_id
WHERE e.id = we.exercise_id;

ALTER TABLE workoutapp.workout_exercise
    ALTER COLUMN muscle_group_name SET NOT NULL,
    ADD CONSTRAINT workout_exercise_muscle_group_name_check
        CHECK (muscle_group_name IN (
            'chest', 'back', 'legs', 'biceps', 'triceps', 'shoulders', 'abs'
        ));
