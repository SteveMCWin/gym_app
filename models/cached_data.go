package models

import "errors"

var cachedExercises map[int]*Exercise
var exercisesByName map[string]int
var cachedTargets map[int]*Target
var cacehdPlansBasic map[int]*WorkoutPlan // NOTE: Only the basic data for a plan is stored: id, name, creator and description

func init() {
	cachedExercises = make(map[int]*Exercise)
	exercisesByName = make(map[string]int)
	cachedTargets = make(map[int]*Target)
	cacehdPlansBasic = make(map[int]*WorkoutPlan)
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
	if t.Id == 0 {
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
	if wp.Id == 0 {
		return errors.New("Cannot cache workout plan without an id")
	}
	cacehdPlansBasic[wp.Id] = wp
	return nil
}

