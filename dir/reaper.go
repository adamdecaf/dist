package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	"sync"
	"errors"
	"github.com/adamdecaf/dist/dist"
)

const (
	ReapInterval = 1 * time.Second
	MaxBatchSize = 100
)

var (
	ErrNodeDown = errors.New("node down")
)

func startReaper() {
	log.Println("starting worker reaper")
	go reaper()
}

func reaper() {
	for _ = range time.Tick(ReapInterval) {
		size := MaxBatchSize

		if countWorkers() > size {
			log.Println("# of workers is large than our reaping batch, need to get a better process here..")
		}

		workers := findWorkers(size)

		if len(workers) <= 0 {
			continue
		}

		wg := sync.WaitGroup{}
		wg.Add(len(workers))

		mux := sync.RWMutex{}
		flushable := make([]dist.Address, 0)

		for i := range workers {
			w := workers[i]
			// Check if we need to drop this worker
			go func() {
				err := checkWorker(w)
				if err != nil {
					mux.Lock()
					flushable = append(flushable, w)
					mux.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait() // wait for all checks to finish

		// Remove all bad workers
		removeWorkers(flushable)
	}
}

// todo: timeouts and other things
func checkWorker(addr dist.Address) error {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/ping", addr.IP.String(), addr.Port))
	if err != nil {
		return ErrNodeDown
	}
	if resp.StatusCode != 200 {
		return ErrNodeDown
	}
	return nil
}
