package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var clients = map[string]string{
	"1": "Alex",
	"2": "Lakeisha",
	"3": "Itzel",
	"4": "Pablo",
	"5": "Ximena",
	"6": "Arturo",
	"7": "Paola",
}

func getClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(clients)
}

func getClientId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	client, found := clients[id]

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Client not found"})
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{id: client})
	}
}

func createClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newClient map[string]string

	reqBody, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(reqBody, &newClient)

	for k, v := range newClient {
		if _, ok := clients[k]; ok {
			w.WriteHeader(http.StatusConflict)
		} else {
			clients[k] = v
			w.WriteHeader(http.StatusCreated)
		}

		json.NewEncoder(w).Encode(clients)
	}
}

func updateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	_, found := clients[id]

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Client not found"})
		return
	}
	var updateClient map[string]string

	if err := json.NewDecoder(r.Body).Decode(&updateClient); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	updateValue, found := updateClient[id]
	if found {
		clients[id] = updateValue
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Definition update successfully"})
}

func deleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	if _, ok := clients[id]; ok {
		delete(clients, id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(clients)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(clients)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/clients", getClients).Methods("GET")
	router.HandleFunc("/client/{id}", getClientId).Methods("GET")
	router.HandleFunc("/postClient", createClient).Methods("POST")
	router.HandleFunc("/updateClient/{id}", updateClient).Methods("PATCH")
	router.HandleFunc("/deleteClient/{id}", deleteClient).Methods("DELETE")

	fmt.Println("The server is running on port 3000...")
	http.ListenAndServe(":3000", router)
}
