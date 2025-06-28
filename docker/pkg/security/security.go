package security

import (
	"net/http"
	"regexp"
	"strings"
)

// SecurityValidator define el comportamiento para validaciones de seguridad
type SecurityValidator interface {
	ContainsBlacklistedImports(code string) (bool, string)
	GetClientIP(r *http.Request) string
	SetSecurityHeaders(w http.ResponseWriter)
}

// CodeValidator implementa validaciones de seguridad para código Go
type CodeValidator struct {
	blacklistedImports []string
	importPattern      *regexp.Regexp
}

// NewCodeValidator crea un nuevo validador de código
func NewCodeValidator() *CodeValidator {
	return &CodeValidator{
		blacklistedImports: []string{
			"os/exec",
			"syscall",
			"unsafe",
			"net",
			"net/http",
			"plugin",
		},
		importPattern: regexp.MustCompile(`(?m)^\s*import\s*(\((?:[^)]+)\)|"[^"]+")`),
	}
}

// ContainsBlacklistedImports verifica si el código contiene imports prohibidos
func (cv *CodeValidator) ContainsBlacklistedImports(code string) (bool, string) {
	// Buscar todos los matches de imports en el código
	matches := cv.importPattern.FindAllStringSubmatch(code, -1)
	
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
			for _, blacklisted := range cv.blacklistedImports {
				if imp == blacklisted {
					return true, blacklisted
				}
			}
		}
	}
	return false, ""
}

// GetClientIP obtiene la dirección IP del cliente desde la solicitud HTTP
func (cv *CodeValidator) GetClientIP(r *http.Request) string {
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

// SetSecurityHeaders establece los encabezados de seguridad en la respuesta HTTP
func (cv *CodeValidator) SetSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}
