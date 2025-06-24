package models

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id                  int
	Name                string
	Email               string
	Password            string
	TrainingSince     float32
	IsTrainer           bool
	GymGoals            string
	CurrentGym          string  // perhaps change this to be an id of a gym in the database
}

type DataBase struct {
	Data        *sql.DB
	is_opened   bool
}

func (dataBase *DataBase) Close() {
	dataBase.Data.Close()
	dataBase.is_opened = false
}

func (dataBase *DataBase) InitDatabase() error {
	if dataBase.is_opened {
		return errors.New("ERROR: Database already opened")
	}
	var err error
	dataBase.Data, err = sql.Open("sqlite3", "models/database.db")
	if err != nil {
		return errors.New("ERROR: Couldn't open database")
	}

	dataBase.is_opened = true

	return nil
}

func (Db *DataBase) CreateUser(c *gin.Context, usr User) (int, error) {
	if usr.Email == "" {
		return 0, errors.New("Cannot store a user without their email")
	}

	err := Db.Data.QueryRow("select id from users where email like ?", usr.Email).Scan(&usr.Id)

	if err != nil {
		// user is signing up
        statement := "insert into users (name, email, password, training_since, is_trainer, gym_goals, current_gym) values (?, ?, ?, ?, ?, ?, ?, ?) returning id"
		var stmt *sql.Stmt
		stmt, err = Db.Data.Prepare(statement)
		if err != nil {
			return 0, errors.New("ERROR: Couldn't prepare statement for storing user")
		}

		defer stmt.Close()
		var usr_id int
		err = stmt.QueryRow(usr.Name, usr.Email, usr.Password, usr.TrainingSince, usr.IsTrainer, usr.GymGoals, usr.CurrentGym).Scan(&usr_id)
		if err != nil {
			return 0, err
		}

		return usr_id, nil
	}

	// user already has an account
	return 0, errors.New("ERROR: user already has an account")
}

func (Db *DataBase) ReadUser(id int) (*User, error) {
	usr := &User{}
	err := Db.Data.QueryRow("select id, name, email, training_since, is_trainer, gym_goals, current_gym where id = ?", id).Scan(
        &usr.Id,
		&usr.Name,
		&usr.Email,
		&usr.TrainingSince,
		&usr.IsTrainer,
		&usr.GymGoals,
		&usr.CurrentGym,
    )

    if err != nil {
        return nil, err
    }

	return usr, nil

}

func (Db *DataBase) UpdateUserPublicData(usr *User) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE users SET name = ?, email = ?, training_since = ?, is_trainer = ?, gym_goals = ?, current_gym = ? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(usr.Name, usr.Email, usr.TrainingSince, usr.IsTrainer, usr.GymGoals, usr.CurrentGym, usr.Id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (Db *DataBase) UpdateUserPassword(usr_id int, encrypted_pass string) (bool, error) { // before this, should send email from which you change your password
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE users SET password = ? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(encrypted_pass, usr_id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (Db *DataBase) DeleteUser(id int) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from users where id = ?")
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
