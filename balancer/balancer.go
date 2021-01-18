package balancer

import "net/http"

// lb load balances the incoming request
func lb(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
	  log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
	  http.Error(w, "Service not available", http.StatusServiceUnavailable)
	  return
	}
  
	peer := serverPool.GetNextPeer()
	if peer != nil {
	  peer.ReverseProxy.ServeHTTP(w, r)
	  return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
  }

proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
	log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
	retries := GetRetryFromContext(request)
	if retries < 3 {
	  select {
		case <-time.After(10 * time.Millisecond):
		  ctx := context.WithValue(request.Context(), Retry, retries+1)
		  proxy.ServeHTTP(writer, request.WithContext(ctx))
		}
		return
	  }
  
	// after 3 retries, mark this backend as down
	serverPool.MarkBackendStatus(serverUrl, false)
  
	// if the same request routing for few attempts with different backends, increase the count
	attempts := GetAttemptsFromContext(request)
	log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
	ctx := context.WithValue(request.Context(), Attempts, attempts+1)
	lb(writer, request.WithContext(ctx))
  }