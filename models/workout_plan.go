package models

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutPlan struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Creator     int    `json:"creator"`
	Description string `json:"description"`
}

// NOTE: Consider ditching these two and just use the PlanJSON

type ExerciseDay struct {
	Id            int     `json:"id"`
	Plan          int     `json:"plan"`
	Exercise      int     `json:"exercise"`
	DayName       string  `json:"day"`
	Weight        float32 `json:"weight"`
	Unit          string  `json:"unit"`
	Sets          int     `json:"sets"`
	MinReps       int     `json:"min_reps"`
	MaxReps       int     `json:"max_reps"` // if == -1 then no max reps
	DayOrder      int     `json:"day_order"`
	ExerciseOrder int     `json:"exercise_order"`
}

type PlanRow struct {
	Name    string  `json:"name"`
	Weight  float32 `json:"weight"`
	Unit    string  `json:"unit"`
	Sets    int     `json:"sets"`
	MinReps int     `json:"min_reps"`
	MaxReps *int    `json:"max_reps"`
}

type PlanColumn struct {
	Name string    `json:"name"`
	Rows []PlanRow `json:"rows"`
}

type PlanJSON struct {
	Id          int          `json:"id"`
	// NOTE: Perhaps add an exercise id field
	Name        string       `json:"name"`
	Description string       `json:"description"`
	MakeCurrent bool         `json:"make_current"`
	Columns     []PlanColumn `json:"columns"`
}

type Exercise struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ExerciseType string `json:"exercise_type"`
	Difficulty int `json:"difficulty"`
}

func (Db *DataBase) CreateWorkoutPlan(wp *WorkoutPlan) (int, error) {

	if wp.Creator == 0 {
		return 0, errors.New("You need to be logged in to create a workout plan")
	}

	if wp.Name == "" {
		wp.Name = "WorkoutPlan"
	}

	err := Db.Data.QueryRow("select id from workout_plan where name like ? AND creator like ?", wp.Name, wp.Creator).Scan()

	if err == nil {
		// A plan with a name like this already exists, add _2 at the end
		// NOTE: Should probably notify the user this is happening
		wp.Name = wp.Name + "_2"
	}

	statement := "insert into workout_plan (name, creator, description) values (?, ?, ?) returning id"
	var stmt *sql.Stmt
	stmt, err = Db.Data.Prepare(statement)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var workout_plan_id int

	err = stmt.QueryRow(
		wp.Name,
		wp.Creator,
		wp.Description,
	).Scan(&workout_plan_id)

	if err != nil {
		return 0, err
	}

	err = Db.AddWorkoutPlanToUser(wp.Creator, workout_plan_id)
	if err != nil {
		return 0, err
	}

	return workout_plan_id, nil
}

func (Db *DataBase) AddWorkoutPlanToUser(usr_id, plan_id int) error { // adds the workout to the list of workouts the user has/uses/whatever
	if plan_id == 0 || usr_id == 0 {
		return errors.New("Cannot add workout without plan_id and usr_id")
	}

	var tmp int

	err := Db.Data.QueryRow("select plan from users_plans where usr = ? AND plan = ?", usr_id, plan_id).Scan(&tmp)

	if err == nil {
		return errors.New("User was already linked with this workout plan")
	}

	tx, err := Db.Data.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into users_plans (usr, plan, date_added) values (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		usr_id,
		plan_id,
		time.Now(),
	)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (Db *DataBase) DeleteWorkoutPlanFromUser(usr_id, plan_id int) error {
	tx, err := Db.Data.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE from workout_plan where user = ? AND plan = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(usr_id, plan_id)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (Db *DataBase) ReadWorkoutPlan(id int) (*WorkoutPlan, error) {
	wp := &WorkoutPlan{Id: id}

	err := Db.Data.QueryRow("select name, creator, description from workout_plan where id = ?", id).Scan(
		&wp.Name,
		&wp.Creator,
		&wp.Description,
	)

	if err != nil {
		return nil, err
	}

	return wp, nil
}

func (Db *DataBase) ReadAllWorkoutsUserUses(usr_id int) ([]*WorkoutPlan, error) {

	// NOTE: figure out how to get the date added. 
	rows, err := Db.Data.Query("select id, name, creator, description from users_plans inner join workout_plan on plan = id where usr = ?", usr_id) 
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]*WorkoutPlan, 0)

	for rows.Next() {
		current_plan := WorkoutPlan{}

		err = rows.Scan(
			&current_plan.Id,
			&current_plan.Name,
			&current_plan.Creator,
			&current_plan.Description,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, &current_plan)
	}

	return res, nil
}

func (Db *DataBase) ReadUsersRecentlyTrackedPlans(user_id int) ([]*WorkoutPlan, error) {

	sql_query := `
	select workout_plan.id, workout_plan.name, workout_plan.creator, workout_plan.description, max(workout_date)
	from workout_track inner join workout_plan on plan = workout_plan.id
	where usr = ?
	group by plan
	order by max(workout_date) asc
	` // Doesn't work because the workout_track table is empty at first

	sql_query2 := `
	select id, name, creator, description
	from workout_plan inner join users_plans on id = plan
	where usr = ?
	order by date_added desc
	`

	_ = sql_query
	_ = sql_query2

	rows, err := Db.Data.Query(sql_query2, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	plans := make([]*WorkoutPlan, 0)

	for rows.Next() {
		plan := WorkoutPlan {}
		// var last_used time.Time

		err = rows.Scan(
			&plan.Id,
			&plan.Name,
			&plan.Creator,
			&plan.Description,
			// &last_used,
		)
		if err != nil {
			return nil, err
		}

		plans = append(plans, &plan)
	}

	return plans, nil
}

func (Db *DataBase) UpdateWorkoutPlan(wp *WorkoutPlan) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE workout_plan SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(wp.Name, wp.Description)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) DeleteWorkoutPlan(id int) (bool, error) { // NOTE: prolly not gonna use this at all tbh but I still gotta figure out how this is all gona work out (pun intended)
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from workout_plan where id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) CreateExerciseDay(ex_day *ExerciseDay) (int, error) {

	err := ValidateExerciseDayInput(ex_day)
	if err != nil {
		return 0, err
	}

	statement := "insert into exercise_day (plan, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order, exercise_order) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) returning id"
	var stmt *sql.Stmt
	stmt, err = Db.Data.Prepare(statement)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var ex_day_id int
	err = stmt.QueryRow(
		ex_day.Plan,
		ex_day.DayName,
		ex_day.Exercise,
		ex_day.Weight,
		ex_day.Unit,
		ex_day.Sets,
		ex_day.MinReps,
		ex_day.MaxReps,
		ex_day.DayOrder,
		ex_day.ExerciseOrder,
	).Scan(&ex_day_id)

	if err != nil {
		return 0, err
	}

	return ex_day_id, nil
}

func ValidateExerciseDayInput(ex_day *ExerciseDay) error {

	if ex_day.Plan == 0 || ex_day.Exercise == 0 {
		return errors.New("Cannot create ExerciseDay without Plan and Exercise ID")
	}

	if ex_day.Sets < 0 {
		ex_day.Sets = 1
	}

	if ex_day.Weight < 0.0 {
		ex_day.Weight = 0.0
	}

	if ex_day.MinReps < 0 {
		ex_day.MinReps = 0
	}

	if ex_day.MaxReps < 0 {
		ex_day.MaxReps = 0
	}

	return nil

}

func (Db *DataBase) ReadExerciseDay(ex_day_id int) (*ExerciseDay, error) {
	ex_day := &ExerciseDay{Id: ex_day_id}

	err := Db.Data.QueryRow("select plan, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order, exercise_order from workout_plan where id = ?", ex_day_id).Scan(
		&ex_day.Plan,
		&ex_day.DayName,
		&ex_day.Exercise,
		&ex_day.Weight,
		&ex_day.Unit,
		&ex_day.Sets,
		&ex_day.MinReps,
		&ex_day.MaxReps,
		&ex_day.DayOrder,
		&ex_day.ExerciseOrder,
	)

	if err != nil {
		return nil, err
	}

	return ex_day, nil
}

func (Db *DataBase) ReadAllExerciseDaysFromPlan(plan_id int) ([]*ExerciseDay, error) {
	rows, err := Db.Data.Query("select id, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order, exercise_order from exercise_day where plan = ? order by day_order asc, exercise_order asc", plan_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]*ExerciseDay, 0)

	for rows.Next() {
		ex_day := ExerciseDay{Plan: plan_id}

		err = rows.Scan(
			&ex_day.Id,
			&ex_day.DayName,
			&ex_day.Exercise,
			&ex_day.Weight,
			&ex_day.Unit,
			&ex_day.Sets,
			&ex_day.MinReps,
			&ex_day.MaxReps,
			&ex_day.DayOrder,
			&ex_day.ExerciseOrder,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, &ex_day)
	}

	return res, nil
}

func (Db *DataBase) UpdateExerciseDayExercise(ex_day *ExerciseDay) error {

	err := ValidateExerciseDayInput(ex_day)
	if err != nil {
		return err
	}

	tx, err := Db.Data.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE exercise_day SET exercise = ?, weight = ?, unit = ?, sets = ?, min_rep = ?, max_reps = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ex_day.Exercise, ex_day.Weight, ex_day.Unit, ex_day.Sets, ex_day.MinReps, ex_day.MaxReps, ex_day.Id)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (Db *DataBase) UpdateExerciseDayOrder(ex_day *ExerciseDay, old_day_order, old_ex_order int) error {

	changed_day_oder := old_day_order != ex_day.DayOrder

	var new_day_name string

	if changed_day_oder {
		err := Db.Data.QueryRow("select day_name from exercise_day where plan = ? and day_order = ? and exercise_order = 0", ex_day.Plan, ex_day.DayOrder).Scan(&new_day_name)
		if err != nil {
			return err
		}
	} else {
		new_day_name = ex_day.DayName
	}

	tx, err := Db.Data.Begin()
	if err != nil {
		return err
	}

	var stmt_update_day *sql.Stmt
	if changed_day_oder {
		stmt_update_day, err = tx.Prepare("UPDATE exercise_day SET exercise_order = exercise_order + 1 WHERE plan = ? AND day_order = ? AND exercise_order >= ?")
		if err != nil {
			return err
		}
	} else {
		stmt_update_day, err = tx.Prepare("UPDATE exercise_day SET exercise_order = exercise_order + 1 WHERE plan = ? AND day_order = ? AND exercise_order >= ? and exercise_order < ?")
		if err != nil {
			return err
		}
	}

	stmt_update_ex, err := tx.Prepare("UPDATE exercise_day SET day_name = ?, day_order = ?, exercise_order = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt_update_day.Close()
	defer stmt_update_ex.Close()

	if changed_day_oder {
		_, err = stmt_update_day.Exec(ex_day.Plan, ex_day.DayOrder, ex_day.ExerciseOrder)
		if err != nil {
			return err
		}
	} else {
		_, err = stmt_update_day.Exec(ex_day.Plan, ex_day.DayOrder, ex_day.ExerciseOrder, old_ex_order)
		if err != nil {
			return err
		}
	}

	_, err = stmt_update_ex.Exec(new_day_name, ex_day.DayOrder, ex_day.ExerciseOrder, ex_day.Id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (Db *DataBase) DeleteExerciseDay(id int) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from exercise_day where id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) ReadAllExercises() ([]*Exercise, error) {

	rows, err := Db.Data.Query("select id, name, description, exercise_type, difficulty from exercises") 
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]*Exercise, 0)

	for rows.Next() {
		current_exercise := Exercise{}

		err = rows.Scan(
			&current_exercise.Id,
			&current_exercise.Name,
			&current_exercise.Description,
			&current_exercise.ExerciseType,
			&current_exercise.Difficulty,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, &current_exercise)
	}

	return res, nil
}
