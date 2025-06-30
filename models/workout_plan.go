package models

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type WorkoutPlan struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Creator     int    `json:"creator"`
	Description string `json:"description"`
}

func (Db *DataBase) CreateWorkoutPlan(c *gin.Context, wp *WorkoutPlan) (int, error) {

	log.Println("CREATING WORKOUT PLANNN")

	if wp.Creator == 0 {
		return 0, errors.New("You need to be logged in to create a workout plan")
	}

	var tmp int
	err := Db.Data.QueryRow("select id from workout_plan where name like ?", wp.Name).Scan(&tmp)

	if wp.Name == "" {
		wp.Name = "WorkoutPlan"
	}

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

	err = Db.AddWorkoutPlanToUser(wp.Creator, wp.Id)
	if err != nil {
		return 0, err
	}

	return workout_plan_id, nil
}

func (Db *DataBase) AddWorkoutPlanToUser(usr_id, wp_id int) error { // adds the workout to the list of workouts the user has/uses/whatever
	if wp_id == 0 || usr_id == 0 {
		return errors.New("Cannot add workout without wp_id and usr_id")
	}

	statement := "insert into workout_plan (usr, plan) values (?, ?)"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(usr_id, wp_id).Scan()
	if err != nil {
		return err
	}

	return nil
}

func (Db *DataBase) RemoveWorkoutPlanFromUser(usr_id, wp_id int) error {
	tx, err := Db.Data.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE from workout_plan where user = ? AND plan = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(usr_id, wp_id)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (Db *DataBase) ReadWorkoutPlan(id int) (*WorkoutPlan, error) {
	wp := &WorkoutPlan{}

	err := Db.Data.QueryRow("select id, name, creator, description from workout_plan where id = ?", id).Scan(
		&wp.Id,
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

