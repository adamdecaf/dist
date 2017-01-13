package main

import (
	"github.com/adamdecaf/dist/dist"
	"github.com/adamdecaf/dist/workers"
	"log"
)

func main() {
	// register
	addr, err := dist.Register()
	if err != nil {
		log.Fatalf("error on startup, err=%s", err)
	}

	// Respond to pings for dir
	workers.Listen(addr)
}
