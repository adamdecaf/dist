package workers

import (
	"fmt"
	"log"
	"net/http"
	"github.com/adamdecaf/dist/dist"
)

func Listen(addr dist.Address) {
	http.HandleFunc("/ping", ping)

	listen := fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}
