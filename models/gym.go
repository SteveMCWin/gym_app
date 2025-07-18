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
	Equipment []*Equipment `json:"equipment"`
}

type Equipment struct {
	Id int `json:"id"`
	Name string `json:"name"`
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

		g.Equipment, err = Db.ReadGymEquipment(g.Id)
		if err != nil {
			return nil, err
		}

		gyms = append(gyms, &g)
	}

	return gyms, nil
}

func (Db *DataBase) ReadGymEquipment(gym_id int) ([]*Equipment, error) {
	rows, err := Db.Data.Query("select id, name from equipment inner join gym_equipment on gym_equipment.equipment = equipment.id where gym_id = ?", gym_id)
	if err != nil {
		return nil, err
	}

	res := make([]*Equipment, 0)

	for rows.Next() {
		eq := Equipment{}
		err = rows.Scan(&eq.Id, &eq.Name)
		if err != nil {
			return nil, err
		}

		res = append(res, &eq)
	}

	return res, nil
}
