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

	if mainIP == "" {
		log.Print("Starting main server")
		handlers.Serve(ip, port)
	} else {
		log.Print("Starting forwarding server")
		handlers.ForwardServe(ip, port, mainIP)
	}
}
