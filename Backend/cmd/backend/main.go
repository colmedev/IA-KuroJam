package main

import (
	"log"

	"github.com/colmedev/IA-KuroJam/Backend/server"
)

var (
	version = "0.0.1"
)

func main() {
	err := server.StartServer()
	if err != nil {
		log.Fatal("error starting app", err)
	}
}
