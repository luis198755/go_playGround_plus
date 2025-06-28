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
	"github.com/luis198755/go_playGround_plus/docker/pkg/security"
)

// Variables globales y constantes se han movido a los paquetes correspondientes

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

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
		cfg.MaxCodeLength,
		cfg.ExecutionTimeout,
	)
	
	// Configurar rutas
	http.HandleFunc("/api/execute", apiHandler.HandleExecuteCode)
	
	// Servir archivos estáticos
	fileServer := handlers.NewFileServer("build", securityValidator)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := securityValidator.GetClientIP(r)
		log.Printf("[IP: %s] Recibida petición: %s %s", clientIP, r.Method, r.URL.Path)

		path := filepath.Join("build", r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			log.Printf("[IP: %s] Archivo no encontrado: %s, sirviendo index.html", clientIP, r.URL.Path)
			http.ServeFile(w, r, "build/index.html")
			return
		}
		log.Printf("[IP: %s] Sirviendo archivo: %s", clientIP, r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})

	// Iniciar servidor
	log.Printf("Servidor iniciado en puerto :%s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
