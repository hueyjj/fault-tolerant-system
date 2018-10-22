package main

import (
	"fmt"
	"os"

	"bitbucket.org/cmps128gofour/homework2/handlers"
)

func main() {

	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	mainIP := os.Getenv("MAINIP")

	if mainIP == "" {
		fmt.Println("MAIN")
		handlers.Serve(ip, port)
	} else {
		fmt.Println("FORWARD")
		handlers.ForwardServe(ip, port, mainIP)
	}

}
