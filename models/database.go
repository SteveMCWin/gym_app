package models

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

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


