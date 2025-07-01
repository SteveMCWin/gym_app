package models

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutPlan struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Creator     int    `json:"creator"`
	Description string `json:"description"`
}

type ExerciseDay struct {
	Id            int     `json:"id"`
	Plan          int     `json:"plan"`
	Exercise      int     `json:"exercise"`
	DayName       string  `json:"day"`
	Weight        float32 `json:"weight"`
	Sets          int     `json:"sets"`
	MinReps       int     `json:"min_reps"`
	MaxReps       int     `json:"max_reps"` // if == -1 then no max reps
	DayOrder      int     `json:"day_order"`
	ExerciseOrder int     `json:"exercise_order"`
}

type PlanColumn struct {
	Name string   `json:"name"`
	Rows []string `json:"rows"`
}

type PlanJSON struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	MakeCurrent bool         `json:"make_current"`
	Columns     []PlanColumn `json:"columns"`
}

func (Db *DataBase) CreateWorkoutPlan(wp *WorkoutPlan) (int, error) {

	log.Println("CREATING WORKOUT PLANNN")

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

	log.Println("WP_ID BEFORE RETURN:", workout_plan_id)

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

	stmt, err := tx.Prepare("insert into users_plans (usr, plan) values (?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		usr_id,
		plan_id,
	)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (Db *DataBase) RemoveWorkoutPlanFromUser(usr_id, plan_id int) error {
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

	tx.Commit()

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

	tx.Commit()

	return true, nil
}

func (Db *DataBase) DeleteWorkoutPlan(id *WorkoutPlan) (bool, error) { // NOTE: prolly not gonna use this at all tbh but I still gotta figure out how this is all gona work out (pun intended)
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

	tx.Commit()

	return true, nil
}

func (Db *DataBase) CreateExerciseDay(ex_day *ExerciseDay) (int, error) {

	err := ValidateExerciseDayInput(ex_day)
	if err != nil {
		return 0, err
	}

	statement := "insert into exercise_day (plan, day_name, exercise, weight, sets, min_reps, max_reps, day_order, exercise_order) values (?, ?, ?, ?, ?, ?, ?, ?, ?) returning id"
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

	err := Db.Data.QueryRow("select plan, day_name, exercise, weight, sets, min_reps, max_reps, day_order, exercise_order from workout_plan where id = ?", ex_day_id).Scan(
		&ex_day.Plan,
		&ex_day.DayName,
		&ex_day.Exercise,
		&ex_day.Weight,
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

func (Db *DataBase) UpdateExerciseDay(ex_day *ExerciseDay) error {

	err := ValidateExerciseDayInput(ex_day)
	if err != nil {
		return err
	}

	tx, err := Db.Data.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE exercise_day SET day_name = ?, exercise = ?, weight = ?, sets = ?, min_rep = ?, max_reps = ?, day_order = ?, exercise_order = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ex_day.DayName, ex_day.Exercise, ex_day.Weight, ex_day.Sets, ex_day.MinReps, ex_day.MaxReps, ex_day.DayOrder, ex_day.ExerciseOrder, ex_day.Id)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
