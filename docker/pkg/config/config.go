package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	// Configuración del servidor
	Port                string
	Host                string
	DebugMode          bool
	StaticFilesDir     string

	// Límites y seguridad
	MaxRequestsPerMinute int
	MaxCodeLength        int
	MaxOutputLength      int
	ExecutionTimeout     time.Duration
	AllowedOrigins       []string

	// Ejecución de código Go
	GoExecutablePath     string
	TempDir              string
	CleanupInterval      time.Duration

	// Logging
	LogLevel            string
	LogFormat           string
}

// NewConfig crea una nueva configuración con valores por defecto
// y los sobrescribe con variables de entorno si están disponibles
func NewConfig() *Config {
	// Valores por defecto
	cfg := &Config{
		// Configuración del servidor
		Port:            getEnvString("SERVER_PORT", "8080"),
		Host:            getEnvString("SERVER_HOST", "0.0.0.0"),
		DebugMode:       getEnvBool("DEBUG_MODE", false),
		StaticFilesDir:  getEnvString("STATIC_FILES_DIR", "/app/build"),

		// Límites y seguridad
		MaxRequestsPerMinute: getEnvInt("MAX_REQUESTS_PER_MINUTE", 30),
		MaxCodeLength:        getEnvInt("MAX_CODE_LENGTH", 10000),
		MaxOutputLength:      getEnvInt("MAX_OUTPUT_LENGTH", 10000),
		ExecutionTimeout:     time.Duration(getEnvInt("EXECUTION_TIMEOUT_SECONDS", 10)) * time.Second,
		AllowedOrigins:       getEnvStringSlice("ALLOWED_ORIGINS", []string{"*"}),

		// Ejecución de código Go
		GoExecutablePath: getEnvString("GO_EXECUTABLE_PATH", "/usr/local/go/bin/go"),
		TempDir:          getEnvString("TEMP_DIR", os.TempDir()),
		CleanupInterval:  time.Duration(getEnvInt("CLEANUP_INTERVAL_MINUTES", 60)) * time.Minute,

		// Logging
		LogLevel:  getEnvString("LOG_LEVEL", "info"),
		LogFormat: getEnvString("LOG_FORMAT", "json"),
	}

	// Validación de la configuración
	validateConfig(cfg)

	return cfg
}

// Funciones auxiliares para obtener valores de variables de entorno

// getEnvString obtiene una variable de entorno string o devuelve el valor por defecto
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt obtiene una variable de entorno int o devuelve el valor por defecto
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool obtiene una variable de entorno bool o devuelve el valor por defecto
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		value = strings.ToLower(value)
		return value == "true" || value == "1" || value == "yes" || value == "y"
	}
	return defaultValue
}

// getEnvStringSlice obtiene una variable de entorno como slice de strings o devuelve el valor por defecto
func getEnvStringSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// validateConfig valida la configuración y ajusta valores si es necesario
func validateConfig(cfg *Config) {
	// Validar límites mínimos
	if cfg.MaxRequestsPerMinute < 1 {
		cfg.MaxRequestsPerMinute = 1
		fmt.Println("WARNING: MAX_REQUESTS_PER_MINUTE ajustado a valor mínimo de 1")
	}

	if cfg.MaxCodeLength < 100 {
		cfg.MaxCodeLength = 100
		fmt.Println("WARNING: MAX_CODE_LENGTH ajustado a valor mínimo de 100")
	}

	if cfg.ExecutionTimeout < time.Second {
		cfg.ExecutionTimeout = time.Second
		fmt.Println("WARNING: EXECUTION_TIMEOUT_SECONDS ajustado a valor mínimo de 1 segundo")
	}

	// Validar que el directorio temporal exista o se pueda crear
	if cfg.TempDir != "" {
		if _, err := os.Stat(cfg.TempDir); os.IsNotExist(err) {
			err := os.MkdirAll(cfg.TempDir, 0755)
			if err != nil {
				fmt.Printf("ERROR: No se pudo crear el directorio temporal %s: %v\n", cfg.TempDir, err)
				cfg.TempDir = os.TempDir()
			}
		}
	}

	// Validar que el ejecutable de Go exista
	if _, err := os.Stat(cfg.GoExecutablePath); os.IsNotExist(err) {
		fmt.Printf("WARNING: El ejecutable de Go no existe en %s\n", cfg.GoExecutablePath)
	}
}

// GetEssentialEnvVars devuelve un mapa con las variables de entorno esenciales
// para la ejecución de código Go
func GetEssentialEnvVars() map[string]string {
	return map[string]string{
		"HOME":           os.Getenv("HOME"),
		"PATH":           os.Getenv("PATH"),
		"GOCACHE":        os.Getenv("GOCACHE"),
		"XDG_CACHE_HOME": os.Getenv("XDG_CACHE_HOME"),
		"GOPATH":         os.Getenv("GOPATH"),
		"GOROOT":         os.Getenv("GOROOT"),
		"PORT":           os.Getenv("SERVER_PORT"),
	}
}

// String devuelve una representación en string de la configuración
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{Port: %s, Host: %s, DebugMode: %v, MaxReqPerMin: %d, MaxCodeLen: %d, ExecTimeout: %v, LogLevel: %s}",
		c.Port, c.Host, c.DebugMode, c.MaxRequestsPerMinute, c.MaxCodeLength, c.ExecutionTimeout, c.LogLevel,
	)
}
