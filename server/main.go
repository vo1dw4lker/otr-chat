package main

import (
	"flag"
	"log"
	"serverenc/core"
)

func main() {
	srv := core.NewServer()
	srv.Start(parsePort())
}

func parsePort() int {
	var port int
	flag.IntVar(&port, "p", 7575, "Port number to listen on")
	flag.Parse()

	if port < 1 || port > 65535 {
		log.Fatalln("Invalid port number. Port must be between 1 and 65535.")
	}

	return port
}
