package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"bitbucket.org/cmps128gofour/homework3/response"
	"bitbucket.org/cmps128gofour/homework3/vectorclock"
	"github.com/gorilla/mux"
)

var (
	vectorClock = vectorclock.New()
)

func subjectPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]
	value := r.PostFormValue("val")

	var resp *response.Response
	replaced := new(bool)
	// Return error message if key is 1 and 200 characters
	if len(key) < 1 || len(key) > 200 {
		resp = &response.Response{
			Msg:    "Key not valid",
			Result: "Error",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else if len(value) > 1e6 {
		// Return error message if value is greater than 1 MB
		resp = &response.Response{
			Msg:    "Object too large. Size limit is 1MB",
			Result: "Error",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else if KVStore.Exists(key) {
		// Replace value in store
		KVStore.Put(key, value)
		*replaced = true
		resp = &response.Response{
			Replaced: replaced,
			Msg:      "Updated successfully",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		// Put key into store if it doesn't exists
		KVStore.Put(key, value)
		*replaced = false
		resp = &response.Response{
			Replaced: replaced,
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
	w.Write(data)
}

func subjectGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	var resp *response.Response
	if KVStore.Exists(key) {
		value, _ := KVStore.Get(key)
		log.Printf(value)
		resp = &response.Response{
			Result: "Success",
			Value:  value,
		}
		w.WriteHeader(http.StatusOK)
	} else {
		resp = &response.Response{
			Result: "Error",
			Msg:    "Key does not exist",
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
	w.Write(data)
}

func subjectSEARCH(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	var resp *response.Response
	isExists := new(bool)
	if KVStore.Exists(key) {
		*isExists = true
		resp = &response.Response{
			Result:   "Success",
			IsExists: isExists,
		}
	} else {
		*isExists = false
		resp = &response.Response{
			Result:   "Success",
			IsExists: isExists,
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
	w.Header().Set("Content-Type", "application/json")

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
			Msg:    "Key deleted",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		resp = &response.Response{
			Result: "Error",
			Msg:    "Key does not exist",
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
	w.Write(data)
}

var iptable []string

func SetViews(views string) {
	iptable = strings.Split(views, ",")
}

func viewGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	iplist := strings.Join(iptable, ",")
	var resp *response.IPTableResponse
	resp = &response.IPTableResponse{
		View: iplist,
	}
	w.WriteHeader(http.StatusOK)

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}

func viewPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the key from url variable and (store) value from the request
	ipport := r.PostFormValue("ip_port")
	isIpportExist := false
	for _, ip := range iptable {
		if ip == ipport {
			isIpportExist = true
		}
	}

	var resp *response.ViewResponse
	if isIpportExist {
		resp = &response.ViewResponse{
			Result: "Error",
			Msg:    fmt.Sprintf("%s is already in view", ipport),
		}
		w.WriteHeader(http.StatusNotFound)
	} else {
		iptable = append(iptable, ipport)
		resp = &response.ViewResponse{
			Result: "Success",
			Msg:    fmt.Sprintf("Successfully added %s to view", ipport),
		}
		w.WriteHeader(http.StatusOK)
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(data)
}
