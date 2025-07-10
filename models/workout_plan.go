package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutPlan struct {
	Id          int       `json:"id"` // id of the workout plan
	Name        string    `json:"name"`
	Creator     int       `json:"creator"`
	Description string    `json:"description"`
	MakeCurrent bool      `json:"make_current"`
	Days        []PlanDay `json:"days"`
}

type PlanDay struct {
	Name      string         `json:"name"`
	Exercises []ExerciseData `json:"exercises"`
}

type ExerciseData struct {
	Id       int      `json:"id"` // id of the exercise_day row
	Exercise Exercise `json:"exercise"`
	Weight   float32  `json:"weight"`
	Unit     string   `json:"unit"`
	Sets     int      `json:"sets"`
	MinReps  int      `json:"min_reps"`
	MaxReps  int      `json:"max_reps"`
}

type Exercise struct {
	Id           int      `json:"id"` // id of the exercise row
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ExerciseType string   `json:"exercise_type"`
	Difficulty   int      `json:"difficulty"`
	Targets      []*Target `json:"targets"`
}

type Target struct {
	Id           int        `json:"id"`
	StandardName string     `json:"standard_name"`
	LatinName    string     `json:"latin_name"`
	Exercises    []*Exercise `json:"exercises"`
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

	wp.Id = workout_plan_id

	err = Db.CreateExerciseDays(wp)
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
		log.Println("WARNING: USER WAS ALREADY LINKED WITH THIS PLAN")
		return nil
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

	wp.Days, err = Db.ReadAllExerciseDaysFromPlan(id)
	if err != nil {
		return nil, err
	}

	return wp, nil
}

func (Db *DataBase) ReadAllWorkoutsUserUses(usr_id int) ([]*WorkoutPlan, error) {

	usr, err := Db.ReadUser(usr_id)
	if err != nil {
		return nil, err
	}

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
		current_plan.MakeCurrent = current_plan.Id != usr.CurrentPlan

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
		plan := WorkoutPlan{}
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

// NOTE: I am assuming this will be called after getting the updated version of the wp as a json
func (Db *DataBase) UpdateWorkoutPlan(wp *WorkoutPlan) (bool, error) { // WARNING: NOT TESTED AT ALL
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt_wp, err := tx.Prepare("UPDATE workout_plan SET name = ?, description = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt_wp.Close()

	stmt_ex, err := tx.Prepare("UPDATE exercise_day SET plan = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt_wp.Close()

	_, err = stmt_wp.Exec(wp.Name, wp.Description, wp.Id)
	if err != nil {
		return false, err
	}

	old_ex_days, err := Db.ReadAllExerciseDaysFromPlan(wp.Id)
	if err != nil {
		return false, err
	}

	diff, err := Db.getExerciseDayDifference(wp, tx)
	if err != nil {
		return false, err
	}

	for i := range min(len(diff), len(old_ex_days)) {
		for j := range min(len(diff[i]), len(old_ex_days[i].Exercises)) {
			if diff[i][j] {
				_, err := stmt_ex.Exec(1, old_ex_days[i].Exercises[j].Id)
				if err != nil {
					return false, err
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) getExerciseDayDifference(new_wp *WorkoutPlan, tx *sql.Tx) ([][]bool, error) { // WARNING: NOTE TESTED AT ALL

	search_query := "select day_order, exercise_order from exercise_day where plan = ? and exercise = ? and weight = ? and unit = ? and sets = ? and min_reps = ? and max_reps = ?"
	search_stmt, err := tx.Prepare(search_query)
	if err != nil {
		return nil, err
	}

	defer search_stmt.Close()

	insert_query := "insert into exercise_day (plan, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order, exercise_order) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) returning id"
	insert_stmt, err := tx.Prepare(insert_query)
	if err != nil {
		return nil, err
	}

	defer insert_stmt.Close()

	diff := make([][]bool, len(new_wp.Days))

	for i, day := range new_wp.Days {
		diff[i] = make([]bool, len(new_wp.Days[i].Exercises))

		var d_order int
		var e_order int
		for j, ex := range day.Exercises {
			err = search_stmt.QueryRow(new_wp.Id, ex.Exercise.Id, ex.Weight, ex.Unit, ex.Sets, ex.MinReps, ex.MaxReps).Scan(&d_order, &e_order)
			if err != nil || d_order != i || e_order != j {
				log.Println("Error in search stmt:", err)
				diff[i][j] = true
				err := insert_stmt.QueryRow(new_wp.Id, day.Name, ex.Exercise.Id, ex.Weight, ex.Unit, ex.Sets, ex.MinReps, ex.MaxReps, i, j).Scan(&d_order)
				if err != nil {
					return nil, err
				}
			} 	
		}

	}

	return diff, nil // NOTE: THE TRUE ONES ARE THE ONES THAT CHANGED
}

func (Db *DataBase) DeleteWorkoutPlan(wp_id int) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt_wp, err := tx.Prepare("DELETE from workout_plan where id = ?")
	if err != nil {
		return false, err
	}

	defer stmt_wp.Close()

	stmt_day, err := tx.Prepare("DELETE from exercise_day where plan = ?")
	if err != nil {
		return false, err
	}

	defer stmt_day.Close()

	stmt_usr, err := tx.Prepare("DELETE from users_plans where plan = ?") // WARNING: I AM ASSUMING THAT USING OTHER PEOPLE'S PLANS WIL GENERATE A PLAN WITH A NEW ID
	if err != nil {
		return false, err
	}

	defer stmt_usr.Close()

	_, err = stmt_usr.Exec(wp_id)
	if err != nil {
		return false, err
	}

	_, err = stmt_day.Exec(wp_id)
	if err != nil {
		return false, err
	}

	_, err = stmt_wp.Exec(wp_id)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) CreateExerciseDays(wp *WorkoutPlan) error {

	statement := "insert into exercise_day (plan, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order, exercise_order) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) returning id"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i, day := range wp.Days {
		for j, ex := range day.Exercises {

			err := ValidateExerciseDayInput(ex)
			if err != nil {
				return err
			}

			tmp, ok := FetchCachedExercise(exercisesByName[ex.Exercise.Name])
			if !ok {
				return errors.New("Couldn't fetch exercise from cached exercises")
			}

			ex.Exercise = *tmp

			err = stmt.QueryRow(
				wp.Id,
				day.Name,
				ex.Exercise.Id,
				ex.Weight,
				ex.Unit,
				ex.Sets,
				ex.MinReps,
				ex.MaxReps,
				i,
				j,
			).Scan(&ex.Id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateExerciseDayInput(ex ExerciseData) error {

	// if ex_day.Plan == 0 || ex_day.Exercise == 0 {
	// 	return errors.New("Cannot create ExerciseDay without Plan and Exercise ID")
	// }

	if ex.Sets <= 0 {
		ex.Sets = 1
	}

	if ex.Weight < 0.0 {
		ex.Weight = 0.0
	}

	if ex.MinReps < 0 {
		ex.MinReps = 0
	}

	if ex.MaxReps < 0 {
		ex.MaxReps = 0
	}

	return nil

}

// func (Db *DataBase) ReadExerciseDay(ex_day_id int) (*ExerciseDay, error) {
// 	ex_day := &ExerciseDay{Id: ex_day_id}
//
// 	err := Db.Data.QueryRow("select plan, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order, exercise_order from workout_plan where id = ?", ex_day_id).Scan(
// 		&ex_day.Plan,
// 		&ex_day.DayName,
// 		&ex_day.Exercise,
// 		&ex_day.Weight,
// 		&ex_day.Unit,
// 		&ex_day.Sets,
// 		&ex_day.MinReps,
// 		&ex_day.MaxReps,
// 		&ex_day.DayOrder,
// 		&ex_day.ExerciseOrder,
// 	)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return ex_day, nil
// }

func (Db *DataBase) ReadAllExerciseDaysFromPlan(plan_id int) ([]PlanDay, error) {
	rows, err := Db.Data.Query("select id, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order from exercise_day where plan = ? order by day_order asc, exercise_order asc", plan_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]PlanDay, 0)

	prev_day := 0
	day := PlanDay{}

	for rows.Next() {
		d := PlanDay{}
		ex := ExerciseData{}
		var curr_day int
		err = rows.Scan(
			&ex.Id,
			&d.Name,
			&ex.Exercise.Id,
			&ex.Weight,
			&ex.Unit,
			&ex.Sets,
			&ex.MinReps,
			&ex.MaxReps,
			&curr_day,
		)
		if err != nil {
			return nil, err
		}


		tmp, ok := FetchCachedExercise(ex.Exercise.Id)
		if !ok {
			return nil, err
		}
		ex.Exercise = *tmp

		if curr_day != prev_day {
			prev_day = curr_day
			res = append(res, day)
			day = PlanDay{}
		}
		day.Name = d.Name
		day.Exercises = append(day.Exercises, ex)
	}
	res = append(res, day)

	return res, nil
}

// func (Db *DataBase) UpdateExerciseDayExercise(ex_day *ExerciseDay) error {
//
// 	// err := ValidateExerciseDayInput(ex_day)
// 	// if err != nil {
// 	// 	return err
// 	// }
//
// 	tx, err := Db.Data.Begin()
// 	if err != nil {
// 		return err
// 	}
//
// 	stmt, err := tx.Prepare("UPDATE exercise_day SET exercise = ?, weight = ?, unit = ?, sets = ?, min_rep = ?, max_reps = ? WHERE id = ?")
// 	if err != nil {
// 		return err
// 	}
//
// 	defer stmt.Close()
//
// 	_, err = stmt.Exec(ex_day.Exercise, ex_day.Weight, ex_day.Unit, ex_day.Sets, ex_day.MinReps, ex_day.MaxReps, ex_day.Id)
//
// 	if err != nil {
// 		return err
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// func (Db *DataBase) UpdateExerciseDayOrder(ex_day *ExerciseDay, old_day_order, old_ex_order int) error {
//
// 	changed_day_oder := old_day_order != ex_day.DayOrder
//
// 	var new_day_name string
//
// 	if changed_day_oder {
// 		err := Db.Data.QueryRow("select day_name from exercise_day where plan = ? and day_order = ? and exercise_order = 0", ex_day.Plan, ex_day.DayOrder).Scan(&new_day_name)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		new_day_name = ex_day.DayName
// 	}
//
// 	tx, err := Db.Data.Begin()
// 	if err != nil {
// 		return err
// 	}
//
// 	var stmt_update_day *sql.Stmt
// 	if changed_day_oder {
// 		stmt_update_day, err = tx.Prepare("UPDATE exercise_day SET exercise_order = exercise_order + 1 WHERE plan = ? AND day_order = ? AND exercise_order >= ?")
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		stmt_update_day, err = tx.Prepare("UPDATE exercise_day SET exercise_order = exercise_order + 1 WHERE plan = ? AND day_order = ? AND exercise_order >= ? and exercise_order < ?")
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	stmt_update_ex, err := tx.Prepare("UPDATE exercise_day SET day_name = ?, day_order = ?, exercise_order = ? WHERE id = ?")
// 	if err != nil {
// 		return err
// 	}
//
// 	defer stmt_update_day.Close()
// 	defer stmt_update_ex.Close()
//
// 	if changed_day_oder {
// 		_, err = stmt_update_day.Exec(ex_day.Plan, ex_day.DayOrder, ex_day.ExerciseOrder)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		_, err = stmt_update_day.Exec(ex_day.Plan, ex_day.DayOrder, ex_day.ExerciseOrder, old_ex_order)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	_, err = stmt_update_ex.Exec(new_day_name, ex_day.DayOrder, ex_day.ExerciseOrder, ex_day.Id)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
//
// }

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

func (Db *DataBase) CacheAllExercises() error {

	rows, err := Db.Data.Query("select id, name, description, exercise_type, difficulty from exercises")
	if err != nil {
		return err
	}

	defer rows.Close()

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
			return err
		}

		err := CacheExercise(&current_exercise)
		if err != nil {
			return err
		}
	}

	return nil
}

func (Db *DataBase) CacheAllTargets() error {

	rows, err := Db.Data.Query("select id, standard_name, latin_name from target")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		current_target := Target{}

		err = rows.Scan(
			&current_target.Id,
			&current_target.StandardName,
			&current_target.LatinName,
		)
		if err != nil {
			return err
		}

		err := CacheTarget(&current_target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (Db *DataBase) LinkCachedExercisesAndTargets() error {

	rows, err := Db.Data.Query("select exercise, target from exercise_targets")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var ex int
		var tar int

		err = rows.Scan(
			&ex,
			&tar,
		)
		if err != nil {
			return err
		}

		_, ok := cachedExercises[ex]
		if !ok {
			return errors.New("Couldn't link exercise and target: missing exercise")
		}
		_, ok = cachedTargets[tar]
		if !ok {
			return errors.New("Couldn't link exercise and target: missing target")
		}

		cachedExercises[ex].Targets = append(cachedExercises[ex].Targets, cachedTargets[tar])
		cachedTargets[tar].Exercises = append(cachedTargets[tar].Exercises, cachedExercises[ex])
	}

	return nil
}
