package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	log.Println("starting dir")

	startReaper()

	// start our http server
	http.HandleFunc("/ping", PingRoute)
	http.HandleFunc("/workers", WorkersRoute)
	http.HandleFunc("/register", RegisterRoute)

	listen := fmt.Sprintf("127.0.0.1:%d", 8080)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
