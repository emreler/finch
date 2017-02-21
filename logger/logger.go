package logger

import (
	"encoding/json"
	"log"

	"github.com/bsphere/le_go"
	"github.com/emreler/finch/config"
)

const (
	levelInfo  = "INFO"
	levelError = "ERROR"
)

// Logger .
type Logger struct {
	conn *le_go.Logger
}

// LogMessage .
type LogMessage struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// NewLogger .
func NewLogger(token config.LogentriesConfig) *Logger {
	le, err := le_go.Connect(string(token))

	if err != nil {
		panic(err)
	}

	log.Println("Connected to Logentries")

	return &Logger{le}
}

// Info .
func (l *Logger) Info(data interface{}) {
	var j []byte
	if str, ok := data.(string); ok {
		logMsg := &LogMessage{
			Level:   levelInfo,
			Message: str,
		}

		j, _ = json.Marshal(logMsg)

		log.Println(string(j))
		l.conn.Println(string(j))
	} else {
		jstring, _ := json.Marshal(data)

		logMsg := &LogMessage{
			Level:   levelInfo,
			Message: string(jstring),
		}

		j, _ = json.Marshal(logMsg)

		log.Println(string(j))
		l.conn.Println(string(j))
	}
}

func (l *Logger) Error(err error) {
	logMsg := &LogMessage{
		Level:   levelError,
		Message: err.Error(),
	}

	j, _ := json.Marshal(logMsg)

	log.Println(string(j))
	l.conn.Println(string(j))
}
