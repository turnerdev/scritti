package main

import (
	"log"
	server "scritti/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := 9090
	server.Server(port)
}
