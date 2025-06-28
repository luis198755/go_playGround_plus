package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/luis198755/go_playGround_plus/docker/pkg/errors"
	"github.com/luis198755/go_playGround_plus/docker/pkg/executor"
	"github.com/luis198755/go_playGround_plus/docker/pkg/limiter"
	"github.com/luis198755/go_playGround_plus/docker/pkg/logger"
	"github.com/luis198755/go_playGround_plus/docker/pkg/security"
	"go.uber.org/zap"
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
	logger           logger.Logger
	maxCodeLength    int
	executionTimeout int // en segundos
}

// NewAPIHandler crea un nuevo manejador de API
func NewAPIHandler(
	limiter limiter.RateLimiterInterface,
	security security.SecurityValidator,
	executor executor.CodeExecutor,
	log logger.Logger,
	maxCodeLength int,
	executionTimeout int,
) *APIHandler {
	return &APIHandler{
		limiter:          limiter,
		security:         security,
		executor:         executor,
		logger:           log,
		maxCodeLength:    maxCodeLength,
		executionTimeout: executionTimeout,
	}
}

// HandleExecuteCode maneja las solicitudes de ejecución de código
func (h *APIHandler) HandleExecuteCode(w http.ResponseWriter, r *http.Request) {
	// Crear logger con contexto para esta solicitud
	reqLogger := h.logger.With(
		zap.String("client_ip", h.security.GetClientIP(r)),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	// Verificar método HTTP
	if r.Method != http.MethodPost {
		err := errors.WithContext(
			errors.New("método no permitido"),
			http.StatusMethodNotAllowed,
			"Método no permitido",
			map[string]interface{}{"method": r.Method},
		)
		errors.HTTPError(w, r, reqLogger, err)
		return
	}

	// Rate limiting
	clientIP := h.security.GetClientIP(r)
	if !h.limiter.IsAllowed(clientIP) {
		reqLogger.Warn("Rate limit exceeded",
			zap.String("client_ip", clientIP),
		)
		err := errors.TooManyRequests(
			errors.New("rate limit exceeded"),
			"Demasiadas peticiones. Por favor, espere un minuto.",
			map[string]interface{}{"client_ip": clientIP},
		)
		errors.HTTPError(w, r, reqLogger, err)
		return
	}

	// Verificar Content-Type
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		err := errors.BadRequest(
			errors.New("content-type inválido"),
			"Content-Type debe ser application/json",
			map[string]interface{}{"content_type": r.Header.Get("Content-Type")},
		)
		errors.HTTPError(w, r, reqLogger, err)
		return
	}

	// Establecer headers de seguridad y para streaming
	h.security.SetSecurityHeaders(w)

	// Verificar que el ResponseWriter soporte flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		err := errors.InternalServerError(
			errors.New("streaming no soportado"),
			"El servidor no soporta streaming de respuestas",
			nil,
		)
		errors.HTTPError(w, r, reqLogger, err)
		return
	}

	// Decodificar la solicitud
	var codeReq CodeRequest
	// Asegurar que el body se cierre adecuadamente
	defer r.Body.Close()
	
	if err := json.NewDecoder(r.Body).Decode(&codeReq); err != nil {
		reqLogger.Error("Error al decodificar la solicitud", zap.Error(err))
		err := errors.BadRequest(
			errors.Wrap(err, "error al decodificar JSON"),
			"Solicitud inválida",
			nil,
		)
		errors.HTTPError(w, r, reqLogger, err)
		return
	}

	// Validar el código
	if codeReq.Code == "" {
		reqLogger.Warn("Código vacío recibido")
		fmt.Fprint(w, "Error: El código no puede estar vacío")
		flusher.Flush()
		return
	}

	if len(codeReq.Code) > h.maxCodeLength {
		reqLogger.Warn("Código excede límite de tamaño",
			zap.Int("code_length", len(codeReq.Code)),
			zap.Int("max_length", h.maxCodeLength),
		)
		fmt.Fprintf(w, "Error: El código excede el límite de %d bytes", h.maxCodeLength)
		flusher.Flush()
		return
	}

	if hasBlacklisted, pkg := h.security.ContainsBlacklistedImports(codeReq.Code); hasBlacklisted {
		reqLogger.Warn("Intento de usar import prohibido",
			zap.String("blacklisted_package", pkg),
		)
		fmt.Fprintf(w, "Error: Import prohibido por seguridad: %s", pkg)
		flusher.Flush()
		return
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.executionTimeout)*time.Second)
	defer cancel()

	// Registrar ejecución
	reqLogger.Info("Ejecutando código Go",
		zap.Int("code_length", len(codeReq.Code)),
		zap.Int("timeout_seconds", h.executionTimeout),
	)

	// Ejecutar el código
	err := h.executor.Execute(ctx, codeReq.Code, w)
	if err != nil {
		reqLogger.Error("Error al ejecutar código", 
			zap.Error(errors.Wrap(err, "error de ejecución")),
		)
		fmt.Fprintf(w, "\nError: %v", err)
		flusher.Flush()
	} else {
		reqLogger.Info("Código ejecutado correctamente")
	}
}

// FileServer representa un servidor de archivos estáticos
type FileServer struct {
	fs      http.Handler
	security security.SecurityValidator
	root     string
}

// NewFileServer crea un nuevo servidor de archivos estáticos
func NewFileServer(root string, security security.SecurityValidator) *FileServer {
	return &FileServer{
		fs:      http.FileServer(http.Dir(root)),
		security: security,
		root:     root,
	}
}

// ServeHTTP implementa la interfaz http.Handler
func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Establecer encabezados de seguridad
	fs.security.SetSecurityHeaders(w)
	
	// Establecer el tipo de contenido correcto según la extensión del archivo
	path := r.URL.Path
	if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	} else if strings.HasSuffix(path, ".svg") {
		w.Header().Set("Content-Type", "image/svg+xml")
	} else if strings.HasSuffix(path, ".html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	
	// Servir el archivo
	fs.fs.ServeHTTP(w, r)
}
