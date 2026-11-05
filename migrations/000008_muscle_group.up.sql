CREATE TABLE workoutapp.muscle_group (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    image_path TEXT
);

INSERT INTO workoutapp.muscle_group (name, title, image_path) VALUES
    ('chest', 'Грудь', '/exercises/groups/chest.webp'),
    ('back', 'Спина', '/exercises/groups/back.webp'),
    ('legs', 'Ноги', '/exercises/groups/legs.webp'),
    ('biceps', 'Бицепс', '/exercises/groups/biceps.webp'),
    ('triceps', 'Трицепс', '/exercises/groups/triceps.webp'),
    ('shoulders', 'Плечи', '/exercises/groups/shoulders.webp'),
    ('abs', 'Пресс', '/exercises/groups/abs.webp');

ALTER TABLE workoutapp.exercise
    ADD CONSTRAINT exercise_muscle_group_id_fkey
    FOREIGN KEY (muscle_group_id)
    REFERENCES workoutapp.muscle_group(id);
