package balancer

import "net/http"

// lb load balances the incoming request
func lb(w http.ResponseWriter, r *http.Request) {
	peer := workerpool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
