package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestSimpleLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &SimpleLogger{level: Info, logger: logger}

	l.Info("This is an info message")
	if !contains(buf.String(), "[INFO] This is an info message") {
		t.Errorf("expected info message, got %s", buf.String())
	}

	buf.Reset()
	l = &SimpleLogger{level: Debug, logger: logger}
	l.Info("This is an info message")
	if !contains(buf.String(), "[INFO] This is an info message") {
		t.Errorf("expected info message, got %s", buf.String())
	}

	buf.Reset()
	l = &SimpleLogger{level: Error, logger: logger}
	l.Info("This is an info message")
	if contains(buf.String(), "[INFO] This is an info message") {
		t.Errorf("did not expect info message, got %s", buf.String())
	}
}

func TestSimpleLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &SimpleLogger{level: Error, logger: logger}

	l.Error("This is an error message")
	if !contains(buf.String(), "[ERROR] This is an error message") {
		t.Errorf("expected error message, got %s", buf.String())
	}
}

func TestSimpleLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &SimpleLogger{level: Debug, logger: logger}

	l.Debug("This is a debug message")
	if !contains(buf.String(), "[DEBUG] This is a debug message") {
		t.Errorf("expected debug message, got %s", buf.String())
	}

	buf.Reset()
	l = &SimpleLogger{level: Info, logger: logger}
	l.Debug("This is a debug message")
	if contains(buf.String(), "[DEBUG] This is a debug message") {
		t.Errorf("did not expect debug message, got %s", buf.String())
	}
}

func TestSimpleLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	l := &SimpleLogger{level: Warn, logger: logger}

	l.Warn("This is a warn message")
	if !contains(buf.String(), "[WARN] This is a warn message") {
		t.Errorf("expected warn message, got %s", buf.String())
	}

	buf.Reset()
	l = &SimpleLogger{level: Debug, logger: logger}
	l.Warn("This is a warn message")
	if !contains(buf.String(), "[WARN] This is a warn message") {
		t.Errorf("expected warn message, got %s", buf.String())
	}

	buf.Reset()
	l = &SimpleLogger{level: Info, logger: logger}
	l.Warn("This is a warn message")
	if contains(buf.String(), "[WARN] This is a warn message") {
		t.Errorf("did not expect warn message, got %s", buf.String())
	}
}

// Helper function to check if a string contains a substring.
func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}
