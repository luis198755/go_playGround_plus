package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/luis198755/go_playGround_plus/docker/pkg/config"
	"github.com/luis198755/go_playGround_plus/docker/pkg/executor"
	"github.com/luis198755/go_playGround_plus/docker/pkg/handlers"
	"github.com/luis198755/go_playGround_plus/docker/pkg/limiter"
	"github.com/luis198755/go_playGround_plus/docker/pkg/logger"
	"github.com/luis198755/go_playGround_plus/docker/pkg/security"
	"go.uber.org/zap"
)

// Variables globales y constantes se han movido a los paquetes correspondientes

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// Inicializar logger estructurado
	appLogger := logger.NewLogger(true) // true para modo desarrollo
	appLogger.Info("Iniciando servidor Go Playground Plus", zap.String("version", "1.0.0"))

	// Cargar configuración
	cfg := config.NewConfig()
	
	// Configurar variables de entorno
	essentialEnvVars := config.GetEssentialEnvVars()
	os.Clearenv()
	for key, value := range essentialEnvVars {
		if value != "" {
			os.Setenv(key, value)
		}
	}

	// Inicializar componentes
	securityValidator := security.NewCodeValidator()
	rateLimiter := limiter.NewRateLimiter(cfg.MaxRequestsPerMinute)
	codeExecutor := executor.NewGoExecutor(cfg.GoExecutablePath, cfg.MaxOutputLength, cfg.TempDir)
	
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
	
	// Servir archivos estáticos
	fileServer := handlers.NewFileServer("build", securityValidator)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := securityValidator.GetClientIP(r)
		appLogger.Info("Petición recibida", 
			zap.String("ip", clientIP),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))

		path := filepath.Join("build", r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			appLogger.Info("Archivo no encontrado, sirviendo index.html", 
				zap.String("ip", clientIP),
				zap.String("path", r.URL.Path))
			http.ServeFile(w, r, "build/index.html")
			return
		}
		log.Printf("[IP: %s] Sirviendo archivo: %s", clientIP, r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})

	// Iniciar servidor
	appLogger.Info("Servidor iniciado", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
