package xlog

import (
	"io"
	"os"
)

// 输出格式
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// Logger can print log to writer.
// You can new a Logger by calling `New()` or `NewWithWriter()`.
type Logger interface {
	// Output is the lowest level function.
	Output(lvl Level, calldepth int, reqID, s string) error

	// return the copy of Config.
	CopyConfig() Config

	// Print is not affected by `Log Level`。 It will print the message regardless of `Log Level`.
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

// New creates a new Logger with the specified config `c`.
func New(c *Config) Logger {
	return NewWithWriter(nil, c)
}

// NewWithWriter creates a Logger with the specified writer `w`.
func NewWithWriter(w io.Writer, c *Config) Logger {
	if w == nil {
		w = os.Stderr
	}
	return newLogger(w, c)
}
