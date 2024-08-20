package router

import (
	"github.com/TeslaMode1X/gormTest/connection"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/users", connection.GetAllUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}", connection.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/api/user/{id}", connection.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/create/user", connection.CreateUser).Methods("POST")

	return r
}
