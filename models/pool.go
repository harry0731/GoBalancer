package models

import (
	"GoBalancer/tools"
	"log"
	"net/url"
	"sync/atomic"
	"time"
)

type ＷorkerPool struct {
	workers []*Worker
	current uint64
}

// AddBackend to the server pool
func (s *ＷorkerPool) AddBackend(worker *Worker) {
	s.workers = append(s.workers, worker)
}

func (s *ＷorkerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.workers)))
}

// MarkBackendStatus changes a status of a backend
func (s *ＷorkerPool) MarkBackendStatus(workerUrl *url.URL, alive bool) {
	for _, b := range s.workers {
		if b.URL.String() == workerUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer returns next active peer to take a connection
func (s *ＷorkerPool) GetNextPeer() *Worker {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	l := len(s.workers) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.workers) // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one
		if s.workers[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx)) // mark the current one
			}
			return s.workers[idx]
		}
	}
	return nil
}

// HealthCheck pings the backends and update the status
func (s *ＷorkerPool) healthCheck() {
	for _, b := range s.workers {
		status := "up"
		alive := tools.IsWorkerAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

// healthCheck runs a routine for check status of the backends every 2 mins
func (s *ＷorkerPool) HealthCheck() {
	t := time.NewTicker(time.Minute * 2)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			s.healthCheck()
			log.Println("Health check completed")
		}
	}
}
