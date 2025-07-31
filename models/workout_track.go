package models

import (
	"database/sql"
	"errors"
	"fitness_app/defs"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutTrack struct {
	Id          int         `json:"id"`
	Plan        WorkoutPlan `json:"plan_id"` // Used when creating the track and comparing the planned and tracked data
	User        int         `json:"user"`
	IsPrivate   bool        `json:"is_private"`
	WorkoutDate time.Time   `json:"workout_date"`
	TrackData   []TrackData `josn:"track_data"` // The actual data the user is interested in
}

type TrackData struct {
	Id           int          `json:"id"`
	ExDayId      int          `json:"ex_day_id"`
	Weight       float32      `json:"weight"`
	SetNum       int          `json:"set_num"`
	RepNum       int          `json:"rep_num"`
	TimeRecorded sql.NullTime `json:"time_recorded"`
}

type ExerciseProgress struct {
	Intensity float32 `json:"intensity"`
	TimePoint time.Time `json:"time_point"`
}

func (Db *DataBase) CreateWorkoutTrack(wt *WorkoutTrack) (int, error) {

	if wt.Plan.Id <= 1 {
		return 0, errors.New("Invalid workout plan used in track")
	}

	if wt.User <= 0 {
		return 0, errors.New("Invalid user used in track")
	}

	statement_wt := "insert into workout_track (plan, usr, is_private, workout_date) values (?, ?, ?, ?) returning id"
	stmt_wt, err := Db.Data.Prepare(statement_wt)
	if err != nil {
		return 0, err
	}

	defer stmt_wt.Close()

	err = stmt_wt.QueryRow(
		wt.Plan.Id,
		wt.User,
		wt.IsPrivate,
		time.Now(),
	).Scan(&wt.Id)

	if err != nil {
		return 0, err
	}

	err = Db.CreateTrackDataForTrack(wt)
	if err != nil {
		return 0, err
	}

	return wt.Id, nil
}

func (Db *DataBase) ReadWorkoutTrack(wt_id int) (*WorkoutTrack, error) {
	wt := &WorkoutTrack{Id: wt_id}

	err := Db.Data.QueryRow("select plan, usr, is_private, workout_date from workout_track where id = ?", wt_id).Scan(
		&wt.Plan.Id,
		&wt.User,
		&wt.IsPrivate,
		&wt.WorkoutDate,
	)

	wt.Plan = *cacehdPlansBasic[wt.Plan.Id]
	wt.Plan.Days, err = Db.ReadAllExerciseDaysFromTrack(wt.Id) // Note: the ex days are not as the plan is currently but as the plan was when the track was created
	if err != nil {
		return nil, err
	}

	wt.TrackData, err = Db.ReadTrackDataForTrack(wt.Id)
	if err != nil {
		return nil, err
	}

	return wt, nil
}

func (Db *DataBase) ReadAllExerciseDaysFromTrack(wt_id int) ([]ExDay, error) {
	querry := `
	select distinct exercise_day.id, day_name, exercise, exercise_day.weight, unit, sets, min_reps, max_reps, day_order
	from exercise_day inner join workout_track_data on exercise_day.id = ex_day
	where track = ?
	order by day_order asc, exercise_order asc
	`
	rows, err := Db.Data.Query(querry, wt_id)
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

func (Db *DataBase) ReadUsersWorkoutTracks(user_id, requesting_user_id int) ([]*WorkoutTrack, error) {
	rows, err := Db.Data.Query("select id, is_private from workout_track where usr = ? order by workout_date desc", user_id) // TODO: add ordering options
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tracks := make([]*WorkoutTrack, 0)

	for rows.Next() {
		var wt_id int
		var is_private bool

		err = rows.Scan(
			&wt_id,
			&is_private,
		)
		if err != nil {
			return nil, err
		}

		if is_private && user_id != requesting_user_id {
			continue
		}

		track, err := Db.ReadWorkoutTrack(wt_id)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

// func (Db *DataBase) GetTracksUserUses(user_id int) ([]int, error) { // TODO: delete either this one or the one above. And while at it go thorugh all the models and delete duplicates
// 	rows, err := Db.Data.Query("select id from workout_track where usr = ?", user_id)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	res := make([]int, 0)
//
// 	for rows.Next() {
// 		var tmp int
// 		err = rows.Scan(&tmp)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		res = append(res, tmp)
// 	}
//
// 	return res, nil
// }

func (Db *DataBase) UpdateWorkoutTrackPrivacy(wt *WorkoutTrack) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE workout_track SET is_private = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(wt.IsPrivate, wt.Id)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			log.Println("Couldn't rollback?")
		}
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) CreateTrackDataForTrack(wt *WorkoutTrack) error {

	if wt.Plan.Id <= defs.PLACEHOLDER_PLAN_ID {
		return errors.New("Cannot create track data for non-existing plan")
	}

	days, err := Db.ReadAllExerciseDaysFromPlan(wt.Plan.Id)
	if err != nil {
		return err
	}

	statement := "insert into workout_track_data (track, ex_day, weight, set_num) values (?, ?, ?, ?)"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, day := range days {
		for _, ex := range day.Exercises {
			for n := range ex.Sets {

				_, err = stmt.Exec(
					wt.Id,
					ex.Id,
					ex.Weight,
					n,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (Db *DataBase) ReadTrackDataForTrack(wt_id int) ([]TrackData, error) {
	query_str := `
	select workout_track_data.id, ex_day, workout_track_data.weight, set_num, rep_num, time_recorded
	from workout_track_data inner join exercise_day on ex_day = exercise_day.id
	where track = ? 
	order by day_order asc, exercise_order asc
	`

	rows, err := Db.Data.Query(query_str, wt_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := make([]TrackData, 0)

	for rows.Next() {
		td := TrackData{}

		err = rows.Scan(
			&td.Id,
			&td.ExDayId,
			&td.Weight,
			&td.SetNum,
			&td.RepNum,
			&td.TimeRecorded,
		)
		if err != nil {
			return nil, err
		}

		data = append(data, td)
	}

	return data, nil
}

func (Db *DataBase) GetAllDataForExerciseAndUser(user_id, exercise_id int) ([]*TrackData, error) {
	query := `
	select wtd.weight, wtd.set_num, wtd.rep_num, wtd.time_recorded
	from workout_track_data wtd inner join workout_track wt on wt.id = wtd.track inner join exercise_day ed on wtd.ex_day = ed.id
	where wt.usr = ? and ed.exercise = ? and wtd.time_recorded is not null
	order by wt.workout_date asc, wtd.time_recorded asc, wtd.set_num asc
	`
	rows, err := Db.Data.Query(query, user_id, exercise_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]*TrackData, 0)

	for rows.Next() {
		tmp := &TrackData{}
		err = rows.Scan(
			&tmp.Weight,
			&tmp.SetNum,
			&tmp.RepNum,
			&tmp.TimeRecorded,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, tmp)
	}

	return res, nil
}

func CalcExerciseProgressFromTrackData(data []*TrackData) []*ExerciseProgress {

	if len(data) <= 0 {
		return nil
	}

	res := make([]*ExerciseProgress, 0)

	curr_set_num := -1
	curr_ex_prog := new(ExerciseProgress)

	for _, td := range data {
		if td.SetNum <= curr_set_num {
			res = append(res, curr_ex_prog)
			curr_ex_prog = new(ExerciseProgress)
		}

		weight := td.Weight
		if td.Weight <= 0 {
			weight = 1.0
		}

		if td.TimeRecorded.Valid {
			curr_ex_prog.TimePoint = td.TimeRecorded.Time
			curr_ex_prog.Intensity += float32(td.RepNum) * weight
		}
		curr_set_num = td.SetNum
	}

	if curr_ex_prog.Intensity > 0.0 {
		res = append(res, curr_ex_prog)
	}

	return res
}

func (Db *DataBase) UpdateMultipleTrackData(tds []TrackData) (bool, error) {

	if len(tds) <= 0 {
		return false, errors.New("Empty track data slice passed")
	}

	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	statement := `
	update workout_track_data
	set 
		weight = ?,
		rep_num = ?,
		time_recorded = case
			when time_recorded is null then ?
			else time_recorded
		end
	where id = ?
	`

	stmt, err := tx.Prepare(statement)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	for _, td := range tds {
		_, err = stmt.Exec(td.Weight, td.RepNum, time.Now(), td.Id)
		if err != nil {
			return false, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (Db *DataBase) DeleteWorkoutTracks(tracks ...*WorkoutTrack) (bool, error) {

	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	// tracks := make([]*WorkoutTrack, 0)

	// for _, t_id := range track_ids {
	// 	track, err := Db.ReadWorkoutTrack(t_id)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	tracks = append(tracks, track)
	// }

	stmt_track, err := tx.Prepare("DELETE from workout_track where id = ?")
	if err != nil {
		return false, err
	}
	defer stmt_track.Close()

	stmt_track_data, err := tx.Prepare("DELETE from workout_track_data where track = ?")
	if err != nil {
		return false, err
	}
	defer stmt_track_data.Close()

	stmt_ex_days, err := tx.Prepare("DELETE from exercise_day where id = ? and plan = 1 and id not in (select ex_day from workout_track_data where ex_day = ?)") // NOTE: remember to test if the ex_days are deleted when not used in any track or plan
	if err != nil {
		return false, err
	}
	defer stmt_ex_days.Close()

	for _, track := range tracks {
		_, err = stmt_track.Exec(track.Id)
		if err != nil {
			return false, err
		}

		_, err = stmt_track_data.Exec(track.Id)
		if err != nil {
			return false, err
		}

		for _, d := range track.Plan.Days {
			for _, e := range d.Exercises {
				_, err = stmt_ex_days.Exec(e.Id, e.Id)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
