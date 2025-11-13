DROP TABLE IF EXISTS equipment;

CREATE TABLE IF NOT EXISTS equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

DROP TABLE IF EXISTS exercise_equipment;

CREATE TABLE IF NOT EXISTS exercise_equipment (
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
    location TEXT NOT NULL,
    description TEXT DEFAULT 'No description'
    -- figure out how to add working hours
);

DROP TABLE IF EXISTS gym_equipment;

CREATE TABLE gym_equipment (
    gym_id INTEGER,
    equipment INTEGER,
    PRIMARY KEY (gym_id, equipment),
    FOREIGN KEY (gym_id) REFERENCES gym(id),
    FOREIGN KEY (equipment) REFERENCES equipment(id)
);

DROP TABLE IF EXISTS gym_users;

CREATE TABLE gym_users (
    gym_id INTEGER,
    user_id INTEGER,
    PRIMARY KEY (gym_id, user_id),
    FOREIGN KEY (gym_id) REFERENCES gym(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

DROP TABLE IF EXISTS gym_tag;

CREATE TABLE gym_tag (
    name TEXT PRIMARY KEY
);

DROP TABLE IF EXISTS gyms_tags;

CREATE TABLE gyms_tags (
    gym INTEGER,
    tag TEXT,
    PRIMARY KEY (gym, tag),
    FOREIGN KEY (gym) REFERENCES gym(id),
    FOREIGN KEY (tag) REFERENCES gym_tag(name)
);
