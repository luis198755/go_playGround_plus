package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/luis198755/go_playGround_plus/docker/pkg/executor"
	"github.com/luis198755/go_playGround_plus/docker/pkg/limiter"
	"github.com/luis198755/go_playGround_plus/docker/pkg/security"
)

// CodeRequest representa la solicitud de ejecución de código
type CodeRequest struct {
	Code string `json:"code"`
}

// Handler define el comportamiento para los manejadores HTTP
type Handler interface {
	HandleExecuteCode(w http.ResponseWriter, r *http.Request)
	HandleStaticFiles(w http.ResponseWriter, r *http.Request)
}

// APIHandler implementa los manejadores HTTP para la API
type APIHandler struct {
	limiter          limiter.RateLimiterInterface
	security         security.SecurityValidator
	executor         executor.CodeExecutor
	maxCodeLength    int
	executionTimeout int // en segundos
}

// NewAPIHandler crea un nuevo manejador de API
func NewAPIHandler(
	limiter limiter.RateLimiterInterface,
	security security.SecurityValidator,
	executor executor.CodeExecutor,
	maxCodeLength int,
	executionTimeout int,
) *APIHandler {
	return &APIHandler{
		limiter:          limiter,
		security:         security,
		executor:         executor,
		maxCodeLength:    maxCodeLength,
		executionTimeout: executionTimeout,
	}
}

// HandleExecuteCode maneja las solicitudes de ejecución de código
func (h *APIHandler) HandleExecuteCode(w http.ResponseWriter, r *http.Request) {
	// Verificar método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Rate limiting
	clientIP := h.security.GetClientIP(r)
	if !h.limiter.IsAllowed(clientIP) {
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
	h.security.SetSecurityHeaders(w)

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
	if len(codeReq.Code) > h.maxCodeLength {
		fmt.Fprintf(w, "Error: El código excede el límite de %d bytes", h.maxCodeLength)
		flusher.Flush()
		return
	}
	if hasBlacklisted, pkg := h.security.ContainsBlacklistedImports(codeReq.Code); hasBlacklisted {
		fmt.Fprintf(w, "Error: Import prohibido por seguridad: %s", pkg)
		flusher.Flush()
		return
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.executionTimeout)*time.Second)
	defer cancel()

	// Ejecutar el código
	err := h.executor.Execute(ctx, codeReq.Code, w)
	if err != nil {
		fmt.Fprintf(w, "\nError: %v", err)
		flusher.Flush()
	}
}

// FileServer representa un servidor de archivos estáticos
type FileServer struct {
	fs      http.Handler
	security security.SecurityValidator
}

// NewFileServer crea un nuevo servidor de archivos estáticos
func NewFileServer(root string, security security.SecurityValidator) *FileServer {
	return &FileServer{
		fs:      http.FileServer(http.Dir(root)),
		security: security,
	}
}

// ServeHTTP implementa la interfaz http.Handler
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.security.SetSecurityHeaders(w)
	fs.fs.ServeHTTP(w, r)
}
