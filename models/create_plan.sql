DROP TABLE IF EXISTS exercise_day;

CREATE TABLE exercise_day (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    day_name TEXT UNIQUE NOT NULL,
    exercise INTEGER NOT NULL,
    sets INTEGER NOT NULL,
    min_reps TEXT NOT NULL, --NOTE: this should support a variety of vals like 2, 14, 30s, 2m
    max_reps TEXT, --NOTE: if this is null then the exercise isn't ranged like 6-12 reps 
    -- but like 5 sets of 5
    FOREIGN KEY(exercise) REFERENCES exercises(id)
); -- NOTE: this seems waaay to inefficient, reconsider how it's done

DROP TABLE IF EXISTS plan_template;

CREATE TABLE plan_template (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- creator INTEGER NOT NULL, -- shouldn't be here, but in another table
    exercise_day INTEGER UNIQUE NOT NULL,

    FOREIGN KEY(exercise_day) REFERENCES exercise_day(id)
)
