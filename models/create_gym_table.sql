-- DROP TABLE IF EXISTS equipment;

CREATE TABLE IF NOT EXISTS equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- DROP TABLE IF EXISTS exercise_equipment;

CREATE TABLE IF NOT EXISTS exercise_equipment (
    -- id INTEGER PRIMARY KEY AUTOINCREMENT,
    equipment INTEGER NOT NULL,
    exercise INTEGER NOT NULL,
    PRIMARY KEY(equipment, exercise),
    FOREIGN KEY(exercise) REFERENCES exercises(id),
    FOREIGN KEY(equipment) REFERENCES equipment(id)
);

DROP TABLE IF EXISTS gym;

CREATE TABLE gym (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    location TEXT NOT NULL
);

DROP TABLE IF EXISTS gym_equipment;

CREATE TABLE gym_equipment (
    gym_id INTEGER,
    gym_equipment INTEGER,
    PRIMARY KEY (gym_id, gym_equipment),
    FOREIGN KEY (gym_id) REFERENCES gym(id),
    FOREIGN KEY (gym_equipment) REFERENCES equipment(id)
);

DROP TABLE IF EXISTS gym_users;

CREATE TABLE gym_users (
    -- id INTEGER PRIMARY KEY AUTOINCREMENT,
    gym_id INTEGER,
    user_id INTEGER,
    PRIMARY KEY (gym_id, user_id),
    FOREIGN KEY (gym_id) REFERENCES gym(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
