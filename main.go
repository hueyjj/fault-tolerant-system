package main

import (
	"os"

	"bitbucket.org/cmps128gofour/homework3/handlers"
)

func main() {
	ipPort := os.Getenv("IP_PORT")
	handlers.Serve(ipPort)
}
