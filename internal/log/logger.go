package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Skryensya/footprint/internal/domain"
)

// Level representa el nivel de severidad del log
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string to a Level.
// Valid values: "debug", "info", "warn", "error" (case insensitive).
// Returns LevelWarn if the string is not recognized.
func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelWarn
	}
}

// Logger maneja el logging a archivo de forma thread-safe
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	minLevel Level
	enabled  bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init inicializa el logger global con el archivo especificado
func Init(logPath string, minLevel Level) error {
	var err error
	once.Do(func() {
		defaultLogger, err = New(logPath, minLevel)
	})
	return err
}

// New crea un nuevo logger que escribe al archivo especificado
func New(logPath string, minLevel Level) (*Logger, error) {
	// Crear directorio si no existe con permisos restrictivos
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0700); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	// Abrir archivo de log con permisos restrictivos
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	// Asegurar permisos correctos si el archivo ya existía
	if err := os.Chmod(logPath, 0600); err != nil {
		file.Close()
		return nil, fmt.Errorf("chmod log file: %w", err)
	}

	return &Logger{
		file:     file,
		minLevel: minLevel,
		enabled:  true,
	}, nil
}

// Close cierra el logger
func (l *Logger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}

// SetEnabled habilita o deshabilita el logging
func (l *Logger) SetEnabled(enabled bool) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = enabled
}

// log escribe un mensaje con el nivel especificado
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if l == nil || !l.enabled || level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, level.String(), message)

	l.file.Write([]byte(logLine))
}

// Debug escribe un mensaje de debug
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info escribe un mensaje informativo
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn escribe un warning
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error escribe un error
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Writer retorna un io.Writer que escribe al log con el nivel especificado
func (l *Logger) Writer(level Level) io.Writer {
	return &logWriter{logger: l, level: level}
}

type logWriter struct {
	logger *Logger
	level  Level
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	w.logger.log(w.level, "%s", string(p))
	return len(p), nil
}

// Funciones de conveniencia para el logger global

// Debug escribe un mensaje de debug al logger global
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

// Info escribe un mensaje informativo al logger global
func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

// Warn escribe un warning al logger global
func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, args...)
	}
}

// Error escribe un error al logger global
func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}

// Close cierra el logger global
func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}

// GetLogger retorna el logger global (puede ser nil si no se inicializó)
func GetLogger() *Logger {
	return defaultLogger
}

// NopLogger is a logger that discards all messages.
// Useful for testing or when logging is disabled.
type NopLogger struct{}

func (NopLogger) Debug(_ string, _ ...any) {}
func (NopLogger) Info(_ string, _ ...any)  {}
func (NopLogger) Warn(_ string, _ ...any)  {}
func (NopLogger) Error(_ string, _ ...any) {}
func (NopLogger) Close() error             { return nil }

// Verify Logger implements domain.Logger
var _ domain.Logger = (*Logger)(nil)
var _ domain.Logger = NopLogger{}
