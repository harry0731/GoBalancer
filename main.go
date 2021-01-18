package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	Attempts int = iota
	Retry
)

// GetAttemptsFromContext returns the attempts for request
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

// isAlive checks whether a backend is Alive by establishing a TCP connection
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	_ = conn.Close() // close it, we dont need to maintain this connection
	return true
}

// HealthCheck pings the backends and update the status
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}

func healthCheck() {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

func main() {
	go healthCheck()
}
