package balancer

import (
	"math/rand"
	"sync/atomic"
)

type Selector interface {
	Select(backends []*Backend) *Backend
}

// RoundRobinSelector plugin
type RoundRobinSelector struct {
	counter uint64
}

func (r *RoundRobinSelector) Select(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	idx := atomic.AddUint64(&r.counter, 1)
	return backends[(idx-1)%uint64(len(backends))]
}

// RandomSelector plugin
type RandomSelector struct{}

func (r *RandomSelector) Select(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	idx := rand.Intn(len(backends))
	return backends[idx]
}
