DROP TABLE IF EXISTS workout_plan;

CREATE TABLE workout_plan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    creator INTEGER NOT NULL,
    description TEXT, --NOTE: add 'date_created' field
    UNIQUE(name, creator),
    FOREIGN KEY(creator) REFERENCES users(id)
);

INSERT INTO workout_plan (name, creator, description) VALUES ('Defualt', 1, '');



DROP TABLE IF EXISTS exercise_day;

CREATE TABLE exercise_day (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plan INTEGER NOT NULL,
    day_name TEXT NOT NULL,
    exercise INTEGER NOT NULL,
    weight FLOAT,
    unit TEXT DEFAULT 'kg' CHECK (unit IN ('kg', 'lbs', 's', 'm', '')),
    sets INTEGER DEFAULT 1, --NOTE: check if >= 1
    min_reps TEXT NOT NULL, --NOTE: this should support a variety of vals like 2, 14, 30s, 2m
    max_reps TEXT, --NOTE: if this is null then the exercise isn't ranged like 6-12 reps but like 5 sets of 5
    day_order INTEGER NOT NULL,
    exercise_order INTEGER NOT NULL,
    FOREIGN KEY(plan) REFERENCES workout_plan(id),
    FOREIGN KEY(exercise) REFERENCES exercises(id)
);



DROP TABLE IF EXISTS users_plans;

CREATE TABLE users_plans (
    usr INTEGER NOT NULL,
    plan INTEGER NOT NULL, --ADD LIKE A DATE_ADDED OR SOMETHING THAT WOULD ALLOW SORTING
    PRIMARY KEY (usr, plan),
    FOREIGN KEY(usr) REFERENCES users(id),
    FOREIGN KEY(plan) REFERENCES workout_plan(id)
);

--------------------------------------------------------------------------------

--------------------------------------------------------------------------------

DROP TABLE IF EXISTS workout_track_data;

CREATE TABLE workout_track_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    track INTEGER NOT NULL,
    ex_day INTEGER NOT NULL,
    weigth FLOAT, -- NOTE: This doesn't neccessarily have to be the weight that is specified in the ex_day, unless it's null
    set_num INTEGER,
    rep_num INTEGER,
    FOREIGN KEY(ex_day) REFERENCES exercise_day(id),
    FOREIGN KEY(track) REFERENCES workout_track(id)
);

DROP TABLE IF EXISTS workout_track;

CREATE TABLE workout_track (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plan INTEGER NOT NULL,
    usr INTEGER NOT NULL,
    workout_date DATE,
    -- NOTE: Consider adding a is_private field, which controls whether other users can see your progress
    FOREIGN KEY(plan) REFERENCES workout_plan(id),
    FOREIGN KEY(usr) REFERENCES users(id),
    UNIQUE(plan, usr, workout_date)
);

