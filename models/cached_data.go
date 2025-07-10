package models

import "errors"

var cachedExercises map[int]*Exercise
var exercisesByName map[string]int
var cachedTargets map[int]*Target

func init() {
	cachedExercises = make(map[int]*Exercise)
	exercisesByName = make(map[string]int)
	cachedTargets = make(map[int]*Target)
}

func FetchCachedExercise(ex_id int) (*Exercise, bool) {
	ex, ok := cachedExercises[ex_id]

	if ok {
		return ex, true
	}
	return &Exercise{}, false
}

func CacheExercise(ex *Exercise) error {
	if ex.Id == 0 {
		return errors.New("Cannot cache exercise without an id")
	}
	cachedExercises[ex.Id] = ex
	exercisesByName[ex.Name] = ex.Id
	return nil
}

func CacheTarget(t *Target) error {
	if t.Id == 0 {
		return errors.New("Cannot cache target without an id")
	}
	cachedTargets[t.Id] = t
	return nil
}

func GetAllCachedExercises() []*Exercise {
	res := make([]*Exercise, 0)
	for _, ex := range cachedExercises {
		res = append(res, ex)
	}

	return res
}

func GetAllCachedTargets() []*Target {
	res := make([]*Target, 0)
	for _, tar := range cachedTargets {
		res = append(res, tar)
	}

	return res
}
