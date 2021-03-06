package routes

import (
	"encoding/json"
	"fmt"
	"goolang-with-docker/src/db"
	"goolang-with-docker/src/helper"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(http.StatusCreated)

	responseJson, erro := helper.ResponseJSON(fmt.Sprintf("User created with success! ID: %d", idInsercao))

	if erro != nil {
		w.Write(helper.RespMessageError("Something wrong happened"))
		return
	}

	w.Write(responseJson)
}

func FindUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write(helper.RespMessageError("Failed when scanning user"))
		return
	}
}

func FindUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)

	// Convert o ID (string) for (int)
	// this function received 3 parameteres
	// value, base => base 10 (ten) why this is int, and 32 (thirteen-two) length of bits or 64 (sixteen-four)
	ID, erro := strconv.ParseUint(parameters["id"], 10, 32)
	if erro != nil {
		w.Write(helper.RespMessageError("Error to convert parameter. Value is not valid"))
		return
	}

	db, erro := db.Connectar()

	if erro != nil {
		w.Write(helper.RespMessageError("Wrong to connect with database"))
		return
	}

	row, erro := db.Query("SELECT * FROM usuarios WHERE id = ?", ID)

	if erro != nil {
		w.Write(helper.RespMessageError("Wrong to connect with database"))
		return
	}

	var u user
	if row.Next() {
		if erro := row.Scan(&u.ID, &u.Nome, &u.Email); erro != nil {
			w.Write(helper.RespMessageError("Failed to scanning user inside the db"))
			return
		}
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		w.Write(helper.RespMessageError("Failed when scanning user"))
		return
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)

	ID, erro := strconv.ParseUint(parameters["id"], 10, 32)
	if erro != nil {
		w.Write(helper.RespMessageError("Error to convert parameter. Value is not valid"))
		return
	}

	bodyRequest, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.Write(helper.RespMessageError("Error to read body request. Please, try again"))
		return
	}

	var user user

	if erro := json.Unmarshal(bodyRequest, &user); erro != nil {
		w.Write(helper.RespMessageError("Error to convert user sent of body request"))
		return
	}

	//open connected with db after that read body request
	db, erro := db.Connectar()

	if erro != nil {
		w.Write(helper.RespMessageError("Wrong to connect with database"))
		return
	}

	//insert, delete, update theses cases is for use statement
	defer db.Close()

	statement, erro := db.Prepare("UPDATE usuarios SET nome = ?, email = ? where id = ?")
	if erro != nil {
		w.Write(helper.RespMessageError("Error to created statement"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(user.Nome, user.Email, ID); erro != nil {
		w.Write(helper.RespMessageError("Error to update user in database"))
		return
	}

	user.ID = uint32(ID)

	// w.WriteHeader(http.StatusNoContent)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write(helper.RespMessageError("Something went wrong. Please try again after"))
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)

	ID, erro := strconv.ParseUint(parameters["id"], 10, 32)
	if erro != nil {
		w.Write(helper.RespMessageError("Error to convert parameter. Value is not valid"))
		return
	}

	db, erro := db.Connectar()

	if erro != nil {
		w.Write(helper.RespMessageError("Wrong to connect with database"))
		return
	}

	//insert, delete, update theses cases is for use statement
	defer db.Close()

	statement, erro := db.Prepare("DELETE FROM usuarios WHERE id = ?")
	if erro != nil {
		w.Write(helper.RespMessageError("Error to created statement"))
		return
	}

	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write(helper.RespMessageError("Error to delete user in database"))
		return
	}

	w.WriteHeader(http.StatusOK)

	responseJson, erro := helper.ResponseJSON(fmt.Sprintf("User delete with success! ID: %d", ID))

	if erro != nil {
		w.Write(helper.RespMessageError("Something wrong happened"))
		return
	}

	w.Write(responseJson)
}
