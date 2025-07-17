package models

import "errors"

var cachedExercises map[int]*Exercise
var exercisesByName map[string]int
var cachedTargets map[int]*Target
var cacehdPlansBasic map[int]*WorkoutPlan // NOTE: Only the basic data for a plan is stored: id, name, creator and description
var cachedGyms map[int]*Gym

func init() {
	cachedExercises = make(map[int]*Exercise)
	exercisesByName = make(map[string]int)
	cachedTargets = make(map[int]*Target)
	cacehdPlansBasic = make(map[int]*WorkoutPlan)
	cachedGyms = make(map[int]*Gym)
}

func FetchCachedExercise(ex_id int) (*Exercise, bool) {
	ex, ok := cachedExercises[ex_id]

	if ok {
		return ex, true
	}
	return &Exercise{}, false
}

func CacheExercise(ex *Exercise) error {
	if ex.Id <= 0 {
		return errors.New("Cannot cache exercise without an id")
	}
	cachedExercises[ex.Id] = ex
	exercisesByName[ex.Name] = ex.Id
	return nil
}

func GetAllCachedExercises() []*Exercise {
	res := make([]*Exercise, 0)
	for _, ex := range cachedExercises {
		res = append(res, ex)
	}

	return res
}

func AddTargetToExercise(tar, ex int) error {
	if _, ok := cachedExercises[ex]; !ok {
		return errors.New("Cannot add target to non cached exercise")
	}
	if _, ok := cachedTargets[tar]; !ok {
		return errors.New("Cannot add non cached target to exercise")
	}

	cachedExercises[ex].Targets = append(cachedExercises[ex].Targets, cachedTargets[tar])
	return nil

}

func CacheTarget(t *Target) error {
	if t.Id <= 0 {
		return errors.New("Cannot cache target without an id")
	}
	cachedTargets[t.Id] = t
	return nil
}

func GetAllCachedTargets() []*Target {
	res := make([]*Target, 0)
	for _, tar := range cachedTargets {
		res = append(res, tar)
	}

	return res
}

func AddExerciseToTarget(tar, ex int) error {
	if _, ok := cachedExercises[ex]; !ok {
		return errors.New("Cannot add target to non cached exercise")
	}
	if _, ok := cachedTargets[tar]; !ok {
		return errors.New("Cannot add non cached target to exercise")
	}

	cachedTargets[tar].Exercises = append(cachedTargets[tar].Exercises, cachedExercises[ex])
	return nil

}

func FetchCachedPlanBasic(wp_id int) (*WorkoutPlan, bool) {
	wp, ok := cacehdPlansBasic[wp_id]

	if ok {
		return wp, true
	}
	return &WorkoutPlan{}, false
}

func CachePlanBasic(wp *WorkoutPlan) error {
	if wp.Id <= 0 {
		return errors.New("Cannot cache workout plan without an id")
	}
	cacehdPlansBasic[wp.Id] = wp
	return nil
}

func (Db *DataBase) CacheAllExercises() error {

	rows, err := Db.Data.Query("select id, name, description, exercise_type, difficulty from exercises")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		current_exercise := Exercise{}

		err = rows.Scan(
			&current_exercise.Id,
			&current_exercise.Name,
			&current_exercise.Description,
			&current_exercise.ExerciseType,
			&current_exercise.Difficulty,
		)
		if err != nil {
			return err
		}

		err := CacheExercise(&current_exercise)
		if err != nil {
			return err
		}
	}

	return nil
}

func (Db *DataBase) CacheAllTargets() error {

	rows, err := Db.Data.Query("select id, standard_name, latin_name from target")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		current_target := Target{}

		err = rows.Scan(
			&current_target.Id,
			&current_target.StandardName,
			&current_target.LatinName,
		)
		if err != nil {
			return err
		}

		err := CacheTarget(&current_target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (Db *DataBase) LinkCachedExercisesAndTargets() error {

	rows, err := Db.Data.Query("select exercise, target from exercise_targets")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var ex int
		var tar int

		err = rows.Scan(
			&ex,
			&tar,
		)
		if err != nil {
			return err
		}

		_, ok := cachedExercises[ex]
		if !ok {
			return errors.New("Couldn't link exercise and target: missing exercise")
		}
		_, ok = cachedTargets[tar]
		if !ok {
			return errors.New("Couldn't link exercise and target: missing target")
		}

		AddTargetToExercise(tar, ex)
		AddExerciseToTarget(ex, tar)
	}

	return nil
}

func (Db *DataBase) CacheAllPlansBasic() error {

	rows, err := Db.Data.Query("select id, name, creator, description from workout_plan")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var wp WorkoutPlan

		err = rows.Scan(&wp.Id, &wp.Name, &wp.Creator, &wp.Description)
		if err != nil {
			return err
		}

		err = CachePlanBasic(&wp)
		if err != nil {
			return err
		}
	}

	return nil
}

func CacheGym(g *Gym) error {
	if g.Id <= 0 {
		return errors.New("Cannot cache gym without an id")
	}

	cachedGyms[g.Id] = g
	return nil
}

func (Db *DataBase) CacheAllGyms() error {
	rows, err := Db.Data.Query("select id, name, location from gym")
	if err != nil {
		return err
	}

	for rows.Next() {
		var g Gym
		err = rows.Scan(&g.Id, &g.Name, &g.Location)
		if err != nil {
			return err
		}

		err = CacheGym(&g)
		if err != nil {
			return err
		}
	}

	return nil
}
