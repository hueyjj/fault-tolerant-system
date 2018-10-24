package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	clientResponse "bitbucket.org/cmps128gofour/homework2/response"
	"github.com/gorilla/mux"
)

var mainIP string
var client = http.Client{}

func proxySubjectGET(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	// make request
	requestString := fmt.Sprintf("%s/keyValue-store/%s", mainIP, key)
	request, err := http.NewRequest(http.MethodGet, requestString, nil)

	if err != nil {
		log.Println("could not make request:", err)
		return
	}

	// Send the request
	response, err := client.Do(request)

	if err != nil {
		log.Println("could not get response:", err)
		// Main server is down
		respondError501(w)
	}
	log.Println("response status code: ", response.StatusCode)

	// Write the header
	w.WriteHeader(response.StatusCode)
	w.Header().Set("content-type", "application/json")
	// Read the body
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("could not read body:", err)
		return
	}

	// Write the body
	_, err = w.Write(body)
	if err != nil {
		log.Println("could not write body:", err)
	}

}

func proxySubjectDEL(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]

	// make request
	requestString := fmt.Sprintf("%s/keyValue-store/%s", mainIP, key)
	request, err := http.NewRequest(http.MethodDelete, requestString, nil)

	if err != nil {
		log.Println("could not make request:", err)
		return
	}

	// Send the request
	response, err := client.Do(request)

	if err != nil {
		log.Println("could not get response:", err)
		// Main server is down
		respondError501(w)
	}

	fmt.Println(response)

	// Write the header
	w.WriteHeader(response.StatusCode)
	w.Header().Set("content-type", "application/json")
	// Read the body
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("could not read body:", err)
		return
	}

	// Write the body
	_, err = w.Write(body)
	if err != nil {
		log.Println("could not write body:", err)
	}

}

func proxySubjectPUT(w http.ResponseWriter, r *http.Request) {
	// Parse the key from url variable and (store) value from the request
	vars := mux.Vars(r)
	key := vars["subject"]
	value := r.PostFormValue("val")

	// make request
	requestString := fmt.Sprintf("%s/keyValue-store/%s?val=%s", mainIP, key, value)
	request, err := http.NewRequest(http.MethodPut, requestString, nil)

	if err != nil {
		log.Println("could not make request:", err)
		return
	}

	// Send the request
	response, err := client.Do(request)

	if err != nil {
		log.Println("could not get response:", err)
		// Main server is down
		respondError501(w)
	}

	fmt.Println(response)

	// Write the header
	w.WriteHeader(response.StatusCode)
	w.Header().Set("content-type", "application/json")
	// Read the body
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("could not read body:", err)
		return
	}

	// Write the body
	_, err = w.Write(body)
	if err != nil {
		log.Println("could not write body:", err)
	}

}

// Have the ResponseWriter write response 501 in JSON format.
func respondError501(w http.ResponseWriter) {

	w.WriteHeader(http.StatusNotImplemented)
	w.Header().Set("content-type", "application/json")
	var resp *clientResponse.Response
	resp = &clientResponse.Response{
		Result: "Error",
		Msg:    "Server unavailable",
	}

	// Convert response into json structure and then into bytes
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println("could not write body:", err)
	}
}

// Serve creates a server that can be gracefully shutdown,
// and handles the routes as defined in the homework 1 spec
func ForwardServe(ip string, port string, mIP string) {
	mainIP = "http://" + mIP
	router := mux.NewRouter()
	// Add handlers here
	router.HandleFunc("/hello", helloGET).Methods("GET")
	router.HandleFunc("/hello", helloPOST).Methods("POST")
	router.HandleFunc("/test", testPOST).Methods("POST")
	router.HandleFunc("/test", testGET).Methods("GET")
	router.HandleFunc("/keyValue-store/{subject}", proxySubjectPUT).Methods("PUT")
	router.HandleFunc("/keyValue-store/{subject}", proxySubjectGET).Methods("GET")
	router.HandleFunc("/keyValue-store/{subject}", proxySubjectDEL).Methods("DELETE")

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
