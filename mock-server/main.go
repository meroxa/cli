package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var memDB map[string]map[string]interface{}

func main() {
	// Init in-memory database
	memDB = make(map[string]map[string]interface{})

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	// Routes
	r.HandleFunc("/{object}", listHandler).Methods("GET")
	r.HandleFunc("/{object}", createHandler).Methods("POST")
	r.HandleFunc("/{object}/{id}", describeHandler).Methods("GET")

	// Run Server
	err := http.ListenAndServe(":8080", r)
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	object := vars["object"]
	list := memDB[object]

	responseJSON, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //nolint:gocritic
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	object := vars["object"]

	if memDB[object] == nil {
		memDB[object] = make(map[string]interface{})
	}

	decoder := json.NewDecoder(r.Body)
	var o map[string]interface{}
	err := decoder.Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextID := strconv.Itoa(len(memDB[object]) + 1)
	memDB[object][nextID] = o
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s created", object)
}

func describeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	object := vars["object"]
	id := vars["id"]

	res := memDB[object][id]

	responseJSON, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //nolint:gocritic
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}
