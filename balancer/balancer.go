package balancer

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

type LoadBalancer struct {
	HealthChecker *HealthChecker
	Selector      Selector
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/lb/metric" {
		w.Write([]byte(lb.HealthChecker.GetMetrics()))
		return
	}

	backend := lb.Selector.Select(lb.HealthChecker.GetHealthyBackends())
	if backend == nil {
		http.Error(w, "No healthy backends", http.StatusServiceUnavailable)
		return
	}

	backendURL, _ := url.Parse(backend.URL)
	proxyReq, err := http.NewRequest(r.Method, backendURL.String()+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}
	proxyReq.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Backend unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	log.Printf("%s %s -> %s", r.Method, r.URL.Path, backend.URL)
}
