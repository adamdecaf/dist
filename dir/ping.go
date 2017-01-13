package dir

import (
	"fmt"
	"net/http"
)

func PingRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}
