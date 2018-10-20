// Package handlers defines the http handler functions
// that will be used in the homework 1 program
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

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
	vars := mux.Vars(r)
	subject := vars["subject"]
	log.Printf(subject)
	log.Println("hi")
	w.WriteHeader(http.StatusOK)
}

func subjectGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := vars["subject"]
	log.Printf(subject)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hi")
}

func subjectDEL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := vars["subject"]
	log.Printf(subject)
	log.Println("hi")
	w.WriteHeader(http.StatusOK)
}

// Serve creates a server that can be gracefully shutdown,
// and handles the routes as defined in the homework 1 spec
func Serve() {

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
		Addr:         "0.0.0.0:8080",
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
