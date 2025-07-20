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

func (Db *DataBase) CheckIfGymHasPlanEquipment(gym_id, plan_id int) ([]int, bool, error) {
	query := `
	select exercise_day.id
	from exercise_day 
	where plan = ? and exercise_day.exercise not in (
	select exercise_equipment.exercise 
	from exercise_equipment inner join gym_equipment on gym_equipment.equipment = exercise_equipment.equipment 
	where gym_equipment.gym_id = ?);
	`
	rows, err := Db.Data.Query(query, plan_id, gym_id)
	if err != nil {
		return nil, false, err
	}

	res := make([]int, 0)

	for rows.Next() {
		var ex_day_id int
		err = rows.Scan(&ex_day_id)
		if err != nil {
			return nil, false, err
		}
		res = append(res, ex_day_id)
	}

	return res, len(res) > 0, nil
}

func (Db *DataBase) CreateGym(g *Gym) error {
	tx, err := Db.Data.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement := "insert into gym (name, location) values (?, ?) returning id"
	stmt_gym, err := tx.Prepare(statement)
	if err != nil {
		return err
	}

	err = stmt_gym.QueryRow(g.Name, g.Location).Scan(&g.Id)
	if err != nil {
		return err
	}

	statement = "insert into gym_equipment (gym_id, equipment) values (?, ?)"
	stmt_eq, err := tx.Prepare(statement)
	if err != nil {
		return err
	}

	for _, eq := range g.Equipment {
		_, err = stmt_eq.Exec(g.Id, eq.Id)
		if err != nil {
			return err
		}
	}

	err = CacheGym(g)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
