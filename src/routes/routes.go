package routes

import (
	"encoding/json"
	"fmt"
	"goolang-with-docker/src/db"
	"goolang-with-docker/src/helper"
	"io/ioutil"
	"net/http"
)

type user struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Entrou aqui")
	bodyRequest, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.Write(helper.RespMessageError("Failed to read body request"))
		return
	}

	var user user

	if erro = json.Unmarshal(bodyRequest, &user); erro != nil {
		w.Write(helper.RespMessageError("Something wrong happened"))
		return
	}

	// fmt.Println(user)

	db, erro := db.Connectar()

	if erro != nil {
		w.Write(helper.RespMessageError("Failed to connecting with DB"))
		return
	}

	defer db.Close()

	// INSERT INTO usuarios (nome, email) values ("nome", "email")
	// PREPARE STATEMENT
	// ELE É UTILIZADO PARA EVITAR O SQL INJECTION

	statement, erro := db.Prepare("INSERT INTO usuarios (nome, email) values (?, ?)")

	if erro != nil {
		w.Write(helper.RespMessageError("Error creating statement in DB"))
		return
	}

	defer statement.Close()

	insercao, erro := statement.Exec(user.Nome, user.Email)

	if erro != nil {
		w.Write(helper.RespMessageError("Error to insert data in DB"))
		return
	}

	// Se chegar aqui, usuário foi inserido
	// Vou retornar o ID que foi inserido

	idInsercao, erro := insercao.LastInsertId()

	if erro != nil {
		w.Write(helper.RespMessageError("Error to get ID insert in DB"))
		return
	}

	// json.NewEncoder(w).Encode(map[string]string{"status": "OK"})

	//STATUS CODE: 201 CREATED, 404 NOT FOUND, 204 NOT CONTENT
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseJson, erro := helper.ResponseJSON(fmt.Sprintf("User created with success! ID: %d", idInsercao))

	if erro != nil {
		w.Write(helper.RespMessageError("Something wrong happened"))
		return
	}

	w.Write(responseJson)
}

func FindUsers(w http.ResponseWriter, r *http.Request) {
	db, erro := db.Connectar()

	if erro != nil {
		w.Write(helper.RespMessageError("Wrong to connect with database"))
		return
	}

	defer db.Close()

	//SELECT * FROM usuarios

	rows, erro := db.Query("SELECT * FROM usuarios")

	if erro != nil {
		w.Write(helper.RespMessageError("Failed to tried find users"))
		return
	}

	defer rows.Close()

	// create slice for users
	var users []user

	// how i am get all users for query
	// this rows return me several rows so i go through that slice of users below
	for rows.Next() {
		var u user

		// ROW = 1 JOÃO EMAIL
		// I GO SCANNING EACH ONE THESES ROWS
		// AND GO THROW THESES ROWS IN INSIDE SLICE THE USERS

		if erro := rows.Scan(&u.ID, &u.Nome, &u.Email); erro != nil {
			w.Write(helper.RespMessageError("Failed to tried find users"))
			return
		}

		// insert user in slice users
		users = append(users, u)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write(helper.RespMessageError("Failed when scanning user"))
		return
	}
}

func FindUser(w http.ResponseWriter, r *http.Request) {

}
