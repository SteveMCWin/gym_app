package models

import (
	"cmp"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"log"
	"slices"

	_ "github.com/mattn/go-sqlite3"
)

type Gym struct {
	Id            int          `json:"id"`
	Name          string       `json:"name"`
	Location      string       `json:"location"`
	NumberOfUsers int          `json:"number_of_users"`
	Equipment     []*Equipment `json:"equipment"`
}

type Equipment struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (Db *DataBase) ReadGym(gym_id int) (*Gym, error) {
	row := Db.Data.QueryRow("select name, location from gym where id = ?", gym_id)
	g := &Gym{ Id: gym_id }
	err := row.Scan(&g.Name, &g.Location)
	if err != nil {
		return nil, err
	}

	g.Equipment, err = Db.ReadGymEquipment(g.Id)
	if err != nil {
		return nil, err
	}

	row = Db.Data.QueryRow("select count(*) from gym_users where gym_id = ?", gym_id)
	err = row.Scan(&g.NumberOfUsers)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (Db *DataBase) ReadAllGyms() ([]*Gym, error) {
	rows, err := Db.Data.Query("select id, name, location from gym")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	gyms := make([]*Gym, 0)

	for rows.Next() {
		g := &Gym{}
		err = rows.Scan(&g.Id, &g.Name, &g.Location)
		if err != nil {
			return nil, err
		}

		g.Equipment, err = Db.ReadGymEquipment(g.Id)
		if err != nil {
			return nil, err
		}

		row := Db.Data.QueryRow("select count(*) from gym_users where gym_id = ?", g.Id)
		err = row.Scan(&g.NumberOfUsers)
		if err != nil {
			return nil, err
		}

		gyms = append(gyms, g)
	}

	return gyms, nil
}

func (Db *DataBase) ReadGymEquipment(gym_id int) ([]*Equipment, error) {
	rows, err := Db.Data.Query("select id, name from equipment inner join gym_equipment on gym_equipment.equipment = equipment.id where gym_id = ?", gym_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func (Db *DataBase) CheckIfGymHasPlanEquipment(gym_id, plan_id int) (map[int]int, error) {
	query := `
	select exercise_day.exercise
	from exercise_day 
	where plan = ? and exercise_day.exercise not in (
	select exercise_equipment.exercise 
	from exercise_equipment inner join gym_equipment on gym_equipment.equipment = exercise_equipment.equipment 
	where gym_equipment.gym_id = ?);
	`
	rows, err := Db.Data.Query(query, plan_id, gym_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make(map[int]int)

	for rows.Next() {
		var ex_id int
		err = rows.Scan(&ex_id)
		if err != nil {
			return nil, err
		}
		res[ex_id] = ex_id
	}

	return res, nil
}

// this one is just for conviniece
func (Db *DataBase) GetPlanGymExDiff(gym_id, wp_id int) ([]Exercise, error) {

	ex_no_eq := make([]Exercise, 0)

	undoable_ex, err := Db.CheckIfGymHasPlanEquipment(gym_id, wp_id)
	if err != nil {
		return nil, err
	}

	for ex := range undoable_ex {
		cached_ex, ok := FetchCachedExercise(ex)
		if ok {
			ex_no_eq = append(ex_no_eq, *cached_ex)
		} else {
			log.Println("No cached exercise with id", ex)
		}
	}

	return ex_no_eq, nil
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

	defer stmt_gym.Close()

	err = stmt_gym.QueryRow(g.Name, g.Location).Scan(&g.Id)
	if err != nil {
		return err
	}

	statement = "insert into gym_equipment (gym_id, equipment) values (?, ?)"
	stmt_eq, err := tx.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt_eq.Close()

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

func SearchForGym(name string) []*Gym {
	all_gyms := FetchAllCachedGyms()
	gym_names := make([]string, len(all_gyms))
	for i, g := range all_gyms {
		gym_names[i] = g.Name
	}

	ranks := fuzzy.RankFindFold(name, gym_names)
	slices.SortFunc(ranks, func(a, b fuzzy.Rank) int {
		return cmp.Compare(a.Distance, b.Distance)
	})

	res := make([]*Gym, 0)
	for _, rank := range ranks {
		res = append(res, all_gyms[rank.OriginalIndex])
	}

	return res
}

func (Db *DataBase) AddUserToGym(gym_id, user_id int) error {
	statement := "insert into gym_users (gym_id, user_id) values (?, ?)"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(gym_id, user_id)
	gym, ok := FetchCachedGym(gym_id)
	if !ok {
		gym, err = Db.ReadGym(gym_id)
		if err != nil {
			return err
		}
	} else {
		gym.NumberOfUsers += 1
	}

	return err
}

func (Db *DataBase) RemoveUserFromAllGyms(user_id int) error {
	rows, err := Db.Data.Query("select gym_id from gym_users where user_id = ?", user_id)
	if err != nil {
		log.Println("The rows query is the problem")
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var g_id int
		err = rows.Scan(&g_id)
		if err != nil {
			return err
		}

		g, ok := FetchCachedGym(g_id)
		if ok {
			g.NumberOfUsers -= 1
		} else {
			g, err = Db.ReadGym(g_id)
			if err != nil {
				return err
			}
			err = CacheGym(g)
		}
	}

	statement := "delete from gym_users where user_id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		log.Println("The statement is the problem")
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user_id)
	return err
}

func (Db *DataBase) UpdateUsersGyms(user_id, old_gym_id, new_gym_id int) error {
	statement := "update gym_users set gym_id = ? where gym_id = ? and user_id = ?"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(new_gym_id, old_gym_id, user_id)

	old_g, ok := FetchCachedGym(old_gym_id)
	if ok {
		old_g.NumberOfUsers -= 1
	} else {
		old_g, err = Db.ReadGym(old_gym_id)
		if err != nil {
			return err
		}
		err = CacheGym(old_g)
		// cache the gym
	}

	new_g, ok := FetchCachedGym(new_gym_id)
	if ok {
		new_g.NumberOfUsers += 1
	} else {
		new_g, err = Db.ReadGym(new_gym_id)
		if err != nil {
			return err
		}
		err = CacheGym(new_g)
		// cache the gym
	}

	return err
}
