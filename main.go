package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Customer struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     uint   `json:"phone"`
	Contacted bool   `json:"contacted"`
}

var customers = map[string]Customer{
	"1": Customer{
		ID:        1,
		Name:      "Alex",
		Role:      "Customer",
		Email:     "alex@gmail.com",
		Phone:     5552223344,
		Contacted: true,
	},
	"2": Customer{
		ID:        2,
		Name:      "Lake",
		Role:      "Customer",
		Email:     "lake@gmail.com",
		Phone:     4445557788,
		Contacted: true,
	},
	"3": Customer{
		ID:        3,
		Name:      "Itzel",
		Role:      "Customer",
		Email:     "itzel@gmail.com",
		Phone:     4449995533,
		Contacted: false,
	},
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	customer, found := customers[id]

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Customer not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCustomer Customer

	reqBody, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(reqBody, &newCustomer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	id := fmt.Sprintf("%d", newCustomer.ID)

	_, ok := customers[id]

	if ok {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Customer already exists"})
		return
	}

	customers[id] = newCustomer

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomer)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	customer, found := customers[id]

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Customer not found"})
		return
	}

	var updateData map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	for key, value := range updateData {
		switch key {
		case "name":
			customer.Name = value.(string)
		case "role":
			customer.Role = value.(string)
		case "email":
			customer.Email = value.(string)
		case "phone":
			customer.Phone = uint(value.(float64))
		case "contacted":
			customer.Contacted = value.(bool)
		}
	}

	customers[id] = customer

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	if _, ok := customers[id]; ok {
		delete(customers, id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customers)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(customers)
	}
}

func getAPIInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	apiInfo := map[string]interface{}{
		"api_name":    "Customers API",
		"description": "This is an API for managing customers.",
		"endpoints": []map[string]string{
			{
				"url":         "/customers",
				"description": "Get a list of customers.",
			},
			{
				"url":         "/customers/{id}",
				"description": "Get details of a customer.",
			},
			{
				"url":         "/customers",
				"description": "Add a new customer.",
			},
			{
				"url":         "/customers/{id}",
				"description": "Update an existing customer.",
			},
			{
				"url":         "/customers/{id}",
				"description": "Delete an existing customer.",
			},
		},
	}
	json.NewEncoder(w).Encode(apiInfo)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PATCH")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/", getAPIInfo).Methods("GET")

	fmt.Println("The server is running on port 3000...")
	http.ListenAndServe(":3000", router)
}
