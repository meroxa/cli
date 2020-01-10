package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var memDB map[string]map[string]interface{}

func listResourcesHandler(w http.ResponseWriter, r *http.Request) {
	resources := memDB["resources"]

	resourcesJSON, err := json.Marshal(resources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(resourcesJSON))
}

func createResourceHandler(w http.ResponseWriter, req *http.Request) {
	if memDB["resources"] == nil {
		memDB["resources"] = make(map[string]interface{})
	}

	decoder := json.NewDecoder(req.Body)
	var r map[string]interface{}
	err := decoder.Decode(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextID := strconv.Itoa(len(memDB["resources"]) + 1)
	memDB["resources"][nextID] = r
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "resource created")
}

func getResourceHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "here are the resources")
}

func main() {

	// Init in-memory database
	memDB = make(map[string]map[string]interface{})

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/resources", listResourcesHandler).Methods("GET")
	r.HandleFunc("/resources", createResourceHandler).Methods("POST")
	r.HandleFunc("/resources/{resource}", listResourcesHandler)

	// Run Server
	http.ListenAndServe(":8080", r)
}
