package models

import (
	// "github.com/gin-gonic/gin"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id                  int
	Name                string
	Email               string
	Password            string
	YearsOfTraining     float32
	IsTrainer           bool
	GymGoals            string
	CurrentGym          string  // perhaps change this to be an id of a gym in the database
}

type DataBase struct {
	data        *sql.DB
	is_opened   bool
}

func (dataBase *DataBase) InitDatabase() {
	var err error
	dataBase.data, err = sql.Open("sqlite3", "users/users.db")
	if err != nil {
		panic(err)
	}

	dataBase.is_opened = true
}

func (Db *DataBase) CreateUser(usr User) (int, error) {
	if usr.Email == "" {
		return 0, errors.New("Cannot store a user without their email")
	}

	err := Db.data.QueryRow("select id from users where email like ?", usr.Email).Scan(&usr.Id)

	if err != nil {
		// user is signing up
        statement := "insert into users (name, email, password, years_of_training, is_trainer, gym_goals, current_gym) values (?, ?, ?, ?, ?, ?, ?, ?) returning id"
		var stmt *sql.Stmt
		stmt, err = Db.data.Prepare(statement)
		if err != nil {
			return 0, errors.New("ERROR: Couldn't prepare statement for storing user")
		}

		defer stmt.Close()
		var usr_id int
		err = stmt.QueryRow(usr.Name, usr.Email, usr.Password, usr.YearsOfTraining, usr.IsTrainer, usr.GymGoals, usr.CurrentGym).Scan(&usr_id)
		if err != nil {
			return 0, err
		}

		return usr_id, nil
	}

	// user already has an account
	return 0, errors.New("ERROR: user already has an account")
}
