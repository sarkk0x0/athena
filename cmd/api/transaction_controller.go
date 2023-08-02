package main

import (
	"encoding/json"
	"net/http"
	"project.lemfi.net/internal/data"
)

func (app *application) createTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := data.Transaction{}
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.transactionQueue <- &transaction
	err = app.writeJSON(w, http.StatusCreated, transaction, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
