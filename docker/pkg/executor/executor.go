package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// CodeExecutor define el comportamiento para ejecutar código Go
type CodeExecutor interface {
	Execute(ctx context.Context, code string, output io.Writer) error
}

// GoExecutor implementa la ejecución de código Go
type GoExecutor struct {
	goExecutablePath string
	maxOutputLength  int
	tempDir          string
}

// NewGoExecutor crea un nuevo ejecutor de código Go
func NewGoExecutor(goExecutablePath string, maxOutputLength int, tempDir string) *GoExecutor {
	return &GoExecutor{
		goExecutablePath: goExecutablePath,
		maxOutputLength:  maxOutputLength,
		tempDir:          tempDir,
	}
}

// Execute ejecuta el código Go y escribe la salida en el writer proporcionado
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
	buf := make([]byte, 1024)
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
