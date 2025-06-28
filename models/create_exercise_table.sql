DROP TABLE IF EXISTS exercises;

CREATE TABLE exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    exercise_type TEXT NOT NULL CHECK (exercise_type IN ('hypertrophy', 'strength', 'mobility', 'stability', 'athleticism', 'cardio', 'recovery', 'functional', '')),
    difficulty INTEGER CHECK (difficulty IN (1, 2, 3))
);

DROP TABLE IF EXISTS exercise_targets;

CREATE TABLE exercise_targets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise INTEGER NOT NULL,
    target TEXT NOT NULL,
    FOREIGN KEY(exercise) REFERENCES exercises(id)
);

DROP TABLE IF EXISTS gym_equipment;

CREATE TABLE gym_equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

DROP TABLE IF EXISTS exercise_equipment;

CREATE TABLE exercise_equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    equipment INTEGER NOT NULL,
    exercise INTEGER NOT NULL,
    FOREIGN KEY(exercise) REFERENCES exercises(id),
    FOREIGN KEY(equipment) REFERENCES gym_equipment(id)
);
