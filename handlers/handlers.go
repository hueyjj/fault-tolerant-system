// Package handlers defines the http handler functions
// that will be used in the homework 1 program
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/cmps128gofour/homework2/response"
	"bitbucket.org/cmps128gofour/homework2/store"
	"github.com/gorilla/mux"
)

// KVStore is the blobal key value store for handling reads, writes, and lookup
var KVStore = store.New()

// Sends and displays a response.
// A status code 200 is given
// when a client makes a GET request to the /hello resource.
func helloGET(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello world!")
}

// Makes sure that a POST request is not allowed to be called to /hello resource.
// Status code 405 is given when attempted.
func helloPOST(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /hello not supported")
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Extracts a message from the request, then responds and displays it.
// Messages should be Alphanumeric.
// A status code 200 is given
// when a client makes a POST request to the /test resource.
func testPOST(w http.ResponseWriter, r *http.Request) {
	// Try to parse the form for values, and error out if there's none
	if err := r.ParseForm(); err != nil {
		log.Printf("could not parse form: %v\n", err)
		http.Error(w, "could not parse attached form", http.StatusBadRequest)
		return
	}

	// Otherwise get the value from msg, again checking for errors
	var msg string
	val, ok := r.Form["msg"]

	if !ok {
		log.Println("received request with missing key \"msg\"")
	} else {
		msg = val[0]
	}

	// If successful, then we can send the OK status, and respond with the correct msg
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "POST message received: %s", msg)
}

// Sends a response and a status code 200 is given
// when a client makes a GET request to the /test resource.
func testGET(w http.ResponseWriter, r *http.Request) {
	// Send the OK status, and respond with the correct msg
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "GET request received")
}

func subjectPUT(w http.ResponseWriter, r *http.Request) {
	var resp *response.Response

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]
	value := r.PostFormValue("val")

	// Return error message if key is 1 and 200 characters
	if len(key) < 1 || len(key) > 200 {
		resp = &response.Response{
			Msg:    "Key not valid",
			Result: "Error",
		}
		w.WriteHeader(http.StatusBadRequest) //(TODO:Jasper Jeng) Check with someone what the status code should be
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
			Replaced: "True",
			Msg:      "Updated successfully",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		// Put key into store if it doesn't exists, or replace key

		KVStore.Put(key, value)
		resp = &response.Response{
			Replaced: "False",
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
	var resp *response.Response

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	if KVStore.Exists(key) {
		value, _ := KVStore.Get(key)
		// We never actually enter here because we check if key exists already
		//if err != nil {
		//	log.Printf("Unable to get value: %v\n", err)
		//	http.Error(w, "Unable to find key", http.StatusBadRequest)
		//}

		resp = &response.Response{
			Result: "Success",
			Msg:    value,
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

func subjectDEL(w http.ResponseWriter, r *http.Request) {
	var resp *response.Response

	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

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

// Serve creates a server that can be gracefully shutdown,
// and handles the routes as defined in the homework 1 spec
func Serve(ip string, port string) {

	router := mux.NewRouter()
	// Add handlers here
	router.HandleFunc("/hello", helloGET).Methods("GET")
	router.HandleFunc("/hello", helloPOST).Methods("POST")
	router.HandleFunc("/test", testPOST).Methods("POST")
	router.HandleFunc("/test", testGET).Methods("GET")
	router.HandleFunc("/keyValue-store/{subject}", subjectPUT).Methods("PUT")
	router.HandleFunc("/keyValue-store/{subject}", subjectGET).Methods("GET")
	router.HandleFunc("/keyValue-store/{subject}", subjectDEL).Methods("DELETE")

	// Run a server as defined by Gorilla mux, with graceful shutdown
	// ref: https://github.com/gorilla/mux#graceful-shutdown

	srv := &http.Server{
		Handler: router,
		// Since we're running in docker, we can't bind to localhost (AKA 127.0.0.1)
		// So we bind to 0.0.0.0 (AKA the global interface) instead
		// that way we can access it outside the docker container
		Addr:         fmt.Sprintf("%s:%s", ip, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run the server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	log.Println("started server on:", srv.Addr)

	// Make a channel, and send a value on that channel
	// whenever we get an os.Interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for, and shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("server shutting down")
}
