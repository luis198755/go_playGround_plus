// Package config proporciona funcionalidades para la configuración de la aplicación Go Playground Plus.
//
// Este paquete maneja la carga de configuración desde variables de entorno con valores por defecto,
// validación de configuración y gestión de variables de entorno esenciales para la ejecución de código Go.
//
// Ejemplo de uso básico:
//
//     // Cargar configuración desde variables de entorno con valores por defecto
//     cfg := config.NewConfig()
//
//     // Utilizar la configuración
//     fmt.Printf("Servidor escuchando en %s:%s\n", cfg.Host, cfg.Port)
//     fmt.Printf("Timeout de ejecución: %v\n", cfg.ExecutionTimeout)
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contiene toda la configuración de la aplicación Go Playground Plus.
//
// Esta estructura agrupa todas las opciones de configuración organizadas por categorías:
// - Configuración del servidor (puerto, host, modo debug, directorio de archivos estáticos)
// - Límites y seguridad (rate limiting, tamaño máximo de código, timeout de ejecución)
// - Ejecución de código Go (ruta del ejecutable, directorio temporal, intervalo de limpieza)
// - Logging (nivel y formato)
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
// y los sobrescribe con variables de entorno si están disponibles.
//
// Este método carga todas las opciones de configuración desde variables de entorno,
// utilizando valores por defecto cuando no están definidas. También realiza validaciones
// para asegurar que la configuración sea válida y segura.
//
// Retorna un puntero a una estructura Config completamente inicializada.
//
// Ejemplo:
//
//     // Establecer algunas variables de entorno para personalizar la configuración
//     os.Setenv("SERVER_PORT", "9000")
//     os.Setenv("DEBUG_MODE", "true")
//
//     // Cargar la configuración
//     cfg := config.NewConfig()
//
//     // La configuración tendrá SERVER_PORT="9000" y DEBUG_MODE=true,
//     // mientras que el resto de opciones tendrán sus valores por defecto
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
// Estas funciones facilitan la obtención de valores tipados desde variables de entorno,
// proporcionando valores por defecto cuando la variable no está definida o su valor no es válido.

// getEnvString obtiene una variable de entorno string o devuelve el valor por defecto.
//
// Parámetros:
//   - key: Nombre de la variable de entorno.
//   - defaultValue: Valor por defecto a utilizar si la variable no existe o está vacía.
//
// Retorna el valor de la variable de entorno o el valor por defecto.
//
// Ejemplo:
//
//     port := getEnvString("SERVER_PORT", "8080")
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt obtiene una variable de entorno int o devuelve el valor por defecto.
//
// Parámetros:
//   - key: Nombre de la variable de entorno.
//   - defaultValue: Valor por defecto a utilizar si la variable no existe o no es un entero válido.
//
// Retorna el valor de la variable de entorno convertido a entero o el valor por defecto.
//
// Ejemplo:
//
//     maxRequests := getEnvInt("MAX_REQUESTS_PER_MINUTE", 30)
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool obtiene una variable de entorno bool o devuelve el valor por defecto.
//
// Parámetros:
//   - key: Nombre de la variable de entorno.
//   - defaultValue: Valor por defecto a utilizar si la variable no existe o no es un booleano válido.
//
// Retorna el valor de la variable de entorno convertido a booleano o el valor por defecto.
// Los valores "true", "1", "yes" y "y" (case-insensitive) se consideran verdaderos.
//
// Ejemplo:
//
//     debugMode := getEnvBool("DEBUG_MODE", false)
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		value = strings.ToLower(value)
		return value == "true" || value == "1" || value == "yes" || value == "y"
	}
	return defaultValue
}

// getEnvStringSlice obtiene una variable de entorno como slice de strings o devuelve el valor por defecto.
//
// Parámetros:
//   - key: Nombre de la variable de entorno.
//   - defaultValue: Valor por defecto a utilizar si la variable no existe o está vacía.
//
// Retorna el valor de la variable de entorno dividido por comas como slice de strings,
// o el valor por defecto si la variable no existe.
//
// Ejemplo:
//
//     // Con ALLOWED_ORIGINS="http://localhost:3000,https://example.com"
//     origins := getEnvStringSlice("ALLOWED_ORIGINS", []string{"*"})
//     // origins = ["http://localhost:3000", "https://example.com"]
func getEnvStringSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// validateConfig valida la configuración y ajusta valores si es necesario.
//
// Esta función realiza comprobaciones de seguridad y validez en la configuración,
// como asegurar que los límites no sean demasiado bajos o altos, verificar la existencia
// de directorios y ejecutables, etc.
//
// Parámetros:
//   - cfg: Puntero a la estructura Config a validar.
//
// La función modifica la estructura Config in-place si es necesario realizar ajustes.
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
// para la ejecución de código Go.
//
// Esta función recopila las variables de entorno que deben estar disponibles
// durante la ejecución de código Go, como PATH, GOPATH, GOROOT, etc.
//
// Retorna un mapa de strings con las variables de entorno esenciales.
//
// Ejemplo:
//
//     envVars := config.GetEssentialEnvVars()
//     for key, value := range envVars {
//         os.Setenv(key, value)
//     }
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

// String devuelve una representación en string de la configuración.
//
// Este método implementa la interfaz Stringer para facilitar el logging
// y depuración de la configuración.
//
// Retorna una cadena con la representación JSON de la configuración.
//
// Ejemplo:
//
//     cfg := config.NewConfig()
//     fmt.Println(cfg.String())
//     // Imprime: {"Port":"8080","Host":"0.0.0.0",...}
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{Port: %s, Host: %s, DebugMode: %v, MaxReqPerMin: %d, MaxCodeLen: %d, ExecTimeout: %v, LogLevel: %s}",
		c.Port, c.Host, c.DebugMode, c.MaxRequestsPerMinute, c.MaxCodeLength, c.ExecutionTimeout, c.LogLevel,
	)
}
