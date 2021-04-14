package routes

import (
	"encoding/json"
	"fmt"
	"goolang-with-docker/src/db"
	"io/ioutil"
	"net/http"
)

type user struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

type response struct {
	Message string `json:"message"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entrou aqui")
	bodyRequest, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.Write([]byte("Failed to read body request"))
		return
	}

	var user user

	if erro = json.Unmarshal(bodyRequest, &user); erro != nil {
		w.Write([]byte("Something wrong happened"))
		return
	}

	// fmt.Println(user)

	db, erro := db.Connectar()

	if erro != nil {
		w.Write([]byte("Failed to connecting with DB"))
		return
	}

	defer db.Close()

	// INSERT INTO usuarios (nome, email) values ("nome", "email")
	// PREPARE STATEMENT
	// ELE É UTILIZADO PARA EVITAR O SQL INJECTION

	statement, erro := db.Prepare("INSERT INTO usuarios (nome, email) values (?, ?)")

	if erro != nil {
		w.Write([]byte("Error creating statement in DB"))
		return
	}

	defer statement.Close()

	insercao, erro := statement.Exec(user.Nome, user.Email)

	if erro != nil {
		w.Write([]byte("Error to insert data in DB"))
		return
	}

	// Se chegar aqui, usuário foi inserido
	// Vou retornar o ID que foi inserido

	idInsercao, erro := insercao.LastInsertId()

	if erro != nil {
		w.Write([]byte("Error to get ID insert in DB"))
		return
	}

	// json.NewEncoder(w).Encode(map[string]string{"status": "OK"})

	//STATUS CODE: 201 CREATED, 404 NOT FOUND, 204 NOT CONTENT
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseRequest := response{Message: fmt.Sprintf("User created with success! ID: %d", idInsercao)}

	responseJson, erro := json.Marshal(responseRequest)

	if erro != nil {
		w.Write([]byte("Something wrong happened"))
		return
	}

	w.Write([]byte(responseJson))

}
