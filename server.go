// main.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Pessoa struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

var db []Pessoa
var proximoID = 1 // Variável para controlar o próximo ID

func getPessoas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(db)
}

func getPessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, p := range db {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func createPessoa(w http.ResponseWriter, r *http.Request) {
	var nova Pessoa
	err := json.NewDecoder(r.Body).Decode(&nova)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nova.ID = proximoID
	proximoID++
	db = append(db, nova)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nova)
}

func updatePessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var update Pessoa
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, p := range db {
		if p.ID == id {
			update.ID = id // Atualiza o ID na estrutura de atualização
			db[i] = update
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(update)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deletePessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, p := range db {
		if p.ID == id {
			db = append(db[:i], db[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pessoas", getPessoas).Methods("GET")
	r.HandleFunc("/pessoas/{id}", getPessoa).Methods("GET")
	r.HandleFunc("/pessoas", createPessoa).Methods("POST")
	r.HandleFunc("/pessoas/{id}", updatePessoa).Methods("PUT")
	r.HandleFunc("/pessoas/{id}", deletePessoa).Methods("DELETE")
	fmt.Println("Servidor rodando na porta :8080")
	http.ListenAndServe(":8080", r)
}
