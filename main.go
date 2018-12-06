package main

import (
	"os"

	"bitbucket.org/cmps128gofour/homework4/handlers"
)

func main() {
	ipPort := os.Getenv("IP_PORT")
	views := os.Getenv("VIEW")
	s := os.Getenv("S")
	handlers.SetShardMap(s, views)
	handlers.SetViews(views)
	handlers.SetMyIpPort(ipPort)
	handlers.Serve(ipPort)
}
