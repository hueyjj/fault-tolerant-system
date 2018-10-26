package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"bitbucket.org/cmps128gofour/homework2/response"
	"github.com/gorilla/mux"
)

func subjectPUT(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]
	value := r.PostFormValue("val")

	var resp *response.Response
	// Return error message if key is 1 and 200 characters
	if len(key) < 1 || len(key) > 200 {
		resp = &response.Response{
			Msg:    "Key not valid",
			Result: "Error",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else if len(value) > 1000000 {
		// Return error message if value is greater than 1 MB

		resp = &response.Response{
			Msg:    "Object too large. Size limit is 1MB",
			Result: "Error",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else if KVStore.Exists(key) {
		// Replace value in store

		KVStore.Put(key, value)
		resp = &response.Response{
			Replaced: true,
			Msg:      "Updated successfully",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		// Put key into store if it doesn't exists, or replace key

		KVStore.Put(key, value)
		resp = &response.Response{
			Replaced: false,
			Msg:      "Added successfully",
		}
		w.WriteHeader(http.StatusCreated)
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func subjectGET(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	var resp *response.Response
	if KVStore.Exists(key) {
		value, _ := KVStore.Get(key)
		resp = &response.Response{
			Result: "Success",
			Value:  value,
		}
		w.WriteHeader(http.StatusOK)
	} else {
		resp = &response.Response{
			Result: "Error",
			Msg:    "Not Found",
		}
		w.WriteHeader(http.StatusNotFound)
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func subjectSEARCH(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	var resp *response.Response
	if KVStore.Exists(key) {
		resp = &response.Response{
			Result:   "Success",
			IsExists: true,
		}
	} else {
		resp = &response.Response{
			Result:   "Success",
			IsExists: false,
		}
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func subjectDEL(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	var resp *response.Response
	if KVStore.Exists(key) {
		err := KVStore.Delete(key)
		if err != nil {
			log.Printf("Unable to delete key: %v\n", err)
			http.Error(w, "Unable to delete key", http.StatusBadRequest)
		}

		resp = &response.Response{
			Result: "Success",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		resp = &response.Response{
			Result: "Error",
			Msg:    "Status code 404",
		}
		w.WriteHeader(http.StatusNotFound)
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
