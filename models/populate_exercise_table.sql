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
