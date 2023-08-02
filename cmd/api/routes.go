package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(writer, "Hello")
	})
	router.HandleFunc("/v1/users", app.createUser).Methods("POST")
	router.HandleFunc("/v1/users", app.getUsers).Methods("GET")
	router.HandleFunc("/v1/transactions", app.createTransaction).Methods("POST")
	return router
}
