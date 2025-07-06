DROP TABLE IF EXISTS workout_plan;

CREATE TABLE workout_plan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    creator INTEGER NOT NULL,
    description TEXT, --NOTE: add 'date_created' field and perhaps a 'last_tracked' field
    UNIQUE(name, creator),
    FOREIGN KEY(creator) REFERENCES users(id)
);

INSERT INTO workout_plan (name, creator, description) VALUES ('/', 1, '');



DROP TABLE IF EXISTS exercise_day;

CREATE TABLE exercise_day (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plan INTEGER NOT NULL,
    day_name TEXT NOT NULL,
    exercise INTEGER NOT NULL,
    weight FLOAT,
    unit TEXT DEFAULT 'kg' CHECK (unit IN ('kg', 'lbs', 's', 'm', '')),
    --NOTE: add a is_dropset field that will allow the user to add more fields than specified by the sets field (and do so when the sets value isn't provided) ((also set default for sets to 0 or smth))
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
    date_added DATE NOT NULL, --NOTE: probably will remove this since I think using other's plans will just actually be copying them and leaving the creator the same
    PRIMARY KEY (usr, plan),
    FOREIGN KEY(usr) REFERENCES users(id),
    FOREIGN KEY(plan) REFERENCES workout_plan(id)
);

