package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	TrainingSince time.Time `json:"training_since"`
	IsTrainer     bool      `json:"is_trainer"`
	GymGoals      string    `json:"gym_goals"`
	CurrentGym    Gym       `json:"current_gym"`
	CurrentPlan   int       `json:"current_plan"`
	DateCreated   time.Time `json:"time_created"`
}

func (Db *DataBase) CreateUser(c *gin.Context, usr User) (int, error) {

	log.Println("CREATING USERRRRR")

	if usr.Email == "" {
		return 0, errors.New("Cannot store a user without their email")
	}

	err := Db.Data.QueryRow("select id from users where email like ?", usr.Email).Scan(&usr.Id)

	if err != nil {
		// user is signing up
		statement := "insert into users (name, email, password, training_since, is_trainer, gym_goals, current_gym, current_plan, date_created) values (?, ?, ?, ?, ?, ?, ?, ?, ?) returning id"
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

		err = stmt.QueryRow(
			usr.Name,
			usr.Email,
			string(encrypted_pass),
			usr.TrainingSince.Format("2006-01-02"),
			usr.IsTrainer,
			usr.GymGoals,
			usr.CurrentGym.Id,
			1, // WARNING: I am assuming there will be a default plan number 1 that is used only for indicating there is no plan
			time.Now().Format("2006-01-02")).Scan(&usr_id)

		if err != nil {
			return 0, err
		}

		return usr_id, nil
	}

	// user already has an account
	return 0, errors.New("ERROR: user already has an account") // NOTE: handle this better, lead the user to the login page
}

func (Db *DataBase) ReadUser(usr_id int) (*User, error) {
	usr := &User{}

	var curr_gym_id int
	err := Db.Data.QueryRow("select id, name, email, training_since, is_trainer, gym_goals, current_gym, current_plan, date_created from users where id = ?", usr_id).Scan(
		&usr.Id,
		&usr.Name,
		&usr.Email,
		&usr.TrainingSince,
		&usr.IsTrainer,
		&usr.GymGoals,
		&curr_gym_id,
		&usr.CurrentPlan,
		&usr.DateCreated,
	) // gets the public data of the user

	if err != nil {
		return nil, err
	}

	curr_gym, ok := FetchCachedGym(curr_gym_id)
	if !ok {
		log.Println("No gym with id:", curr_gym_id)
		return nil, errors.New("Invalid gym id")
	}

	usr.CurrentGym = *curr_gym

	return usr, nil
}

func (Db *DataBase) ReadUserCurrentGymId(usr_id int) (int, error) {
	var curr_gym_id int
	err := Db.Data.QueryRow("select current_gym from users where id = ?", usr_id).Scan(
		&curr_gym_id,
	) // gets the public data of the user

	if err != nil {
		return 0, err
	}

	return curr_gym_id, nil
}

func (Db *DataBase) ReadUserShallow(usr_id int) (*User, error) {
	usr := &User{}

	err := Db.Data.QueryRow("select id, name, training_since, is_trainer, current_plan from users where id = ?", usr_id).Scan(
		&usr.Id,
		&usr.Name,
		&usr.TrainingSince,
		&usr.IsTrainer,
		&usr.CurrentPlan,
	) // gets most of the public data of the user

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (Db *DataBase) ReadUserIdByEmail(email string) (int, error) {
	var usr_id int

	err := Db.Data.QueryRow("select id from users where email like ?", email).Scan(
		&usr_id,
	)

	if err != nil {
		return 0, err
	}

	return usr_id, nil
}

func (Db *DataBase) AuthUserByEmail(email, password string) (int, error) { // returns the id if the credentials are right, 0 if not (as well as an error)
	if email == "" || password == "" {
		return 0, errors.New("Empty email or password provided")
	}

	var usr_id int
	var stored_password string

	err := Db.Data.QueryRow("select id, password from users where email like ?", email).Scan(&usr_id, &stored_password)

	if err != nil {
		return 0, errors.New("ERROR: no account with provided email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(stored_password), []byte(password))

	if err != nil {
		return 0, errors.New("ERROR: wrong password")
	}

	return usr_id, nil
}

func (Db *DataBase) AuthUserByID(usr_id int, password string) error { // returns the id if the credentials are right, 0 if not (as well as an error)
	if usr_id == 0 || password == "" {
		return errors.New("Invalid ID or password provided")
	}

	var stored_password string

	err := Db.Data.QueryRow("select password from users where id = ?", usr_id).Scan(&stored_password)

	if err != nil { // not account with provided email
		return errors.New("ERROR: no account with provided id")
	}

	err = bcrypt.CompareHashAndPassword([]byte(stored_password), []byte(password))

	if err != nil {
		return errors.New("ERROR: wrong password")
	}

	return nil
}

func (Db *DataBase) EmailExists(email string) bool {
	var tmp int
	err := Db.Data.QueryRow("select id from users where email = ?", email).Scan(&tmp)
	return err == nil
}

func (Db *DataBase) UpdateUserPublicData(usr *User) (bool, error) {
	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE users SET name = ?, training_since = ?, is_trainer = ?, gym_goals = ?, current_gym = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(usr.Name, usr.TrainingSince, usr.IsTrainer, usr.GymGoals, usr.CurrentGym.Id, usr.Id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (Db *DataBase) UpdateUserCurrentPlan(usr_id, plan_id int) (bool, error) {

	if !Db.CheckIfUserUsesPlan(usr_id, plan_id) {
		return false, errors.New("Cannot make plan current if user doesn't even use it")
	}

	tx, err := Db.Data.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE users SET current_plan = ? WHERE Id = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(plan_id, usr_id)

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

	defer tx.Rollback()

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

	defer tx.Rollback()

	stmt_usr, err := tx.Prepare("DELETE from users where id = ?")
	if err != nil {
		return false, err
	}

	defer stmt_usr.Close()

	_, err = stmt_usr.Exec(id)

	if err != nil {
		return false, err
	}

	_, err = Db.DeleteAllWorkoutsForUser(id)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (Db *DataBase) SearchForUsers(username string, requesting_user_id int) ([]User, error) {
	rows, err := Db.Data.Query("select id, name from spellfix_users inner join users on word = name where word match ? and id != ?", username, requesting_user_id)
	if err != nil {
		return nil, err
	}

	matches := make([]User, 0)

	for rows.Next() {
		var usr User
		err = rows.Scan(&usr.Id, &usr.Name)
		if err != nil {
			return nil, err
		}

		matches = append(matches, usr)
	}

	return matches, nil
}
