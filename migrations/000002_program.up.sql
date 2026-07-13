CREATE TABLE workoutapp.program (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES workoutapp.users(id),
    name VARCHAR(100) NOT NULL CHECK(char_length(name) BETWEEN 3 AND 100),
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ
);

CREATE TABLE workoutapp.training_days (
    id SERIAL PRIMARY KEY,
    program_id INT NOT NULL REFERENCES workoutapp.program(id),
    day_number INT NOT NULL CHECK(day_number BETWEEN 1 AND 100),
    UNIQUE (program_id, day_number)
);

CREATE TABLE workoutapp.exercise (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL CHECK(char_length(name) BETWEEN 3 AND 100),
    muscle_group VARCHAR(100) NOT NULL CHECK(char_length(muscle_group) BETWEEN 3 AND 100),
    description TEXT
);

CREATE TABLE workoutapp.training_day_exercises (
    id SERIAL PRIMARY KEY,
    training_day_id INT NOT NULL REFERENCES workoutapp.training_days(id),
    exercise_id INT NOT NULL REFERENCES workoutapp.exercise(id),
    sets INT CHECK(sets BETWEEN 1 AND 100),
    reps INT CHECK(reps BETWEEN 1 AND 1000),
    duration_sec INT CHECK(duration_sec BETWEEN 1 AND 1000)
);
