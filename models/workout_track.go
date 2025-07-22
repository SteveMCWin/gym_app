package models

import (
	"log"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutTrack struct {
	Id int `json:"id"`
	Plan WorkoutPlan `json:"plan"`
	User int `json:"user"`
	IsPrivate bool `json:"is_private"`
	WorkoutDate time.Time `json:"workout_date"`
	ExDays []ExDay `json:"ex_days"`
}

type TrackData struct {
	Id int `json:"id"`
	Track int `json:"track"`
	ExDay int `json:"ex_day"`
	Weight float32 `json:"weight"`
	SetNum int `json:"set_num"`
	RepNum int `json:"rep_num"`
}

type TrackJSON struct {
	WTrack WorkoutTrack `json:"wt"`
	Data []TrackData `json:"td"`
}

func (Db *DataBase) CreateWorkoutTrack(wt *WorkoutTrack) (int, error) {

	if wt.Plan.Id <= 1 {
		return 0, errors.New("Invalid workout plan used in track")
	}

	if wt.User <= 0 {
		return 0, errors.New("Invalid user used in track")
	}

	statement_wt := "insert into workout_track (plan, usr, is_private, workout_date) values (?, ?, ?, ?) returning id"
	var stmt_wt *sql.Stmt
	stmt_wt, err := Db.Data.Prepare(statement_wt)
	if err != nil {
		return 0, err
	}

	defer stmt_wt.Close()

	statement_te := "insert into track_exercise (track, ex_day) values (?, ?) returning 1"
	stmt_te, err := Db.Data.Prepare(statement_te)
	if err != nil {
		return 0, err
	}

	defer stmt_te.Close()

	var workout_track_id int

	err = stmt_wt.QueryRow(
		wt.Plan.Id,
		wt.User,
		wt.IsPrivate,
		time.Now(),
	).Scan(&workout_track_id)

	if err != nil {
		log.Println("HERE 1")
		return 0, err
	}

	ex_days, err := Db.ReadAllExerciseDaysFromPlan(wt.Plan.Id)
	if err != nil {
		log.Println("HERE 2")
		return 0, err
	}

	var tmp int
	for _, day := range ex_days {
		for _, ex := range day.Exercises {
			err = stmt_te.QueryRow(workout_track_id, ex.Id).Scan(&tmp) // Perhaps this should be an Exec and not QueryRow
			if err != nil {
				log.Println("HERE 3")
				return 0, err
			}
		}
	}

	return workout_track_id, nil
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

	if err != nil {
		return nil, err
	}

	wt.ExDays, err = Db.ReadAllExerciseDaysFromTrack(wt_id)

	if err != nil {
		return nil, err
	}

	return wt, nil
}

func (Db *DataBase) ReadAllExerciseDaysFromTrack(wt_id int) ([]ExDay, error) {
	querry := `
	select id, day_name, exercise, weight, unit, sets, min_reps, max_reps, day_order
	from exercise_day inner join track_exercise on id = ex_day
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
	rows, err := Db.Data.Query("select id, plan, workout_date, is_private from workout_track where usr = ? order by workout_date desc", user_id) // NOTE: add ordering options
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tracks := make([]*WorkoutTrack, 0)

	for rows.Next() {
		track := WorkoutTrack { User: user_id }

		err = rows.Scan(
			&track.Id,
			&track.Plan.Id,
			&track.WorkoutDate,
			&track.IsPrivate,
		)
		if err != nil {
			return nil, err
		}

		if track.IsPrivate && user_id != requesting_user_id {
			continue
		}

		plan_basic, ok := FetchCachedPlanBasic(track.Plan.Id)
		if !ok {
			return nil, errors.New("Couldn't read cached plan data")
		}
		track.Plan = *plan_basic

		tracks = append(tracks, &track)
	}

	return tracks, nil
}

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

func (Db *DataBase) CreateTrackData(td *TrackData) (int, error) {

	if td.Track <= 1 {
		return 0, errors.New("Invalid workout plan used in track")
	}

	if td.ExDay <= 0 {
		return 0, errors.New("Invalid user used in track")
	}

	statement := "insert into workout_track_data (track, ex_day, weight, set_num) values (?, ?, ?, ?) returning id"
	var stmt *sql.Stmt
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var track_data_id int

	err = stmt.QueryRow(
		td.Track,
		td.ExDay,
		td.Weight,
		td.SetNum,
	).Scan(&track_data_id)

	if err != nil {
		return 0, err
	}

	return track_data_id, nil
}

func (Db *DataBase) CreateTrackDataForTrack(wt *WorkoutTrack) error {

	if wt.Plan.Id <= 1 {
		return errors.New("Cannot create track data for non-existing plan")
	}

	days, err := Db.ReadAllExerciseDaysFromPlan(wt.Plan.Id)
	if err != nil {
		return err
	}

	statement := "insert into workout_track_data (track, ex_day, weight, set_num) values (?, ?, ?, ?) returning id"
	var stmt *sql.Stmt
	stmt, err = Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, day := range days {
		for _, ex := range day.Exercises {
			for n := range ex.Sets {
				// td := TrackData {
				// 	Track: wt.Id,
				// 	ExDay: ex_day.Id,
				// 	Weight: ex_day.Weight,
				// 	SetNum: j,
				// }
				// _, err = Db.CreateTrackData(&td) // NOTE: perhaps may be done in a goroutine

				var tmp int

				err = stmt.QueryRow(
					wt.Id,
					ex.Id,
					ex.Weight,
					n,
				).Scan(&tmp)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (Db *DataBase) ReadTrackDataForTrack(wt_id int) ([]*TrackData, error){
	query_str := `
	select id, ex_day, weight, set_num, rep_num 
	from workout_track_data 
	where track = ? 
	` // WARNING: will probably have to order this depending on the exercise order that is defined by the set and day order in the ex_day
	// But perhaps may be realized with just comparing the ex_day from here with the ex_day.id from the ex_day, we'll see

	rows, err := Db.Data.Query(query_str, wt_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := make([]*TrackData, 0)

	for rows.Next() {
		d := TrackData{ Track: wt_id }

		err = rows.Scan(
			&d.Id,
			&d.ExDay,
			&d.Weight,
			&d.SetNum,
			&d.RepNum,
		)
		if err != nil {
			return nil, err
		}

		data = append(data, &d)
	}

	return data, nil
}

func (Db *DataBase) UpdateTrackData(td *TrackData) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE workout_track_data SET weight = ?, rep_num = ? WHERE track = ? and ex_day = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(td.Weight, td.RepNum, td.Track, td.ExDay)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
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

	// stmt, err := tx.Prepare("UPDATE workout_track_data SET weight = ?, rep_num = ? WHERE track = ? and ex_day = ? and set_num = ?")
	stmt, err := tx.Prepare("UPDATE workout_track_data SET weight = ?, rep_num = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	for _, td := range tds {
		_, err = stmt.Exec(td.Weight, td.RepNum, td.Id)
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

func (Db *DataBase) GetTracksUserUses(user_id int) ([]int, error) {
	rows, err := Db.Data.Query("select id from workout_track where usr = ?", user_id)
	if err != nil {
		return nil, err
	}

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

func (Db *DataBase) DeleteWorkoutTracks(track_ids ...int) (bool, error) {

	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	tracks := make([]*WorkoutTrack, 0)

	for _, t_id := range track_ids {
		track, err := Db.ReadWorkoutTrack(t_id)
		if err != nil {
			return false, err
		}
		tracks = append(tracks, track)
	}

	stmt_track, err := tx.Prepare("DELETE from workout_track where id = ?")
	if err != nil {
		return false, err
	}
	defer stmt_track.Close()

	stmt_track_ex, err := tx.Prepare("DELETE from track_exercise where track = ?")
	if err != nil {
		return false, err
	}
	defer stmt_track_ex.Close()

	stmt_track_data, err := tx.Prepare("DELETE from workout_track_data where track = ?")
	if err != nil {
		return false, err
	}
	defer stmt_track_data.Close()

	stmt_ex_days, err := tx.Prepare("DELETE from exercise_day where id = ? and plan = 1 and id not in (select ex_day from track_exercise where ex_day = ?)") // NOTE: remember to test if the ex_days are deleted when not used in any track or plan
	if err != nil {
		return false, err
	}
	defer stmt_ex_days.Close()

	for _, track := range tracks {
		_, err = stmt_track.Exec(track.Id)
		if err != nil {
			return false, err
		}

		_, err = stmt_track_ex.Exec(track.Id)
		if err != nil {
			return false, err
		}

		_, err = stmt_track_data.Exec(track.Id)
		if err != nil {
			return false, err
		}

		for _, d := range track.ExDays {
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
