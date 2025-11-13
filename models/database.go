package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

type DataBase struct {
	Data    *sql.DB
	is_open bool
}

func (dataBase *DataBase) Close() {
	dataBase.Data.Close()
	dataBase.is_open = false
}

// initializes the database
// if any bool parameters (whether true or false) are passed, uses the test_database
func (Db *DataBase) InitDatabase(is_test ...bool) error {
	if Db.is_open {
		return errors.New("ERROR: Database already open")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	spellfix_path := filepath.Join(dir, "extensions", "spellfix.so")

	sql.Register("sqlite3_with_extension",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				spellfix_path,
			},
		},
	)

	data_dir := "data"
	if err := os.MkdirAll(data_dir, 0755); err != nil {
		return err
	}

	db_name := "database.db"

	db_path := filepath.Join(data_dir, db_name)
	// db_path := "models/"
	if len(is_test) != 0 {
		db_name = "test_database.db"
	}
	// if len(is_test) == 0 {
	// 	db_path = db_path + "database.db"
	// } else {
	// 	db_path = db_path + "test_database.db"
	// }

	_, err = os.Stat(db_path)
	dbExists := !os.IsNotExist(err)

	Db.Data, err = sql.Open("sqlite3_with_extension", db_path)
	if err != nil {
		return err
	}

	if dbExists != true {
		sqlFiles := []string{
			"data/create_exercise_table.sql",
			"data/create_user_table.sql",
			"data/create_plan_table.sql",
			"data/create_gym_table.sql",
			"data/create_track_table.sql",
			"data/create_sessions_table.sql",
			"data/create_spellfix_users.sql",
			"data/populate_exercise_table.sql",
		}

		for _, sqlFile := range sqlFiles {
			log.Printf("Running %s...\n", sqlFile)

			// Open the SQL file
			file, err := os.Open(sqlFile)
			if err != nil {
				log.Println("failed to open sql file:", sqlFile, err)
				return err
			}

			// Execute: sqlite3 database.db < sqlFile
			cmd := exec.Command("sqlite3", db_path)
			cmd.Stdin = file
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				file.Close()
				log.Println("failed to execute sql file:", sqlFile, err)
				return err
			}

			file.Close()
		}
	}

	Db.is_open = true

	return nil
}

// stores some stuff from the database to runtime memory for convenience and speed
func (Db *DataBase) CacheData() error {
	if !Db.is_open {
		return errors.New("Cannot cache data from closed database")
	}

	err := Db.CacheAllExercises()
	if err != nil {
		return err
	}

	err = Db.CacheAllTargets()
	if err != nil {
		return err
	}

	err = Db.LinkCachedExercisesAndTargets()
	if err != nil {
		return err
	}

	err = Db.CacheAllPlansBasic()
	if err != nil {
		return err
	}

	err = Db.CacheAllGyms()
	if err != nil {
		return err
	}

	return nil

}
