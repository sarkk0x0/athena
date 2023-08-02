package main

import (
	"net/http"
	"project.lemfi.net/internal"
	"project.lemfi.net/internal/data"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	// get next id
	// persist user
	user := &data.User{
		ID:                 app.store.GetNextID(),
		Name:               internal.GenerateRandomName(),
		Balance:            1000,
		VerificationStatus: false,
	}
	err := app.store.CreateUser(user)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	app.verificationQueue <- user
	err = app.writeJSON(w, http.StatusCreated, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.store.GetUsers()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, users, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
