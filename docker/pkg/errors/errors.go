package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luis198755/go_playGround_plus/docker/pkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// AppError representa un error de la aplicación con contexto adicional
type AppError struct {
	Err        error
	StatusCode int
	Message    string
	Context    map[string]interface{}
}

// Error implementa la interfaz error
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// Unwrap devuelve el error original
func (e *AppError) Unwrap() error {
	return e.Err
}

// ErrorResponse es la estructura que se envía como respuesta HTTP en caso de error
type ErrorResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// New crea un nuevo error con contexto
func New(message string) error {
	return errors.New(message)
}

// Wrap envuelve un error con un mensaje adicional
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf envuelve un error con un mensaje formateado
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// WithContext añade contexto a un error
func WithContext(err error, statusCode int, message string, context map[string]interface{}) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
		Context:    context,
	}
}

// IsNotFound verifica si un error es de tipo "no encontrado"
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsBadRequest verifica si un error es de tipo "solicitud incorrecta"
func IsBadRequest(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode == http.StatusBadRequest
	}
	return false
}

// IsUnauthorized verifica si un error es de tipo "no autorizado"
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsForbidden verifica si un error es de tipo "prohibido"
func IsForbidden(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode == http.StatusForbidden
	}
	return false
}

// HTTPError responde con un error HTTP y registra el error
func HTTPError(w http.ResponseWriter, r *http.Request, log logger.Logger, err error) {
	var appErr *AppError
	statusCode := http.StatusInternalServerError
	message := "Error interno del servidor"
	details := make(map[string]interface{})

	if errors.As(err, &appErr) {
		statusCode = appErr.StatusCode
		message = appErr.Message
		details = appErr.Context
	}

	// Registrar el error con contexto
	log.Error("Error HTTP",
		zap.Int("status_code", statusCode),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.Error(err),
	)

	// Crear respuesta de error
	resp := ErrorResponse{
		Status:  statusCode,
		Message: message,
		Details: details,
	}

	// Enviar respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("Error al codificar respuesta JSON", zap.Error(err))
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}
}

// BadRequest crea un error de tipo "solicitud incorrecta"
func BadRequest(err error, message string, context map[string]interface{}) *AppError {
	return WithContext(err, http.StatusBadRequest, message, context)
}

// NotFound crea un error de tipo "no encontrado"
func NotFound(err error, message string, context map[string]interface{}) *AppError {
	return WithContext(err, http.StatusNotFound, message, context)
}

// Unauthorized crea un error de tipo "no autorizado"
func Unauthorized(err error, message string, context map[string]interface{}) *AppError {
	return WithContext(err, http.StatusUnauthorized, message, context)
}

// Forbidden crea un error de tipo "prohibido"
func Forbidden(err error, message string, context map[string]interface{}) *AppError {
	return WithContext(err, http.StatusForbidden, message, context)
}

// InternalServerError crea un error de tipo "error interno del servidor"
func InternalServerError(err error, message string, context map[string]interface{}) *AppError {
	return WithContext(err, http.StatusInternalServerError, message, context)
}

// TooManyRequests crea un error de tipo "demasiadas solicitudes"
func TooManyRequests(err error, message string, context map[string]interface{}) *AppError {
	return WithContext(err, http.StatusTooManyRequests, message, context)
}
