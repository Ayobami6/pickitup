package main

import (
	"log"

	"github.com/Ayobami6/pickitup/cmd/api"
)

func main() {
	addr := "localhost:2300"
	server := api.NewAPIServer(addr)

	if err := server.Run(); err!= nil {
        log.Fatal(err)
    }
}