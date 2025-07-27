DROP TABLE IF EXISTS workout_track_data;

CREATE TABLE workout_track_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- NOTE: NOT NECCESARY, CAN USE TRACK AND EX_DAY AND SET_NUM AS PRIMARY KEY
    track INTEGER NOT NULL,
    ex_day INTEGER NOT NULL,
    weight FLOAT, -- NOTE: This doesn't neccessarily have to be the weight that is specified in the ex_day, unless it's null
    set_num INTEGER,
    rep_num INTEGER DEFAULT 0,
    FOREIGN KEY(ex_day) REFERENCES exercise_day(id),
    FOREIGN KEY(track) REFERENCES workout_track(id)
);



DROP TABLE IF EXISTS workout_track;

CREATE TABLE workout_track (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plan INTEGER NOT NULL,
    usr INTEGER NOT NULL,
    workout_date DATE,
    is_private BOOLEAN DEFAULT TRUE,
    -- NOTE: Consider adding a is_private field, which controls whether other users can see your progress
    FOREIGN KEY(plan) REFERENCES workout_plan(id),
    FOREIGN KEY(usr) REFERENCES users(id),
    UNIQUE(plan, usr, workout_date)
);



DROP TABLE IF EXISTS track_exercise;

-- NOTE: this table is used to track which ex_days are for which track since ex days can be changed in a plan, but the track needs to remember them
-- CREATE TABLE track_exercise (
--     track INTEGER,
--     ex_day INTEGER,
--     PRIMARY KEY(track, ex_day),
--     FOREIGN KEY(track) REFERENCES workout_track(id),
--     FOREIGN KEY(ex_day) REFERENCES exercise_day(id)
-- );
