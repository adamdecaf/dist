package main

import (
	"encoding/json"
	"github.com/adamdecaf/dist/dist"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

var (
	directory []*dist.Address
	dirLock sync.RWMutex
)

// todo:
// filter out bad ips
// - lookback, invalid, etc
// - don't add duplicates, or bad addresses

func register(address dist.Address) {
	dirLock.Lock()
	defer dirLock.Unlock()

	if address.Valid() {
		directory = append(directory, &address)
		log.Printf("registered %s", address.String())
	}
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

	// reset count if they've asked for too much
	max := len(directory)
	if count > max {
		count = max
	}

	out := make([]dist.Address, 0)
	for _,v := range directory {
		out = append(out, *v)
	}

	return out
}

// countWorkers returns the current count of workers
func countWorkers() int {
	dirLock.RLock()
	defer dirLock.RUnlock()
	return len(directory)
}

// removeWorkers is to remove workers from the active pool
func removeWorkers(addrs []dist.Address) error {
	dirLock.Lock()
	defer dirLock.Unlock()

	// cs 101 here we come
	res := make([]*dist.Address, 0)
	for i := range directory {
		rm := false
		for k := range addrs {
			if directory[i].Equal(addrs[k]) {
				rm = true
			}
		}
		if !rm {
			if directory[i].Valid() {
				res = append(res, directory[i])
			}
		}
	}
	directory = res
	return nil
}

// Json for returning workers
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
