// Package handlers defines the http handler functions
// that will be used in the cs128 homeworks
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/cmps128gofour/homework3/store"
	"github.com/gorilla/mux"
)

// KVStore is the global key value store for handling reads, writes, and lookup
var KVStore = store.New()

// Serve creates a server that can be gracefully shutdown,
// and handles the routes as defined in the homework 1 spec
func Serve(ipPort string) {

	router := mux.NewRouter()
	// Add handlers here
	router.HandleFunc("/hello", helloGET).Methods("GET")
	router.HandleFunc("/hello", helloPOST).Methods("POST")
	router.HandleFunc("/test", testPOST).Methods("POST")
	router.HandleFunc("/test", testGET).Methods("GET")
	router.HandleFunc("/keyValue-store/{subject}", subjectPUT).Methods("PUT")
	router.HandleFunc("/keyValue-store/{subject}", subjectGET).Methods("GET")
	router.HandleFunc("/keyValue-store/search/{subject}", subjectSEARCH).Methods("GET")
	router.HandleFunc("/keyValue-store/{subject}", subjectDEL).Methods("DELETE")
	router.HandleFunc("/keyValue-store/{subject}", subjectDEL).Methods("POST")
	router.HandleFunc("/view", viewGET).Methods("GET")
	router.HandleFunc("/view", viewPUT).Methods("PUT")
	router.HandleFunc("/view", viewDELETE).Methods("DELETE")
	router.HandleFunc("/view", viewDELETE).Methods("POST")

	// Run a server as defined by Gorilla mux, with graceful shutdown
	// ref: https://github.com/gorilla/mux#graceful-shutdown

	srv := &http.Server{
		Handler: router,
		// Since we're running in docker, we can't bind to localhost (AKA 127.0.0.1)
		// So we bind to 0.0.0.0 (AKA the global interface) instead
		// that way we can access it outside the docker container
		Addr:         ipPort,
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

// ForwardServe creates a forwarding server that can be gracefully shutdown,
// and handles the routes as defined in the homework 1 spec
func ForwardServe(ip string, port string, mIP string) {
	mainIP = "http://" + mIP
	router := mux.NewRouter()
	// Add handlers here
	router.HandleFunc("/hello", helloGET).Methods("GET")
	router.HandleFunc("/hello", helloPOST).Methods("POST")
	router.HandleFunc("/test", testPOST).Methods("POST")
	router.HandleFunc("/test", testGET).Methods("GET")
	// Homework 2
	router.HandleFunc("/keyValue-store/{subject}", proxySubjectPUT).Methods("PUT")
	router.HandleFunc("/keyValue-store/{subject}", proxySubjectGET).Methods("GET")
	router.HandleFunc("/keyValue-store/search/{subject}", proxySubjectSEARCH).Methods("GET")
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
