---------------------
--- WOKROUT PLANS ---
---------------------

DELETE from workout_plan;
DELETE FROM sqlite_sequence WHERE name='workout_plan';

INSERT INTO workout_plan (name, creator, description) VALUES ('/', 1, '');

-----------------
--- EXERCISES ---
-----------------

DELETE from exercises;
DELETE FROM sqlite_sequence WHERE name='exercises';

-- CHEST
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('dips', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('barbellbench press', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('incline bench press', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('cable flies', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('pushups', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('machine flies', '', '', 2);
-- BACK
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('barbell bent-over row', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('dumbbell row', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('cable row', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('pullups', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('t-bar row', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('lat pulldowns', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('lat pullovers', '', '', 2);
-- UPPER ARMS
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('skull crushers', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('triceps extensions', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('dumbbell kickback', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('cable pushdowns', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('dumbbell curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('barbell curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('cable curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('spider curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('incline curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('preacher curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('focus curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('hammer curls', '', '', 2);
-- DETLS
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('lateral raises', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('military press', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('shoulder press', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('pike pushups', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('front raises', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('upright row', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('rear delt flies', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('rear delt row', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('lu raises', '', '', 2);
-- QUADS
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('lunges', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('bulgarian split squat', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('barbell squat', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('goblet squat', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('leg extensions', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('sissy squats', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('pistol squats', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('hack sqaut', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('leg press', '', '', 2);
-- HAMSTRINGS
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('hamstring curls', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('romanian deadlift', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('good mornings', '', '', 2);
INSERT into exercises (name, description, exercise_type, difficulty) VALUES ('hamstring raise', '', '', 2);


---------------
--- TARGETS ---
---------------

DELETE from target;
DELETE FROM sqlite_sequence WHERE name='target';

INSERT into target (standard_name, latin_name) VALUES ('chest', '');
INSERT into target (standard_name, latin_name) VALUES ('back', '');
INSERT into target (standard_name, latin_name) VALUES ('biceps', '');
INSERT into target (standard_name, latin_name) VALUES ('triceps', '');
INSERT into target (standard_name, latin_name) VALUES ('shoulders', '');
INSERT into target (standard_name, latin_name) VALUES ('quads', '');
INSERT into target (standard_name, latin_name) VALUES ('hamstrings', '');

------------------------
--- EXERCISE TARGETS ---
------------------------

DELETE from exercise_targets;
DELETE FROM sqlite_sequence WHERE name='exercise_targets';

INSERT into exercise_targets (exercise, target) VALUES (1, 1);
INSERT into exercise_targets (exercise, target) VALUES (2, 1);
INSERT into exercise_targets (exercise, target) VALUES (3, 1);
INSERT into exercise_targets (exercise, target) VALUES (4, 1);
INSERT into exercise_targets (exercise, target) VALUES (5, 1);
INSERT into exercise_targets (exercise, target) VALUES (6, 1);
INSERT into exercise_targets (exercise, target) VALUES (7, 2);
INSERT into exercise_targets (exercise, target) VALUES (8, 2);
INSERT into exercise_targets (exercise, target) VALUES (9, 2);
INSERT into exercise_targets (exercise, target) VALUES (10, 2);
INSERT into exercise_targets (exercise, target) VALUES (11, 2);
INSERT into exercise_targets (exercise, target) VALUES (12, 2);
INSERT into exercise_targets (exercise, target) VALUES (13, 2);
INSERT into exercise_targets (exercise, target) VALUES (14, 4);
INSERT into exercise_targets (exercise, target) VALUES (15, 4);
INSERT into exercise_targets (exercise, target) VALUES (16, 4);
INSERT into exercise_targets (exercise, target) VALUES (17, 4);
INSERT into exercise_targets (exercise, target) VALUES (18, 3);
INSERT into exercise_targets (exercise, target) VALUES (19, 3);
INSERT into exercise_targets (exercise, target) VALUES (20, 3);
INSERT into exercise_targets (exercise, target) VALUES (21, 3);
INSERT into exercise_targets (exercise, target) VALUES (22, 3);
INSERT into exercise_targets (exercise, target) VALUES (23, 3);
INSERT into exercise_targets (exercise, target) VALUES (24, 3);
INSERT into exercise_targets (exercise, target) VALUES (25, 3);
INSERT into exercise_targets (exercise, target) VALUES (26, 5);
INSERT into exercise_targets (exercise, target) VALUES (27, 5);
INSERT into exercise_targets (exercise, target) VALUES (28, 5);
INSERT into exercise_targets (exercise, target) VALUES (29, 5);
INSERT into exercise_targets (exercise, target) VALUES (30, 5);
INSERT into exercise_targets (exercise, target) VALUES (31, 5);
INSERT into exercise_targets (exercise, target) VALUES (32, 5);
INSERT into exercise_targets (exercise, target) VALUES (33, 5);
INSERT into exercise_targets (exercise, target) VALUES (34, 5);
INSERT into exercise_targets (exercise, target) VALUES (35, 6);
INSERT into exercise_targets (exercise, target) VALUES (36, 6);
INSERT into exercise_targets (exercise, target) VALUES (37, 6);
INSERT into exercise_targets (exercise, target) VALUES (38, 6);
INSERT into exercise_targets (exercise, target) VALUES (39, 6);
INSERT into exercise_targets (exercise, target) VALUES (40, 6);
INSERT into exercise_targets (exercise, target) VALUES (41, 6);
INSERT into exercise_targets (exercise, target) VALUES (42, 6);
INSERT into exercise_targets (exercise, target) VALUES (43, 6);
INSERT into exercise_targets (exercise, target) VALUES (44, 7);
INSERT into exercise_targets (exercise, target) VALUES (45, 7);
INSERT into exercise_targets (exercise, target) VALUES (46, 7);
INSERT into exercise_targets (exercise, target) VALUES (47, 7);

-----------------
--- EQUIPMENT ---
-----------------

DELETE from equipment;
DELETE FROM sqlite_sequence WHERE name='equipment';

INSERT into equipment (name) VALUES ('cable');
INSERT into equipment (name) VALUES ('lat pulldown machine');
INSERT into equipment (name) VALUES ('flat bench');
INSERT into equipment (name) VALUES ('incline bench');
INSERT into equipment (name) VALUES ('pullup bar');
INSERT into equipment (name) VALUES ('squat rack');
INSERT into equipment (name) VALUES ('hack squat');
INSERT into equipment (name) VALUES ('leg extension');
INSERT into equipment (name) VALUES ('leg curl');
INSERT into equipment (name) VALUES ('hip trust');
INSERT into equipment (name) VALUES ('barbell');
INSERT into equipment (name) VALUES ('dumbbell');
INSERT into equipment (name) VALUES ('cable row');
INSERT into equipment (name) VALUES ('bike');
INSERT into equipment (name) VALUES ('treadmill');
INSERT into equipment (name) VALUES ('dip bar');
INSERT into equipment (name) VALUES ('floor');
INSERT into equipment (name) VALUES ('parallettes');
INSERT into equipment (name) VALUES ('t-bar');
INSERT into equipment (name) VALUES ('z-bar');
INSERT into equipment (name) VALUES ('smith machine');
INSERT into equipment (name) VALUES ('leg press');

-------------
--- EX-EQ ---
-------------

DELETE from exercise_equipment;
DELETE FROM sqlite_sequence WHERE name='exercise_equipment';

-- Chest
INSERT INTO exercise_equipment (equipment, exercise) VALUES (4, 3); -- incline bench press
INSERT INTO exercise_equipment (equipment, exercise) VALUES (16, 1); -- dips
INSERT INTO exercise_equipment (equipment, exercise) VALUES (17, 5); -- pushups
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 6); -- machine flies
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 2); -- barbell bench press
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 6); -- dumbbell flies

-- Back
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 7); -- barbell bent-over row
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 8); -- dumbbell row
INSERT INTO exercise_equipment (equipment, exercise) VALUES (5, 10); -- pullups
INSERT INTO exercise_equipment (equipment, exercise) VALUES (19, 11); -- t-bar row
INSERT INTO exercise_equipment (equipment, exercise) VALUES (1, 9);  -- cable row (already added)
INSERT INTO exercise_equipment (equipment, exercise) VALUES (2, 12); -- lat pulldowns (already added)
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 13); -- lat pullovers
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 13); -- lat pullovers

-- Triceps
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 14); -- skull crushers
INSERT INTO exercise_equipment (equipment, exercise) VALUES (1, 15);  -- triceps extensions
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 16); -- dumbbell kickback
INSERT INTO exercise_equipment (equipment, exercise) VALUES (1, 17);  -- cable pushdowns (already added)

-- Biceps
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 18); -- dumbbell curls
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 19); -- barbell curls
INSERT INTO exercise_equipment (equipment, exercise) VALUES (1, 20);  -- cable curls (already added)
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 21); -- spider curls
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 22); -- incline curls
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 23); -- preacher curls (often on preacher bench / machine)
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 24); -- focus curls
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 25); -- hammer curls

-- Shoulders
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 26); -- lateral raises
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 27); -- military press
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 28); -- dumbbell shoulder press
INSERT INTO exercise_equipment (equipment, exercise) VALUES (17, 29); -- pike pushups
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 30); -- front raises
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 31); -- upright row
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 32); -- rear delt flies
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 33); -- rear delt row
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 34); -- LU raises

-- Legs
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 35); -- lunges
INSERT INTO exercise_equipment (equipment, exercise) VALUES (3, 36);  -- bulgarian split squat
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 37); -- barbell squat
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 38); -- goblet squat
INSERT INTO exercise_equipment (equipment, exercise) VALUES (8, 39);  -- leg extensions
INSERT INTO exercise_equipment (equipment, exercise) VALUES (17, 40); -- sissy squats
INSERT INTO exercise_equipment (equipment, exercise) VALUES (17, 41); -- pistol squats
INSERT INTO exercise_equipment (equipment, exercise) VALUES (7, 42);  -- hack squat
INSERT INTO exercise_equipment (equipment, exercise) VALUES (22, 43); -- leg press
INSERT INTO exercise_equipment (equipment, exercise) VALUES (9, 44);  -- hamstring curls
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 45); -- romanian deadlift
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 46); -- good mornings
INSERT INTO exercise_equipment (equipment, exercise) VALUES (17, 47); -- hamstring raise

-- Alts

-- Chest
INSERT INTO exercise_equipment (equipment, exercise) VALUES (18, 1);  -- dips on parallettes
INSERT INTO exercise_equipment (equipment, exercise) VALUES (18, 5);  -- pushups on parallettes

-- Back
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 7);  -- bent-over row on smith machine

-- Triceps
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 15); -- triceps extensions with barbell
INSERT INTO exercise_equipment (equipment, exercise) VALUES (12, 15); -- triceps extensions with dumbbell
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 14); -- skull crushers on smith machine

-- Biceps
INSERT INTO exercise_equipment (equipment, exercise) VALUES (20, 19); -- barbell curls with z-bar
INSERT INTO exercise_equipment (equipment, exercise) VALUES (20, 23); -- preacher curls with z-bar

-- Shoulders
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 27); -- military press on smith machine
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 28); -- shoulder press on smith machine
INSERT INTO exercise_equipment (equipment, exercise) VALUES (11, 28); -- shoulder press with barbell
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 31); -- upright row on smith machine

-- Legs
INSERT INTO exercise_equipment (equipment, exercise) VALUES (6, 37);  -- squat rack for barbell squat
INSERT INTO exercise_equipment (equipment, exercise) VALUES (21, 45); -- romanian deadlift on smith machine

-----------
--- GYM ---
-----------

DELETE from gym;
DELETE FROM sqlite_sequence WHERE name='gym';

INSERT into gym (name, location, description) values ('No Gym', '', '');
INSERT into gym (name, location, description) values ('Prime', 'Svetozara Miletica 43', 'Old school style gym with a great atmosphere');
INSERT into gym (name, location, description) values ('Muscle Gym', 'Kisacka 5', 'Idk, prolly a decent gym');
INSERT into gym (name, location, description) values ('Stevina Garaza', 'Nema ulice 130', 'Where the magic happens');

----------------
--- GYM TAGS ---
----------------

DELETE from gym_tag;
DELETE FROM sqlite_sequence WHERE name='gym_tag';

INSERT into gym_tag (name) values ('crossfit'), ('women only'), ('powerlifting'), ('mma'), ('fitpass'), ('comercial'), ('home gym');

-----------------
--- GYMS TAGS ---
-----------------

DELETE from gyms_tags;
DELETE FROM sqlite_sequence WHERE(name='gyms_tags');

INSERT into gyms_tags(gym, tag) values (2, 'powerlifting'), (2, 'fitpass'), (3, 'comerial'), (3, 'crossfit'), (3, 'fitpass'), (4, 'home gym');

-------------
--- GYM-EQ---
-------------

DELETE from gym_equipment;
DELETE FROM sqlite_sequence WHERE name='gym_equipment';

-- First gym: Prime — has **all equipment**
INSERT INTO gym_equipment (gym_id, equipment) VALUES
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6), (2, 7), (2, 8), (2, 9), (2, 10),
(2, 11), (2, 12), (2, 13), (2, 14), (2, 15), (2, 16), (2, 17), (2, 18), (2, 19), (2, 20), (2, 21), (2, 22);

-- Second gym: Muscle Gym — ~75% of equipment (random ~16 of 22)
INSERT INTO gym_equipment (gym_id, equipment) VALUES
(3, 1),  -- cable
(3, 2),  -- lat pulldown machine
(3, 3),  -- flat bench
(3, 4),  -- incline bench
(3, 5),  -- pullup bar
(3, 6),  -- squat rack
(3, 7),  -- hack squat
(3, 8),  -- leg extension
-- skipping 9 (leg curl)
(3, 10), -- hip trust
(3, 11), -- barbell
(3, 12), -- dumbbell
(3, 13), -- cable row
-- skipping 14 (bike)
(3, 15), -- treadmill
(3, 16), -- dip bar
(3, 17), -- floor
-- skipping 18 (parallettes)
(3, 19), -- t-bar
-- skipping 20 (z-bar)
(3, 21); -- smith machine
-- skipping 22 (leg press)

-- Third gym: Stevina Garaza — minimal
INSERT INTO gym_equipment (gym_id, equipment) VALUES
(4, 5),  -- pullup bar
(4, 12), -- dumbbell
(4, 17); -- floor





















































