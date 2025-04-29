package balancer

import (
	"net/http"
	"sync"
	"time"
)

type HealthChecker struct {
	Backends []*Backend
	Interval time.Duration
	mutex    sync.RWMutex
}

func (h *HealthChecker) Start() {
	go func() {
		for {
			h.checkHealth()
			time.Sleep(h.Interval)
		}
	}()
}

func (h *HealthChecker) checkHealth() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, backend := range h.Backends {
		resp, err := http.Get(backend.URL + "/")
		if err != nil || resp.StatusCode >= 400 {
			backend.Healthy = false
		} else {
			backend.Healthy = true
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func (h *HealthChecker) GetHealthyBackends() []*Backend {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var healthy []*Backend
	for _, backend := range h.Backends {
		if backend.Healthy {
			healthy = append(healthy, backend)
		}
	}
	return healthy
}

func (h *HealthChecker) GetMetrics() string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	metrics := ""
	for _, backend := range h.Backends {
		status := "healthy"
		if !backend.Healthy {
			status = "unhealthy"
		}
		metrics += backend.URL + " -> " + status + "\n"
	}
	return metrics
}
