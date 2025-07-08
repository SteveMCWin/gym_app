package models

import (
)

var cachedExercises map[int]Exercise
var exercisesByName map[string]int

func init() {
	cachedExercises = make(map[int]Exercise)
	exercisesByName = make(map[string]int)
}

func FetchCachedExercise(ex_id int) (Exercise, bool) {
	ex, ok := cachedExercises[ex_id]

	if ok {
		return ex, true
	}
	return Exercise{}, false
}

func CacheExercise(ex Exercise) bool {
	if ex.Id == 0 {
		return false
	}
	cachedExercises[ex.Id] = ex
	exercisesByName[ex.Name] = ex.Id
	return true
}

func GetAllCachedExercises() []*Exercise {
	res := make([]*Exercise, 0)
	for _, ex := range cachedExercises {
		res = append(res, &ex)
	}

	return res
}
