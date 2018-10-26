package handlers

import (
	"fmt"
	"log"
	"net/http"
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
