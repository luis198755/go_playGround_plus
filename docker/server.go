package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"regexp"
)

type CodeRequest struct {
	Code string `json:"code"`
}

var blacklistedImports = []string{
	"os/exec",
	"syscall",
	"unsafe",
	"net",
	"net/http",
	"plugin",
}

type RateLimiter struct {
	requests map[string]*IPRequests
	mu       sync.RWMutex
}

type IPRequests struct {
	count    int
	lastTime time.Time
}

var rateLimiter = &RateLimiter{
	requests: make(map[string]*IPRequests),
}

const (
	maxRequestsPerMinute = 30
	maxCodeLength        = 10000 // 5KB
	maxOutputLength      = 10000 // 5KB
)

// Expresión regular mejorada para detectar imports en bloque o línea única
var importPattern = regexp.MustCompile(`(?m)^\s*import\s*(\((?:[^)]+)\)|"[^"]+")`)

func containsBlacklistedImports(code string) (bool, string) {
	// Buscar todos los matches de imports en el código
	matches := importPattern.FindAllStringSubmatch(code, -1)
	
	for _, match := range matches {
		importStatement := match[1] // Captura lo que está dentro de `import (...)` o `import "..."`

		// Eliminar paréntesis si es un bloque
		importStatement = strings.ReplaceAll(importStatement, "(", "")
		importStatement = strings.ReplaceAll(importStatement, ")", "")

		// Separar los imports en líneas individuales y limpiar espacios
		imports := strings.Split(importStatement, "\n")
		for _, imp := range imports {
			imp = strings.TrimSpace(strings.Split(imp, "//")[0]) // Eliminar comentarios en línea
			imp = strings.Trim(imp, `"`)                         // Eliminar comillas si existen

			// Comparar con la lista de imports prohibidos
			for _, blacklisted := range blacklistedImports {
				if imp == blacklisted {
					return true, blacklisted
				}
			}
		}
	}
	return false, ""
}

func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if req, exists := rl.requests[ip]; exists {
		if now.Sub(req.lastTime) > time.Minute {
			req.count = 1
			req.lastTime = now
			return true
		}
		if req.count < maxRequestsPerMinute {
			req.count++
			return true
		}
		return false
	}
	rl.requests[ip] = &IPRequests{count: 1, lastTime: now}
	return true
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

// handleExecuteCode ahora envía la salida del comando en streaming
func handleExecuteCode(w http.ResponseWriter, r *http.Request) {
	// Verificar método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Rate limiting
	clientIP := getClientIP(r)
	if !rateLimiter.isAllowed(clientIP) {
		log.Printf("[SECURITY] Rate limit exceeded for IP: %s", clientIP)
		http.Error(w, "Demasiadas peticiones. Por favor, espere un minuto.", http.StatusTooManyRequests)
		return
	}

	// Verificar Content-Type
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Content-Type debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Establecer headers de seguridad y para streaming
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Verificar que el ResponseWriter soporte flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming no soportado", http.StatusInternalServerError)
		return
	}

	// Decodificar la solicitud
	var codeReq CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&codeReq); err != nil {
		log.Printf("Error al decodificar la solicitud: %v", err)
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	// Validar el código
	if codeReq.Code == "" {
		fmt.Fprint(w, "Error: El código no puede estar vacío")
		flusher.Flush()
		return
	}
	if len(codeReq.Code) > maxCodeLength {
		fmt.Fprintf(w, "Error: El código excede el límite de %d bytes", maxCodeLength)
		flusher.Flush()
		return
	}
	if hasBlacklisted, pkg := containsBlacklistedImports(codeReq.Code); hasBlacklisted {
		fmt.Fprintf(w, "Error: Import prohibido por seguridad: %s", pkg)
		flusher.Flush()
		return
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Crear archivo temporal para el código
	tmpFile, err := os.CreateTemp("", "code-*.go")
	if err != nil {
		fmt.Fprint(w, "Error creando archivo temporal: ", err.Error())
		flusher.Flush()
		return
	}
	tmpPath := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		// Intentar eliminar el archivo temporal
		for i := 0; i < 3; i++ {
			if err := os.Remove(tmpPath); err == nil || os.IsNotExist(err) {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	if _, err := tmpFile.WriteString(codeReq.Code); err != nil {
		fmt.Fprint(w, "Error escribiendo código: ", err.Error())
		flusher.Flush()
		return
	}
	tmpFile.Close()

	// Configurar y ejecutar el comando
	cmd := exec.CommandContext(ctx, "/usr/local/go/bin/go", "run", tmpPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprint(w, "Error obteniendo salida del comando: ", err.Error())
		flusher.Flush()
		return
	}
	// Combinar stderr con stdout
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		fmt.Fprint(w, "Error iniciando el comando: ", err.Error())
		flusher.Flush()
		return
	}

	totalBytes := 0
	buf := make([]byte, 1024)
	for {
		n, err := stdoutPipe.Read(buf)
		if n > 0 {
			// Limitar la cantidad total de bytes enviados
			if totalBytes+n > maxOutputLength {
				allowed := maxOutputLength - totalBytes
				if allowed > 0 {
					w.Write(buf[:allowed])
					totalBytes += allowed
				}
				fmt.Fprint(w, "\n... (output truncated)")
				flusher.Flush()
				break
			} else {
				w.Write(buf[:n])
				totalBytes += n
				flusher.Flush()
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("Error leyendo salida: %v", err)
			}
			break
		}
	}

	// Esperar a que el comando finalice
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(w, "\nError: %v", err)
		flusher.Flush()
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	essentialEnvVars := map[string]string{
		"HOME":           os.Getenv("HOME"),
		"PATH":           os.Getenv("PATH"),
		"GOCACHE":        os.Getenv("GOCACHE"),
		"XDG_CACHE_HOME": os.Getenv("XDG_CACHE_HOME"),
		"GOPATH":         os.Getenv("GOPATH"),
		"GOROOT":         os.Getenv("GOROOT"),
		"PORT":           os.Getenv("WEB_PORT"),
	}

	os.Clearenv()
	for key, value := range essentialEnvVars {
		if value != "" {
			os.Setenv(key, value)
		}
	}

	port := essentialEnvVars["PORT"]
	if port == "" {
		port = "8080"
	}

	fs := http.FileServer(http.Dir("build"))
	http.HandleFunc("/api/execute", handleExecuteCode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)
		log.Printf("[IP: %s] Recibida petición: %s %s", clientIP, r.Method, r.URL.Path)

		path := filepath.Join("build", r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			log.Printf("[IP: %s] Archivo no encontrado: %s, sirviendo index.html", clientIP, r.URL.Path)
			http.ServeFile(w, r, "build/index.html")
			return
		}
		log.Printf("[IP: %s] Sirviendo archivo: %s", clientIP, r.URL.Path)
		fs.ServeHTTP(w, r)
	})

	log.Printf("Servidor iniciado en puerto :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
