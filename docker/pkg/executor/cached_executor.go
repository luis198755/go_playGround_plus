// Package executor proporciona funcionalidades para ejecutar código Go de forma segura.
//
// Este paquete implementa diferentes ejecutores de código que permiten ejecutar
// código Go en un entorno controlado, con límites de tiempo y recursos.
// También proporciona un sistema de caché para optimizar ejecuciones repetidas.
//
// Ejemplo de uso básico:
//
//     // Crear un ejecutor básico
//     baseExecutor := executor.NewGoExecutor("/usr/local/go/bin/go", 10000, "/tmp")
//
//     // Envolver con caché para optimizar ejecuciones repetidas
//     cachedExecutor := executor.NewCachedExecutor(baseExecutor, 100, 30*time.Minute)
//
//     // Ejecutar código
//     var output bytes.Buffer
//     err := cachedExecutor.Execute(context.Background(), "package main\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}", &output)
//     if err != nil {
//         log.Fatalf("Error ejecutando código: %v", err)
//     }
//     fmt.Println(output.String())
package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sync"
	"time"
)

// CacheEntry representa una entrada en el caché de ejecuciones.
// Contiene el resultado de la ejecución, la última vez que fue accedida
// y un contador de accesos para estadísticas y políticas de reemplazo.
type CacheEntry struct {
	Result      []byte
	LastAccess  time.Time
	AccessCount int
}

// CachedExecutor implementa un ejecutor con caché para código frecuentemente ejecutado.
// Utiliza un sistema de caché basado en el hash SHA-256 del código fuente para
// identificar ejecuciones idénticas y evitar la re-ejecución innecesaria.
// Incluye políticas de expiración (TTL) y reemplazo (LRU) para gestionar el tamaño del caché.
type CachedExecutor struct {
	executor     CodeExecutor
	cache        map[string]*CacheEntry
	cacheMutex   sync.RWMutex
	maxCacheSize int
	ttl          time.Duration
}

// NewCachedExecutor crea un nuevo ejecutor con caché que envuelve a otro ejecutor.
//
// Parámetros:
//   - executor: El ejecutor base que se utilizará para las ejecuciones que no estén en caché.
//   - maxCacheSize: El número máximo de entradas que se almacenarán en el caché.
//   - ttl: El tiempo de vida de las entradas en el caché antes de ser consideradas expiradas.
//
// Ejemplo:
//
//     baseExecutor := executor.NewGoExecutor("/usr/local/go/bin/go", 10000, os.TempDir())
//     cachedExecutor := executor.NewCachedExecutor(baseExecutor, 100, 30*time.Minute)
//     // Ahora cachedExecutor puede usarse como cualquier otro CodeExecutor
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

// Execute ejecuta el código Go, utilizando el caché si está disponible.
// Si el código ya ha sido ejecutado anteriormente y la entrada no ha expirado,
// devuelve el resultado almacenado en caché. De lo contrario, ejecuta el código
// utilizando el ejecutor base y almacena el resultado en el caché para futuras ejecuciones.
//
// Parámetros:
//   - ctx: Contexto para control de cancelación y timeout.
//   - code: El código Go a ejecutar.
//   - output: Writer donde se escribirá la salida de la ejecución.
//
// Retorna error si hay algún problema durante la ejecución.
//
// Ejemplo:
//
//     var output bytes.Buffer
//     err := cachedExecutor.Execute(ctx, "fmt.Println(\"Hello\");", &output)
//     if err != nil {
//         log.Printf("Error: %v", err)
//     } else {
//         fmt.Println("Resultado:", output.String())
//     }
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

// hashCode genera un hash SHA-256 del código.
// Este hash se utiliza como clave para identificar entradas únicas en el caché.
func (ce *CachedExecutor) hashCode(code string) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))
	return hex.EncodeToString(hasher.Sum(nil))
}

// updateCacheStats actualiza las estadísticas de uso del caché.
// Incrementa el contador de accesos y actualiza el timestamp de último acceso.
// Esta información se utiliza para la política de reemplazo LRU.
func (ce *CachedExecutor) updateCacheStats(codeHash string) {
	ce.cacheMutex.Lock()
	defer ce.cacheMutex.Unlock()
	
	if entry, exists := ce.cache[codeHash]; exists {
		entry.LastAccess = time.Now()
		entry.AccessCount++
	}
}

// evictLeastRecentlyUsed elimina la entrada menos recientemente usada del caché.
// Se llama cuando el caché está lleno y es necesario hacer espacio para una nueva entrada.
// Implementa la política de reemplazo Least Recently Used (LRU).
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

// cleanupRoutine limpia periódicamente las entradas expiradas del caché.
// Se ejecuta en una goroutine separada y se activa cada ttl/2 tiempo.
func (ce *CachedExecutor) cleanupRoutine() {
	ticker := time.NewTicker(ce.ttl / 2)
	defer ticker.Stop()
	
	for range ticker.C {
		ce.cleanupCache()
	}
}

// cleanupCache elimina las entradas expiradas del caché.
// Una entrada se considera expirada si ha pasado más tiempo que el TTL desde su último acceso.
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

// cachingWriter es un escritor que almacena los datos en un buffer.
// Se utiliza para capturar la salida de la ejecución y almacenarla en el caché.
type cachingWriter struct {
	buffer []byte
}

// Write implementa la interfaz io.Writer.
// Almacena los datos escritos en el buffer interno para su posterior almacenamiento en el caché.
func (cw *cachingWriter) Write(p []byte) (n int, err error) {
	cw.buffer = append(cw.buffer, p...)
	return len(p), nil
}
