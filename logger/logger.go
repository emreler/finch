package logger

import (
	"encoding/json"
	"io"
	"log"
)

const (
	levelInfo  = "INFO"
	levelError = "ERROR"
	prefix     = ""
)

// Logger .
type Logger struct {
	logger *log.Logger
}

// NewLogger returns a new Logger.
func NewLogger(logDest io.Writer) *Logger {
	return &Logger{logger: log.New(logDest, prefix, log.Ldate|log.Ltime)}
}

func (l *Logger) print(level string, msg string) {
	l.logger.Printf("level=%s %s", level, msg)
}

// Info logs messages with INFO level. Parameter must be either string or JSON serializable structs.
func (l *Logger) Info(data interface{}) {
	if str, ok := data.(string); ok {
		l.print(levelInfo, str)
	} else {
		jstring, _ := json.Marshal(data)
		l.print(levelInfo, string(jstring))
	}
}

// Error logs error with ERROR level.
func (l *Logger) Error(err error) {
	l.print(levelError, err.Error())
}
