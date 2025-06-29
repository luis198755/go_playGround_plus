// Package executor proporciona funcionalidades para ejecutar código Go de forma segura.
//
// Este paquete implementa diferentes ejecutores de código que permiten ejecutar
// código Go en un entorno controlado, con límites de tiempo y recursos.
// También proporciona un sistema de caché para optimizar ejecuciones repetidas.
package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// CodeExecutor define la interfaz para ejecutar código Go.
//
// Esta interfaz permite implementar diferentes estrategias de ejecución de código,
// como ejecución directa, con caché, con sandbox, etc., manteniendo una API consistente.
//
// Ejemplo de uso:
//
//     var executor CodeExecutor = NewGoExecutor("/usr/local/go/bin/go", 10000, os.TempDir())
//     var output bytes.Buffer
//     err := executor.Execute(context.Background(), "fmt.Println(\"Hello\")", &output)
//     if err != nil {
//         log.Fatalf("Error: %v", err)
//     }
//     fmt.Println(output.String())
type CodeExecutor interface {
	Execute(ctx context.Context, code string, output io.Writer) error
}

// GoExecutor implementa la ejecución de código Go mediante el comando 'go run'.
//
// Esta implementación crea un archivo temporal con el código proporcionado,
// ejecuta 'go run' sobre ese archivo, y captura la salida estándar y de error.
// Incluye límites para la cantidad de salida generada y utiliza un pool de buffers
// para optimizar el uso de memoria.
type GoExecutor struct {
	goExecutablePath string
	maxOutputLength  int
	tempDir          string
	bufferPool       sync.Pool
}

// NewGoExecutor crea un nuevo ejecutor de código Go.
//
// Parámetros:
//   - goExecutablePath: Ruta al ejecutable de Go (ej. "/usr/local/go/bin/go").
//   - maxOutputLength: Tamaño máximo en bytes de la salida permitida.
//   - tempDir: Directorio temporal donde se crearán los archivos de código.
//
// Retorna un nuevo GoExecutor configurado con los parámetros especificados.
//
// Ejemplo:
//
//     executor := executor.NewGoExecutor("/usr/local/go/bin/go", 10000, os.TempDir())
//     var output bytes.Buffer
//     err := executor.Execute(context.Background(), "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}", &output)
func NewGoExecutor(goExecutablePath string, maxOutputLength int, tempDir string) *GoExecutor {
	return &GoExecutor{
		goExecutablePath: goExecutablePath,
		maxOutputLength:  maxOutputLength,
		tempDir:          tempDir,
		bufferPool: sync.Pool{
			New: func() interface{} {
				// Crear un buffer de 1KB por defecto
				buf := make([]byte, 1024)
				return &buf
			},
		},
	}
}

// Execute ejecuta el código Go y escribe la salida en el writer proporcionado.
//
// Este método crea un archivo temporal con el código proporcionado, ejecuta 'go run'
// sobre ese archivo, y escribe la salida en el writer proporcionado. Utiliza el contexto
// para controlar timeouts y cancelación. Limita la cantidad de salida generada según
// maxOutputLength y utiliza un pool de buffers para optimizar el uso de memoria.
//
// Parámetros:
//   - ctx: Contexto para control de cancelación y timeout.
//   - code: El código Go a ejecutar.
//   - output: Writer donde se escribirá la salida de la ejecución.
//
// Retorna error si hay algún problema durante la ejecución.
//
// Ejemplo:
//
//     var output bytes.Buffer
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     defer cancel()
//     err := executor.Execute(ctx, "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}", &output)
//     if err != nil {
//         log.Printf("Error: %v", err)
//     } else {
//         fmt.Println("Resultado:", output.String())
//     }
func (ge *GoExecutor) Execute(ctx context.Context, code string, output io.Writer) error {
	// Crear archivo temporal para el código
	tmpFile, err := os.CreateTemp(ge.tempDir, "code-*.go")
	if err != nil {
		return fmt.Errorf("error creando archivo temporal: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		// Intentar eliminar el archivo temporal
		for i := 0; i < 3; i++ {
			if err := os.Remove(tmpPath); err == nil || os.IsNotExist(err) {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	if _, err := tmpFile.WriteString(code); err != nil {
		return fmt.Errorf("error escribiendo código: %w", err)
	}
	tmpFile.Close()

	// Configurar y ejecutar el comando
	cmd := exec.CommandContext(ctx, ge.goExecutablePath, "run", tmpPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error obteniendo salida del comando: %w", err)
	}
	// Combinar stderr con stdout
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error iniciando el comando: %w", err)
	}

	totalBytes := 0
	
	// Obtener un buffer del pool
	bufPtr := ge.bufferPool.Get().(*[]byte)
	buf := *bufPtr
	
	// Asegurar que el buffer se devuelva al pool
	defer ge.bufferPool.Put(bufPtr)
	
	for {
		n, err := stdoutPipe.Read(buf)
		if n > 0 {
			// Limitar la cantidad total de bytes enviados
			if totalBytes+n > ge.maxOutputLength {
				allowed := ge.maxOutputLength - totalBytes
				if allowed > 0 {
					output.Write(buf[:allowed])
					totalBytes += allowed
				}
				fmt.Fprint(output, "\n... (output truncated)")
				break
			} else {
				output.Write(buf[:n])
				totalBytes += n
			}
		}
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error leyendo salida: %w", err)
			}
			break
		}
	}

	// Esperar a que el comando finalice
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error en la ejecución: %w", err)
	}
	
	return nil
}
