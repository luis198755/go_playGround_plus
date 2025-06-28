package limiter

import (
	"sync"
	"time"
)

// RateLimiterInterface define el comportamiento de un limitador de tasa
type RateLimiterInterface interface {
	IsAllowed(ip string) bool
}

// TokenBucket implementa el algoritmo de token bucket para rate limiting
type TokenBucket struct {
	tokens        float64    // Tokens actuales en el bucket
	capacity      float64    // Capacidad máxima del bucket
	refillRate    float64    // Tokens por segundo que se añaden
	lastRefillTime time.Time // Última vez que se rellenaron tokens
}

// RateLimiter implementa un limitador de tasa basado en IP usando token bucket
type RateLimiter struct {
	buckets       map[string]*TokenBucket
	mu           sync.RWMutex
	capacity     float64 // Capacidad máxima del bucket
	refillRate   float64 // Tokens por segundo que se añaden
}

// NewRateLimiter crea un nuevo limitador de tasa con algoritmo token bucket
func NewRateLimiter(maxRequestsPerMin int) *RateLimiter {
	// Convertimos solicitudes por minuto a tokens por segundo
	refillRate := float64(maxRequestsPerMin) / 60.0
	
	// La capacidad del bucket es igual al máximo de solicitudes por minuto
	// para permitir ráfagas controladas
	return &RateLimiter{
		buckets:     make(map[string]*TokenBucket),
		capacity:    float64(maxRequestsPerMin),
		refillRate:  refillRate,
	}
}

// IsAllowed verifica si una IP está permitida para hacer una solicitud usando token bucket
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	
	// Obtener o crear el bucket para esta IP
	bucket, exists := rl.buckets[ip]
	if !exists {
		// Para nuevas IPs, crear un bucket lleno
		bucket = &TokenBucket{
			tokens:        rl.capacity,
			capacity:      rl.capacity,
			refillRate:    rl.refillRate,
			lastRefillTime: now,
		}
		rl.buckets[ip] = bucket
		return true
	}
	
	// Calcular cuánto tiempo ha pasado desde la última recarga
	elapsed := now.Sub(bucket.lastRefillTime).Seconds()
	
	// Añadir tokens basados en el tiempo transcurrido
	newTokens := elapsed * bucket.refillRate
	bucket.tokens += newTokens
	
	// Limitar tokens a la capacidad máxima
	if bucket.tokens > bucket.capacity {
		bucket.tokens = bucket.capacity
	}
	
	// Actualizar el tiempo de la última recarga
	bucket.lastRefillTime = now
	
	// Verificar si hay suficientes tokens para esta solicitud
	if bucket.tokens >= 1.0 {
		// Consumir un token
		bucket.tokens -= 1.0
		return true
	}
	
	return false
}
