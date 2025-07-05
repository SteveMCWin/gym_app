package models

import (
	"database/sql"
	"errors"
	// "log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutTrack struct {
	Id int
	Plan int
	User int
	IsPrivate bool
	WorkoutDate time.Time
}

type TrackData struct {
	Id int
	Track int
	ExDay int
	Weight float32
	SetNum int
	RepNum int
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
	).Scan(&workout_track_id)

	if err != nil {
		return 0, err
	}

	return workout_track_id, nil
}

func (Db *DataBase) ReadUsersWorkoutTracks(user_id, requesting_user_id int) ([]*WorkoutTrack, error) {
	rows, err := Db.Data.Query("select id, plan, workout_date, is_private from workout_track where usr = ? order by workout_date asc", user_id) // NOTE: add ordering options
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

func (Db *DataBase) ReadUsersRecentlyUsedTracks(user_id int) ([]*WorkoutTrack, error) {

	sql_query := `
	select workout_plan.id, workout_plan.name, workout_plan.creator workout_plan.description, max(workout_date)
	from workout_track inner join workout_plan on plan = workout_plan.id
	group by plan
	where usr = ?
	order by max(workout_date) asc
	`

	rows, err := Db.Data.Query(sql_query, user_id)
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

	statement := "insert into workout_track_data (track, ex_day, weight, set_num, rep_num) values (?, ?, ?, ?, ?) returning id"
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
		td.RepNum,
	).Scan(&track_data_id)

	if err != nil {
		return 0, err
	}

	return track_data_id, nil
}

func (Db *DataBase) UpdateTrackData(td *TrackData) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE workout_track_data SET weight = ?, rep_num = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(td.Weight, td.RepNum)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
