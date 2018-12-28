package xlog

import (
	"fmt"
	"strings"
)

// Log Level
const (
	LevelPrint Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

// Level type
type Level uint8

// Convert the Level to a string.
func (level Level) String() string {
	switch level {
	case LevelPrint:
		return "print"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}

	return "unknown"
}

// LogStr returns the string that applices to the log.
func (level Level) LogStr() string {
	switch level {
	case LevelPrint:
		return ""
	case LevelDebug:
		return "[DEBU]"
	case LevelInfo:
		return "[INFO]"
	case LevelWarn:
		return "[WARN]"
	case LevelError:
		return "[ERRO]"
	case LevelFatal:
		return "[FATA]"
	case LevelPanic:
		return "[PANI]"
	}
	return "[UNKN]"
}

// Color returns the corresponding color.
func (level Level) Color() string {
	switch level {
	case LevelPrint:
		return ColorPrint
	case LevelDebug:
		return ColorDebug
	case LevelInfo:
		return ColorInfo
	case LevelWarn:
		return ColorWarn
	case LevelError:
		return ColorError
	case LevelFatal:
		return ColorFatal
	case LevelPanic:
		return ColorPanic
	}
	return ColorPrint
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return LevelPanic, nil
	case "fatal":
		return LevelFatal, nil
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	case "print":
		return LevelPrint, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %q", lvl)
}
