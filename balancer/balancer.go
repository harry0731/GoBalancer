package balancer

import (
	"GoBalancer/models"
	"GoBalancer/tools"
	"log"
	"net/http"
)

var WorkersPool models.ＷorkerPool

// lb load balances the incoming request
func Lb(w http.ResponseWriter, r *http.Request) {
	attempts := tools.GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := WorkersPool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
