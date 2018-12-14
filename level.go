package xlog

import (
	"fmt"
	"strings"
)

// Log Level
const (
	PrintLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

// Level type
type Level uint8

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case PrintLevel:
		return "print"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

// LogStr 返回适用于日志的字符串
func (level Level) LogStr() string {
	switch level {
	case PrintLevel:
		return ""
	case DebugLevel:
		return "[DEBUG]"
	case InfoLevel:
		return "[INFO ]"
	case WarnLevel:
		return "[WARN ]"
	case ErrorLevel:
		return "[ERROR]"
	case FatalLevel:
		return "[FATAL]"
	case PanicLevel:
		return "[PANIC]"
	}
	return "[UNKNO]"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "print":
		return PrintLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %q", lvl)
}
