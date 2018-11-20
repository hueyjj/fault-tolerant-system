package main

import (
	"os"

	"bitbucket.org/cmps128gofour/homework3/handlers"
)

func main() {
	ipPort := os.Getenv("IP_PORT")
	views := os.Getenv("VIEW")
	handlers.SetViews(views)
	handlers.SetMyIpPort(ipPort)
	handlers.Serve(ipPort)
}
