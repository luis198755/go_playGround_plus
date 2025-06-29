package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/luis198755/go_playGround_plus/docker/pkg/config"
	"github.com/luis198755/go_playGround_plus/docker/pkg/executor"
	"github.com/luis198755/go_playGround_plus/docker/pkg/handlers"
	"github.com/luis198755/go_playGround_plus/docker/pkg/limiter"
	"github.com/luis198755/go_playGround_plus/docker/pkg/logger"
	"github.com/luis198755/go_playGround_plus/docker/pkg/security"
	"go.uber.org/zap"
)

// Variables globales y constantes se han movido a los paquetes correspondientes

// getEnvInt obtiene una variable de entorno int o devuelve el valor por defecto
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// Cargar configuración
	cfg := config.NewConfig()

	// Inicializar logger estructurado con nivel basado en configuración
	debugMode := cfg.DebugMode
	appLogger := logger.NewLogger(debugMode)
	appLogger.Info("Iniciando servidor Go Playground Plus", 
		zap.String("version", "1.0.0"),
		zap.String("config", cfg.String()))
	
	// Configurar variables de entorno para la ejecución del código Go
	essentialEnvVars := config.GetEssentialEnvVars()
	appLogger.Info("Configurando variables de entorno para ejecución de código")
	
	// En lugar de limpiar todas las variables de entorno (os.Clearenv),
	// establecemos solo las variables esenciales que necesitamos
	for key, value := range essentialEnvVars {
		if value != "" {
			os.Setenv(key, value)
			appLogger.Debug("Variable de entorno configurada", 
				zap.String("key", key))
		}
	}

	// Inicializar componentes
	securityValidator := security.NewCodeValidator()
	
	// Verificar que el directorio temporal existe
	if _, err := os.Stat(cfg.TempDir); os.IsNotExist(err) {
		appLogger.Info("Creando directorio temporal", zap.String("dir", cfg.TempDir))
		if err := os.MkdirAll(cfg.TempDir, 0755); err != nil {
			appLogger.Fatal("Error al crear directorio temporal", zap.Error(err))
		}
	}
	
	// Inicializar rate limiter con configuración
	rateLimiter := limiter.NewRateLimiter(cfg.MaxRequestsPerMinute)
	appLogger.Info("Rate limiter configurado", 
		zap.Int("max_requests_per_minute", cfg.MaxRequestsPerMinute))
	
	// Inicializar ejecutor de código Go
	baseExecutor := executor.NewGoExecutor(
		cfg.GoExecutablePath,
		cfg.MaxOutputLength,
		cfg.TempDir,
	)
	
	// Configurar el ejecutor con caché
	maxCacheSize := getEnvInt("MAX_CACHE_SIZE", 100) // Número máximo de entradas en caché
	cacheTTL := time.Duration(getEnvInt("CACHE_TTL_MINUTES", 30)) * time.Minute
	
	appLogger.Info("Configurando caché de ejecución", 
		zap.Int("max_size", maxCacheSize),
		zap.Duration("ttl", cacheTTL))
		
	codeExecutor := executor.NewCachedExecutor(baseExecutor, maxCacheSize, cacheTTL)
	appLogger.Info("Ejecutor de código configurado", 
		zap.String("go_path", cfg.GoExecutablePath),
		zap.String("temp_dir", cfg.TempDir))
	
	// Inicializar handlers
	apiHandler := handlers.NewAPIHandler(
		rateLimiter,
		securityValidator,
		codeExecutor,
		appLogger,
		cfg.MaxCodeLength,
		cfg.ExecutionTimeout,
	)
	
	// Configurar rutas
	http.HandleFunc("/api/execute", apiHandler.HandleExecuteCode)
	
	// Servir archivos estáticos desde la ruta configurada
	staticDir := cfg.StaticFilesDir
	appLogger.Info("Configurando servidor de archivos estáticos", 
		zap.String("static_dir", staticDir))
	
	// Verificar que el directorio de archivos estáticos exista
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		appLogger.Error("El directorio de archivos estáticos no existe", 
			zap.String("static_dir", staticDir),
			zap.Error(err))
		// Intentar crear el directorio
		if err := os.MkdirAll(staticDir, 0755); err != nil {
			appLogger.Fatal("No se pudo crear el directorio de archivos estáticos", 
				zap.String("static_dir", staticDir),
				zap.Error(err))
		}
		appLogger.Info("Directorio de archivos estáticos creado", 
			zap.String("static_dir", staticDir))
	}
	
	fileServer := handlers.NewFileServer(staticDir, securityValidator)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := securityValidator.GetClientIP(r)
		appLogger.Info("Petición recibida", 
			zap.String("ip", clientIP),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))

		path := filepath.Join(staticDir, r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			appLogger.Info("Archivo no encontrado, sirviendo index.html", 
				zap.String("ip", clientIP),
				zap.String("path", r.URL.Path))
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}
		appLogger.Info("Sirviendo archivo", 
			zap.String("ip", clientIP),
			zap.String("path", r.URL.Path))
		fileServer.ServeHTTP(w, r)
	})

	// Iniciar servidor
	serverAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	appLogger.Info("Servidor iniciado", 
		zap.String("address", serverAddr),
		zap.String("static_dir", staticDir))
	
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		appLogger.Fatal("Error al iniciar el servidor", 
			zap.String("address", serverAddr),
			zap.Error(err))
	}
}
