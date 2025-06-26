DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    training_since DATETIME,
    is_trainer BOOLEAN NOT NULL DEFAULT 0,
    gym_goals TEXT CHECK (gym_goals IN ('health', 'strength', 'mobility', 'athleticism', 'bodybuilding', '')),
    current_gym TEXT
);
