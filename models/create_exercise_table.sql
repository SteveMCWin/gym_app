DROP TABLE IF EXISTS exercises;

CREATE TABLE exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    exercise_type TEXT NOT NULL CHECK (exercise_type IN ('hypertrophy', 'strength', 'mobility', 'stability', 'athleticism', 'cardio', 'recovery', 'functional', '')),
    difficulty INTEGER CHECK (difficulty IN (1, 2, 3))
);

DROP TABLE IF EXISTS target;

CREATE TABLE target (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    standard_name TEXT,
    latin_name TEXT,
    CHECK( standard_name IS NOT NULL OR latin_name IS NOT NULL)
)

DROP TABLE IF EXISTS exercise_targets;

CREATE TABLE exercise_targets (
    exercise INTEGER NOT NULL,
    target INTEGER NOT NULL,
    PRIMARY KEY (exercise, target),
    FOREIGN KEY(exercise) REFERENCES exercises(id),
    FOREIGN KEY(target) REFERENCES target(id)
);
