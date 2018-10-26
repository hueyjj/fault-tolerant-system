package main

import (
	"log"
	"os"

	"bitbucket.org/cmps128gofour/homework2/handlers"
)

func main() {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	mainIP := os.Getenv("MAINIP")

	if ip == "" {
		ip = "0.0.0.0"
	}

	if port == "" {
		port = "8080"
	}

	if mainIP == "" {
		log.Print("Starting main server")
		handlers.Serve(ip, port)
	} else {
		log.Print("Starting forwarding server")
		handlers.ForwardServe(ip, port, mainIP)
	}
}
