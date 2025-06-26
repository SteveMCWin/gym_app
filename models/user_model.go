package models

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	TrainingSince float32 `json:"training_since"`
	IsTrainer     bool    `json:"is_trainer"`
	GymGoals      string  `json:"gym_goals"`
	CurrentGym    string  `json:"current_gym"` // perhaps change this to be an id of a gym in the database
}

type DataBase struct {
	Data      *sql.DB
	is_opened bool
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
			return 0, err
		}

		defer stmt.Close()

		var usr_id int

		encrypted_pass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
		if err != nil {
			return 0, err
		}

		err = stmt.QueryRow(usr.Name, usr.Email, string(encrypted_pass), usr.TrainingSince, usr.IsTrainer, usr.GymGoals, usr.CurrentGym).Scan(&usr_id)
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
	) // gets the public data of the user

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (Db *DataBase) AuthUser(email, password string) (int, error) { // returns the id if the credentials are right, 0 if not (as well as an error)
	if email == "" || password == "" {
		return 0, errors.New("Empty email and password provided")
	}

	var usr_id int
	var stored_password string

	err := Db.Data.QueryRow("select id, password from users where email like ?", email).Scan(&usr_id, &stored_password)

	if err != nil { // not account with provided email
		return 0, errors.New("ERROR: no account with provided email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(stored_password), []byte(password))

	if err != nil {
		return 0, errors.New("ERROR: wrong password")
	}

	return usr_id, nil
}

func (Db *DataBase) EmailExists(email string) bool {
	var tmp int
	err := Db.Data.QueryRow("select id form users where email like ?", email).Scan(&tmp)

	return err == nil
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

func (Db *DataBase) UpdateUserPassword(usr_id int, pass string) (bool, error) { // before this, should send email from which you change your password
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE users SET password = ? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	encrypted_pass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(string(encrypted_pass), usr_id)

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
