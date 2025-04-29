package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"loadbalancer/balancer"
	"loadbalancer/config"
)

func main() {
	cfg := config.LoadConfig("config.yml")

	backends := []*balancer.Backend{}
	for _, b := range cfg.Backends {
		backends = append(backends, &balancer.Backend{
			URL:     b.URL,
			Healthy: true,
		})
	}

	healthChecker := &balancer.HealthChecker{
		Backends: backends,
		Interval: time.Duration(cfg.HealthCheckIntervalSeconds) * time.Second,
	}
	healthChecker.Start()

	// Load selection algorithm plugin
	var selector balancer.Selector
	switch cfg.SelectionAlgorithm {
	case "round_robin":
		selector = &balancer.RoundRobinSelector{}
	case "random":
		selector = &balancer.RandomSelector{}
	default:
		log.Fatalf("Unknown selection algorithm: %s", cfg.SelectionAlgorithm)
	}

	lb := &balancer.LoadBalancer{
		HealthChecker: healthChecker,
		Selector:      selector,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: lb,
	}

	// Graceful shutdown
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Printf("Load balancer listening on port %d (TLS=%v)\n", cfg.Port, cfg.TLSEnabled)

	var err error
	if cfg.TLSEnabled {
		err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
