package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	testCases := []struct {
		level Level
		msg   string
		want  string
	}{
		{level: Info, msg: "This is an info message", want: "[INFO] This is an info message"},
		{level: Debug, msg: "This is an info message", want: "[INFO] This is an info message"},
		{level: Warn, msg: "This is an info message", want: ""},
		{level: Error, msg: "This is an info message", want: ""},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		logger := log.New(&buf, "", log.LstdFlags)
		l := &SimpleLogger{level: tc.level, logger: logger}

		l.Info(tc.msg)
		if !contains(buf.String(), tc.want) {
			t.Errorf("expected info message, got %s", buf.String())
		}
	}

	testCases = []struct {
		level Level
		msg   string
		want  string
	}{
		{level: Info, msg: "This is an error message", want: "[ERROR] This is an error message"},
		{level: Debug, msg: "This is an error message", want: "[ERROR] This is an error message"},
		{level: Warn, msg: "This is an error message", want: "[ERROR] This is an error message"},
		{level: Error, msg: "This is an error message", want: "[ERROR] This is an error message"},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		logger := log.New(&buf, "", log.LstdFlags)
		l := &SimpleLogger{level: tc.level, logger: logger}

		l.Error(tc.msg)
		if !contains(buf.String(), tc.want) {
			t.Errorf("expected error message, got %s", buf.String())
		}
	}

	testCases = []struct {
		level Level
		msg   string
		want  string
	}{
		{level: Info, msg: "This is a debug message", want: ""},
		{level: Debug, msg: "This is a debug message", want: "[DEBUG] This is a debug message"},
		{level: Warn, msg: "This is a debug message", want: ""},
		{level: Error, msg: "This is a debug message", want: ""},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		logger := log.New(&buf, "", log.LstdFlags)
		l := &SimpleLogger{level: tc.level, logger: logger}

		l.Debug(tc.msg)
		if !contains(buf.String(), tc.want) {
			t.Errorf("expected debug message, got %s", buf.String())
		}
	}

	testCases = []struct {
		level Level
		msg   string
		want  string
	}{
		{level: Info, msg: "This is a warn message", want: ""},
		{level: Debug, msg: "This is a warn message", want: "[WARN] This is a warn message"},
		{level: Warn, msg: "This is a warn message", want: "[WARN] This is a warn message"},
		{level: Error, msg: "This is a warn message", want: ""},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		logger := log.New(&buf, "", log.LstdFlags)
		l := &SimpleLogger{level: tc.level, logger: logger}

		l.Warn(tc.msg)
		if !contains(buf.String(), tc.want) {
			t.Errorf("expected warn message, got %s", buf.String())
		}
	}
}

func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}
