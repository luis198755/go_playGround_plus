package limiter

import (
	"sync"
	"time"
)

// RateLimiterInterface define el comportamiento de un limitador de tasa
type RateLimiterInterface interface {
	IsAllowed(ip string) bool
}

// IPRequests almacena información sobre las solicitudes de una IP
type IPRequests struct {
	count    int
	lastTime time.Time
}

// RateLimiter implementa un limitador de tasa basado en IP
type RateLimiter struct {
	requests          map[string]*IPRequests
	mu                sync.RWMutex
	maxRequestsPerMin int
}

// NewRateLimiter crea un nuevo limitador de tasa
func NewRateLimiter(maxRequestsPerMin int) *RateLimiter {
	return &RateLimiter{
		requests:          make(map[string]*IPRequests),
		maxRequestsPerMin: maxRequestsPerMin,
	}
}

// IsAllowed verifica si una IP está permitida para hacer una solicitud
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if req, exists := rl.requests[ip]; exists {
		if now.Sub(req.lastTime) > time.Minute {
			req.count = 1
			req.lastTime = now
			return true
		}
		if req.count < rl.maxRequestsPerMin {
			req.count++
			return true
		}
		return false
	}
	rl.requests[ip] = &IPRequests{count: 1, lastTime: now}
	return true
}
