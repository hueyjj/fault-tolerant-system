package main

import (
	"os"

	"bitbucket.org/cmps128gofour/homework2/handlers"
)

func main() {

	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	mainIP := os.Getenv("MAINIP")

	if mainIp != "" {
		handlers.Serve(ip, port)
	} else {
		handlers.ForwardServe(ip, port, mainIP)
	}

}
