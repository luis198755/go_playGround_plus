package config

import (
	"os"
	"strconv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Port                string
	MaxRequestsPerMinute int
	MaxCodeLength        int
	MaxOutputLength      int
	GoExecutablePath     string
	TempDir             string
	ExecutionTimeout    int // en segundos
}

// NewConfig crea una nueva configuración con valores por defecto
// y los sobrescribe con variables de entorno si están disponibles
func NewConfig() *Config {
	cfg := &Config{
		Port:                "8080",
		MaxRequestsPerMinute: 30,
		MaxCodeLength:        10000,
		MaxOutputLength:      10000,
		GoExecutablePath:     "/usr/local/go/bin/go",
		TempDir:             "",
		ExecutionTimeout:    10,
	}

	// Sobrescribir con variables de entorno
	if port := os.Getenv("WEB_PORT"); port != "" {
		cfg.Port = port
	}

	if maxReq := os.Getenv("MAX_REQUESTS_PER_MINUTE"); maxReq != "" {
		if val, err := strconv.Atoi(maxReq); err == nil {
			cfg.MaxRequestsPerMinute = val
		}
	}

	if maxCode := os.Getenv("MAX_CODE_LENGTH"); maxCode != "" {
		if val, err := strconv.Atoi(maxCode); err == nil {
			cfg.MaxCodeLength = val
		}
	}

	if maxOutput := os.Getenv("MAX_OUTPUT_LENGTH"); maxOutput != "" {
		if val, err := strconv.Atoi(maxOutput); err == nil {
			cfg.MaxOutputLength = val
		}
	}

	if goPath := os.Getenv("GO_EXECUTABLE_PATH"); goPath != "" {
		cfg.GoExecutablePath = goPath
	}

	if tempDir := os.Getenv("TEMP_DIR"); tempDir != "" {
		cfg.TempDir = tempDir
	}

	if timeout := os.Getenv("EXECUTION_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			cfg.ExecutionTimeout = val
		}
	}

	return cfg
}

// GetEssentialEnvVars devuelve un mapa con las variables de entorno esenciales
func GetEssentialEnvVars() map[string]string {
	return map[string]string{
		"HOME":           os.Getenv("HOME"),
		"PATH":           os.Getenv("PATH"),
		"GOCACHE":        os.Getenv("GOCACHE"),
		"XDG_CACHE_HOME": os.Getenv("XDG_CACHE_HOME"),
		"GOPATH":         os.Getenv("GOPATH"),
		"GOROOT":         os.Getenv("GOROOT"),
		"PORT":           os.Getenv("WEB_PORT"),
	}
}
