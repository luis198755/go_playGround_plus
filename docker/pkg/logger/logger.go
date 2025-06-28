package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Logger es la interfaz para el logging estructurado
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}

// zapLogger implementa la interfaz Logger usando zap
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger crea una nueva instancia de Logger
func NewLogger(development bool) Logger {
	once.Do(func() {
		var config zap.Config
		if development {
			// Configuración para desarrollo: más verbosa, salida legible por humanos
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			// Configuración para producción: JSON estructurado
			config = zap.NewProductionConfig()
		}
		
		var err error
		log, err = config.Build()
		if err != nil {
			// Si hay un error al construir el logger, fallback a un logger básico
			log = zap.New(zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			))
		}
	})
	
	return &zapLogger{
		logger: log,
	}
}

// Info registra un mensaje a nivel INFO
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Error registra un mensaje a nivel ERROR
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Debug registra un mensaje a nivel DEBUG
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Warn registra un mensaje a nivel WARN
func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Fatal registra un mensaje a nivel FATAL y termina la aplicación
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// With crea un nuevo logger con campos adicionales
func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
	}
}

// Field crea un campo para el logger
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
