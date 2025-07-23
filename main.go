package main

import (
	"net/http"

	"fitness_app/handlers"
	"fitness_app/models"
)

func main() {

	var db models.DataBase
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.CacheAllExercises()
	if err != nil {
		panic(err)
	}

	err = db.CacheAllTargets()
	if err != nil {
		panic(err)
	}

	err = db.LinkCachedExercisesAndTargets()
	if err != nil {
		panic(err)
	}

	err = db.CacheAllPlansBasic()
	if err != nil {
		panic(err)
	}

	err = db.CacheAllGyms()
	if err != nil {
		panic(err)
	}

	router := handlers.SetUpRouter(db)

	http.ListenAndServe(":8080", router)
}
