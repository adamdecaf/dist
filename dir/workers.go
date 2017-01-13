package main

import (
	"encoding/json"
	"github.com/adamdecaf/dist/dist"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"sync"
)

var (
	directory []dist.Address
	dirLock sync.RWMutex
)

// todo:
// filter out bad ips
// - lookback, invalid, etc

func register(address dist.Address) {
	dirLock.Lock()
	defer dirLock.Unlock()
	directory = append(directory, address)
	log.Printf("registered %s", address.String())
}

// RegisterRoute accepts a worker and adds it to dir
func RegisterRoute(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("address")
	port, err := strconv.Atoi(r.URL.Query().Get("port"))
	if (addr == "" || port <= 0) || err != nil {
		http.Error(w, "'address' and/or 'port' query params are bad, please fix", http.StatusBadRequest)
		return
	}
	ip := net.ParseIP(addr)
	register(dist.Address{
		IP: ip,
		Port: port,
	})
	w.WriteHeader(204) // 'No Content'
}

// findWorkers attempts to find as many workers as possible, but cannot ensure
// that the number workers found is always what the client has requested.
// However, more workers than requested will never be returned.
func findWorkers(count int) []dist.Address {
	dirLock.RLock()
	defer dirLock.RUnlock()

	if count >= len(directory) {
		return directory
	}

	start := rand.Int() % len(directory)
	return directory[start:start+count]
}

type Workers struct {
	Workers []dist.Address `json:"workers"`
}

// WorkersRoute reads the requested worker count and returns it for clients.
func WorkersRoute(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		http.Error(w, "query param 'count' is not an integer", http.StatusInternalServerError)
		return
	}

	workers := Workers{
		Workers: findWorkers(count),
	}
	json, err := json.Marshal(workers)
	if err != nil {
		http.Error(w, "Unable to find workers, something went wrong.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
