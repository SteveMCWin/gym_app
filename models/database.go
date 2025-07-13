package models

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	Data      *sql.DB
	is_open bool
}

func (dataBase *DataBase) Close() {
	dataBase.Data.Close()
	dataBase.is_open = false
}

func (dataBase *DataBase) InitDatabase() error {
	if dataBase.is_open {
		return errors.New("ERROR: Database already open")
	}
	var err error
	dataBase.Data, err = sql.Open("sqlite3", "models/database.db")
	if err != nil {
		return err
	}

	_, err = dataBase.Data.Exec(`PRAGMA enable_load_extension = 1;`)
	if err != nil {
		return err
	}

	_, err = dataBase.Data.Exec(`SELECT load_extension('./spellfix.so')`)
    if err != nil {
        return err
    }

	dataBase.is_open = true

	return nil
}


