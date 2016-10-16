package logger

import (
	"encoding/json"
	"log"

	"github.com/bsphere/le_go"
	"gitlab.com/emreler/finch/config"
)

// Logger .
type Logger struct {
	conn *le_go.Logger
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
	if str, ok := data.(string); ok {
		l.conn.Println(str)
	} else {
		jstring, _ := json.Marshal(data)
		l.conn.Println(string(jstring))
	}
}

func (l *Logger) Error(err error) {
	l.conn.Println(err.Error())
}
