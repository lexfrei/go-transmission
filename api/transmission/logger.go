package transmission

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value any
}

// Logger is the interface for structured logging.
// Implementations can wrap slog, zap, logrus, or any other logging library.
type Logger interface {
	// Debug logs a message at debug level with optional fields.
	Debug(msg string, fields ...Field)
	// Warn logs a message at warn level with optional fields.
	Warn(msg string, fields ...Field)
	// Error logs a message at error level with optional fields.
	Error(msg string, fields ...Field)
}

// noopLogger is a logger that does nothing.
type noopLogger struct{}

// NoopLogger returns a logger that discards all log messages.
// This is used as the default when no logger is provided.
func NoopLogger() Logger {
	return &noopLogger{}
}

func (n *noopLogger) Debug(_ string, _ ...Field) {}
func (n *noopLogger) Warn(_ string, _ ...Field)  {}
func (n *noopLogger) Error(_ string, _ ...Field) {}
