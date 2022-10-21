package main

import (
	"Go-server/server"
)

func main() {
	server := server.NewServer()
	if err := server.Start(); err != nil {
		panic("Couldn't start the server")
	}
}
