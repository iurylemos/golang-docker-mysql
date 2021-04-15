package main

import (
	"fmt"
	"goolang-with-docker/src/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// CRUD - "C"REATE (POST), "R"EAD (GET), "U"PDATE (UPDATE), "D"ELETE (DELETE)

	router := mux.NewRouter()

	router.HandleFunc("/usuarios", routes.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", routes.FindUsers).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", routes.FindUser).Methods(http.MethodGet)
	router.HandleFunc("/usuarios/{id}", routes.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/usuarios/{id}", routes.UpdateUser).Methods(http.MethodDelete)

	fmt.Printf("Listening in port %s", "5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
