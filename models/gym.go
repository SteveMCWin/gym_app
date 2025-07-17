package models

import (
	// "database/sql"
	// "errors"

	_ "github.com/mattn/go-sqlite3"
)

type Gym struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Location string `json:"location"`
}

func (Db *DataBase) ReadAllGyms() ([]*Gym, error) {
	rows, err := Db.Data.Query("select id, name, location from gym")
	if err != nil {
		return nil, err
	}

	gyms := make([]*Gym, 0)

	for rows.Next() {
		var g Gym
		err = rows.Scan(&g.Id, &g.Name, &g.Location)
		if err != nil {
			return nil, err
		}

		gyms = append(gyms, &g)
	}

	return gyms, nil
}
