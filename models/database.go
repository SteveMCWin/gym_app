package models

import (
	"database/sql"
	"errors"
	"os"

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
// if there are no parameters passed in the database used is one for testing purposes
func (Db *DataBase) InitDatabase(is_test ...bool) error {
	if Db.is_open {
		return errors.New("ERROR: Database already open")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	spellfix_relative_path := "/models/spellfix.so"

	sql.Register("sqlite3_with_extension",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				dir + spellfix_relative_path,
			},
		},
	)

	db_path := "models/"
	if len(is_test) == 0 {
		db_path = db_path + "database.db"
	} else {
		db_path = db_path + "test_database.db"
	}

	Db.Data, err = sql.Open("sqlite3_with_extension", db_path)
	if err != nil {
		return err
	}

	Db.is_open = true

	return nil
}

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
