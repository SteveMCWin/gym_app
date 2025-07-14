package models

import (
	"os"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
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

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	sql.Register("sqlite3_with_extension",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				dir+"/models/spellfix.so",
			},
		},
	)

	dataBase.Data, err = sql.Open("sqlite3_with_extension", "models/database.db")
	if err != nil {
		return err
	}

	dataBase.is_open = true

	return nil
}


