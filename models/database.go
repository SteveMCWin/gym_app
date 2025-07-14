package models

import (
	"log"
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

	sql.Register("sqlite3_with_spellfix",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"spellfix1",
			},
		},
	)

	var err error
	dataBase.Data, err = sql.Open("sqlite3", "models/database.db")
	if err != nil {
		return err
	}

	// rows, err := dataBase.Data.Query(`PRAGMA compile_options;`)
	// if err != nil {
	// 	log.Println("Error:", err)
	// 	return err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var opt string
	// 	rows.Scan(&opt)
	// 	log.Println(opt)
	// }

	_, err = dataBase.Data.Exec("PRAGMA enable_load_extension = 1;")
	if err != nil {
		return err
	}

	_, err = dataBase.Data.Exec("SELECT load_extension('./models/spellfix.so')")
    if err != nil {
		log.Println("HEREEEE???")
        return err
    }

	dataBase.is_open = true

	return nil
}


