package logger

import (
	"log"
	"os"
)

type Level string

const (
	Info  Level = "info"
	Error Level = "error"
	Debug Level = "debug"
	Warn  Level = "warn"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type SimpleLogger struct {
	level  Level
	logger *log.Logger
}

func New(level Level) *SimpleLogger {
	return &SimpleLogger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *SimpleLogger) Info(msg string) {
	if l.level == Info || l.level == Debug {
		l.logger.Println("[INFO] " + msg)
	}
}

func (l *SimpleLogger) Error(msg string) {
	if l.level != "" {
		l.logger.Println("[ERROR] " + msg)
	}
}

func (l *SimpleLogger) Debug(msg string) {
	if l.level == Debug {
		l.logger.Println("[DEBUG] " + msg)
	}
}

func (l *SimpleLogger) Warn(msg string) {
	if l.level == Warn || l.level == Debug {
		l.logger.Println("[WARN] " + msg)
	}
}
