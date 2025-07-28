package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"fitness_app/defs"
)

type WorkoutPlan struct {
	Id          int     `json:"id"` // id of the workout plan
	Name        string  `json:"name"`
	Creator     int     `json:"creator"`
	Description string  `json:"description"`
	MakeCurrent bool    `json:"make_current"`
	Days        []ExDay `json:"days"`
}

type ExDay struct {
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
	Id           int       `json:"id"` // id of the exercise row
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ExerciseType string    `json:"exercise_type"`
	Difficulty   int       `json:"difficulty"`
	Targets      []*Target `json:"targets"`
}

type Target struct {
	Id           int         `json:"id"`
	StandardName string      `json:"standard_name"`
	LatinName    string      `json:"latin_name"`
	Exercises    []*Exercise `json:"exercises"`
}

type PlanAnalysis struct {
	SetsPerTarget map[*Target]int
	// TODO: expand analysis. check for exercise types, is max_reps high etc.
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

	err = Db.CachePlanBasic(wp.Id)
	if err != nil {
		return 0, err
	}

	err = Db.CreateExerciseDays(wp)
	if err != nil {
		return 0, err
	}

	Db.AddWorkoutPlanToUser(wp.Creator, wp.Id)

	return workout_plan_id, nil
}

func (Db *DataBase) CheckIfUserUsesPlan(usr_id, plan_id int) bool {
	var tmp int
	err := Db.Data.QueryRow("select 1 from users_plans where usr = ? AND plan = ?", usr_id, plan_id).Scan(&tmp)
	return err == nil
}

func (Db *DataBase) AddWorkoutPlanToUser(usr_id, plan_id int) error { // adds the workout to the list of workouts the user has/uses/whatever
	if plan_id == 0 || usr_id == 0 {
		return errors.New("Cannot add workout without plan_id and usr_id")
	}

	if Db.CheckIfUserUsesPlan(usr_id, plan_id) {
		log.Println("WARNING: USER WAS ALREADY LINKED WITH THIS PLAN")
		return nil
	}

	statement := "insert into users_plans (usr, plan, date_added) values (?, ?, ?)"
	stmt, err := Db.Data.Prepare(statement)
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
	rows, err := Db.Data.Query("select id, name, creator, description from users_plans inner join workout_plan on plan = id where usr = ? order by id desc", usr_id)
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

// basically a simpler version of the funciton above for when I only need the ids
func (Db *DataBase) GetPlansUserUses(user_id int) ([]int, error) {
	rows, err := Db.Data.Query("select plan from users_plans where usr = ?", user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]int, 0)

	for rows.Next() {
		var tmp int
		err = rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}

		res = append(res, tmp)
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

func (Db *DataBase) UpdateWorkoutPlan(wp *WorkoutPlan) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

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

	log.Println()
	log.Println("dif:", diff)
	log.Println()

	for i := range len(old_ex_days) {
		for j := range len(old_ex_days[i].Exercises) {
			if i >= len(diff) || j >= len(diff[i]) {
				_, err := stmt_ex.Exec(defs.PLACEHOLDER_PLAN_ID, old_ex_days[i].Exercises[j].Id)
				if err != nil {
					return false, err
				}
			} else if diff[i][j] {
				_, err := stmt_ex.Exec(defs.PLACEHOLDER_PLAN_ID, old_ex_days[i].Exercises[j].Id)
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

	err = Db.CachePlanBasic(wp.Id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) getExerciseDayDifference(new_wp *WorkoutPlan, tx *sql.Tx) ([][]bool, error) {

	search_query := "select day_name, exercise, weight, unit, sets, min_reps, max_reps from exercise_day where plan = ? and day_order = ? and exercise_order = ?"
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

	day_name_update := "update exercise_day set day_name = ? where plan = ? and day_order = ?"
	day_name_stmt, err := tx.Prepare(day_name_update)
	if err != nil {
		return nil, err
	}

	defer day_name_stmt.Close()

	diff := make([][]bool, len(new_wp.Days))

	for i, day := range new_wp.Days {
		diff[i] = make([]bool, len(new_wp.Days[i].Exercises))

		for j, ex := range day.Exercises {
			ex.Exercise.Id = exercisesByName[ex.Exercise.Name] // Needed because the new wp is just read from json so they only have their names

			var old_day_name string
			var tmp_ex ExerciseData
			err = search_stmt.QueryRow(new_wp.Id, i, j).Scan(&old_day_name, &tmp_ex.Exercise.Id, &tmp_ex.Weight, &tmp_ex.Unit, &tmp_ex.Sets, &tmp_ex.MinReps, &tmp_ex.MaxReps)
			if err != nil ||
				ex.Exercise.Id != tmp_ex.Exercise.Id ||
				ex.Weight != tmp_ex.Weight ||
				ex.Unit != tmp_ex.Unit ||
				ex.Sets != tmp_ex.Sets ||
				ex.MinReps != tmp_ex.MinReps ||
				ex.MaxReps != tmp_ex.MaxReps {
				log.Println("Error in search stmt:", err)
				diff[i][j] = true
				err := insert_stmt.QueryRow(new_wp.Id, day.Name, ex.Exercise.Id, ex.Weight, ex.Unit, ex.Sets, ex.MinReps, ex.MaxReps, i, j).Scan(&tmp_ex.Id)
				if err != nil {
					return nil, err
				}
			} else if old_day_name != day.Name {
				_, err = day_name_stmt.Exec(day.Name, new_wp.Id, i)
				if err != nil {
					return nil, err
				}
			}
		}

	}

	return diff, nil // NOTE: THE TRUE ONES ARE THE ONES THAT CHANGED
}

func (Db *DataBase) DeleteWorkoutPlans(plan_ids ...int) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

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

	stmt_usr, err := tx.Prepare("DELETE from users_plans where plan = ?")
	if err != nil {
		return false, err
	}

	defer stmt_usr.Close()

	for _, p_id := range plan_ids {
		_, err = stmt_usr.Exec(p_id)
		if err != nil {
			return false, err
		}

		_, err = stmt_day.Exec(p_id)
		if err != nil {
			return false, err
		}

		_, err = stmt_wp.Exec(p_id)
		if err != nil {
			return false, err
		}
		log.Println("Deleted a plan with id:", p_id)
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

// func (Db *DataBase) ReadExerciseDay(ex_day_id int) (*ExerciseData, error) {
// 	row, err := Db.Data.QueryRow("select plan, day")
// }

func (Db *DataBase) ReadAllExerciseDaysFromPlan(plan_id int) ([]ExDay, error) {
	rows, err := Db.Data.Query("select id, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order from exercise_day where plan = ? order by day_order asc, exercise_order asc", plan_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]ExDay, 0)

	prev_day := 0
	day := ExDay{}

	for rows.Next() {
		d := ExDay{}
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
			day = ExDay{}
		}
		day.Name = d.Name
		day.Exercises = append(day.Exercises, ex)
	}
	res = append(res, day)

	return res, nil
}

func (Db *DataBase) DeleteExerciseDay(id int) (bool, error) {

	statement := "DELETE from exercise_day where id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) CachePlanBasic(wp_id int) error {
	row := Db.Data.QueryRow("select name, creator, description from workout_plan where id = ?", wp_id)
	wp := WorkoutPlan{Id: wp_id}
	err := row.Scan(&wp.Name, &wp.Creator, &wp.Description)
	if err != nil {
		return err
	}

	err = CachePlanBasic(&wp)
	return err
}

func (wp *WorkoutPlan) GetAnalysis() *PlanAnalysis {
	if wp.Id == 1 {
		return nil
	}

	analysis := &PlanAnalysis{SetsPerTarget: make(map[*Target]int)}

	for _, day := range wp.Days {
		for _, ex := range day.Exercises {
			for _, tar := range ex.Exercise.Targets {
				analysis.SetsPerTarget[tar] += 1 * ex.Sets // NOTE: this will probably change if i add a multiplier of how much a target is being worked by each exercise
			}
		}
	}

	return analysis
}
