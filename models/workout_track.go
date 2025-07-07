package models

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutTrack struct {
	Id int `json:"id"`
	Plan int `json:"plan"`
	User int `json:"user"`
	IsPrivate bool `json:"is_private"`
	WorkoutDate time.Time `json:"workout_date"`
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

	if wt.Plan <= 1 {
		return 0, errors.New("Invalid workout plan used in track")
	}

	if wt.User <= 0 {
		return 0, errors.New("Invalid user used in track")
	}

	statement := "insert into workout_track (plan, usr, is_private, workout_date) values (?, ?, ?, ?) returning id"
	var stmt *sql.Stmt
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var workout_track_id int

	err = stmt.QueryRow(
		wt.Plan,
		wt.User,
		wt.IsPrivate,
		time.Now(),
	).Scan(&workout_track_id)

	if err != nil {
		return 0, err
	}

	return workout_track_id, nil
}

func (Db *DataBase) ReadWorkoutTrack(wt_id int) (*WorkoutTrack, error) {
	wt := &WorkoutTrack{Id: wt_id}

	err := Db.Data.QueryRow("select plan, usr, is_private, workout_date from workout_track where id = ?", wt_id).Scan(
		&wt.Plan,
		&wt.User,
		&wt.IsPrivate,
		&wt.WorkoutDate,
	)

	if err != nil {
		return nil, err
	}

	return wt, nil
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
			&track.Plan,
			&track.WorkoutDate,
			&track.IsPrivate,
		)
		if err != nil {
			return nil, err
		}

		if track.IsPrivate && user_id != requesting_user_id {
			continue
		}

		tracks = append(tracks, &track)
	}

	return tracks, nil
}

func (Db *DataBase) UpdateWorkoutTrackPrivacy(wt *WorkoutTrack) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE workout_track SET is_private = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(wt.IsPrivate, wt.Id)
	if err != nil {
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

	if wt.Plan <= 1 {
		return errors.New("Cannot create track data for non-existing plan")
	}

	ex_days, err := Db.ReadAllExerciseDaysFromPlan(wt.Plan)
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
	for _, ex_day := range ex_days {
		for j := range ex_day.Sets {
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
				ex_day.Id,
				ex_day.Weight,
				j,
			).Scan(&tmp)
			if err != nil {
				return err
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

	// stmt, err := tx.Prepare("UPDATE workout_track_data SET weight = ?, rep_num = ? WHERE track = ? and ex_day = ? and set_num = ?")
	stmt, err := tx.Prepare("UPDATE workout_track_data SET weight = ?, rep_num = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	// track := tds[0].Track // WARNING: this works only assuming all of the track data is related to the same track and exercise day
	// ex_day := tds[0].ExDay

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
