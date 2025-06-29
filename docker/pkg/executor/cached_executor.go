package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sync"
	"time"
)

// CacheEntry representa una entrada en el caché
type CacheEntry struct {
	Result      []byte
	LastAccess  time.Time
	AccessCount int
}

// CachedExecutor implementa un ejecutor con caché para código frecuentemente ejecutado
type CachedExecutor struct {
	executor     CodeExecutor
	cache        map[string]*CacheEntry
	cacheMutex   sync.RWMutex
	maxCacheSize int
	ttl          time.Duration
}

// NewCachedExecutor crea un nuevo ejecutor con caché
func NewCachedExecutor(executor CodeExecutor, maxCacheSize int, ttl time.Duration) *CachedExecutor {
	ce := &CachedExecutor{
		executor:     executor,
		cache:        make(map[string]*CacheEntry),
		maxCacheSize: maxCacheSize,
		ttl:          ttl,
	}
	
	// Iniciar rutina de limpieza periódica
	go ce.cleanupRoutine()
	
	return ce
}

// Execute ejecuta el código, utilizando el caché si está disponible
func (ce *CachedExecutor) Execute(ctx context.Context, code string, output io.Writer) error {
	// Generar hash del código como clave del caché
	codeHash := ce.hashCode(code)
	
	// Intentar obtener del caché
	ce.cacheMutex.RLock()
	entry, found := ce.cache[codeHash]
	if found {
		// Verificar si la entrada no ha expirado
		if time.Since(entry.LastAccess) <= ce.ttl {
			ce.cacheMutex.RUnlock()
			
			// Actualizar estadísticas del caché (en una goroutine separada para no bloquear)
			go ce.updateCacheStats(codeHash)
			
			// Escribir resultado desde el caché
			_, err := output.Write(entry.Result)
			return err
		}
		// La entrada ha expirado
		found = false
	}
	ce.cacheMutex.RUnlock()
	
	if !found {
		// Crear un buffer para capturar la salida
		buffer := &cachingWriter{
			buffer: make([]byte, 0, 4096), // Buffer inicial de 4KB
		}
		
		// Crear un escritor multi-destino
		multiWriter := io.MultiWriter(output, buffer)
		
		// Ejecutar el código
		err := ce.executor.Execute(ctx, code, multiWriter)
		if err != nil {
			return err
		}
		
		// Guardar en caché
		ce.cacheMutex.Lock()
		defer ce.cacheMutex.Unlock()
		
		// Verificar si necesitamos hacer espacio en el caché
		if len(ce.cache) >= ce.maxCacheSize {
			ce.evictLeastRecentlyUsed()
		}
		
		// Almacenar resultado en caché
		ce.cache[codeHash] = &CacheEntry{
			Result:      buffer.buffer,
			LastAccess:  time.Now(),
			AccessCount: 1,
		}
	}
	
	return nil
}

// hashCode genera un hash SHA-256 del código
func (ce *CachedExecutor) hashCode(code string) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))
	return hex.EncodeToString(hasher.Sum(nil))
}

// updateCacheStats actualiza las estadísticas de uso del caché
func (ce *CachedExecutor) updateCacheStats(codeHash string) {
	ce.cacheMutex.Lock()
	defer ce.cacheMutex.Unlock()
	
	if entry, exists := ce.cache[codeHash]; exists {
		entry.LastAccess = time.Now()
		entry.AccessCount++
	}
}

// evictLeastRecentlyUsed elimina la entrada menos recientemente usada del caché
func (ce *CachedExecutor) evictLeastRecentlyUsed() {
	var oldestKey string
	var oldestTime time.Time
	
	// Inicializar con el primer elemento
	for k, v := range ce.cache {
		oldestKey = k
		oldestTime = v.LastAccess
		break
	}
	
	// Encontrar la entrada más antigua
	for k, v := range ce.cache {
		if v.LastAccess.Before(oldestTime) {
			oldestKey = k
			oldestTime = v.LastAccess
		}
	}
	
	// Eliminar la entrada más antigua
	if oldestKey != "" {
		delete(ce.cache, oldestKey)
	}
}

// cleanupRoutine limpia periódicamente las entradas expiradas del caché
func (ce *CachedExecutor) cleanupRoutine() {
	ticker := time.NewTicker(ce.ttl / 2)
	defer ticker.Stop()
	
	for range ticker.C {
		ce.cleanupCache()
	}
}

// cleanupCache elimina las entradas expiradas del caché
func (ce *CachedExecutor) cleanupCache() {
	ce.cacheMutex.Lock()
	defer ce.cacheMutex.Unlock()
	
	now := time.Now()
	for k, v := range ce.cache {
		if now.Sub(v.LastAccess) > ce.ttl {
			delete(ce.cache, k)
		}
	}
}

// cachingWriter es un escritor que almacena los datos en un buffer
type cachingWriter struct {
	buffer []byte
}

// Write implementa la interfaz io.Writer
func (cw *cachingWriter) Write(p []byte) (n int, err error) {
	cw.buffer = append(cw.buffer, p...)
	return len(p), nil
}
