package xlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// define the color of each level.
var (
	ColorPrint = "\x1b[0m"  // none
	ColorDebug = "\x1b[37m" // gray
	ColorInfo  = "\x1b[34m" // blue
	ColorWarn  = "\x1b[33m" // yellow
	ColorError = "\x1b[31m" // red
	ColorFatal = "\x1b[31m" // red
	ColorPanic = "\x1b[31m" // red
)

// Config for Logger.
type Config struct {
	Prefix        string
	Flag          int
	BaseCalldepth int
	Level         Level
	ForceColors   bool
	InitBufSize   int
}

// logger is the default implementation of the Logger interface.
type logger struct {
	wmut    sync.Mutex
	w       io.Writer
	bufPool sync.Pool
	Config
}

func newLogger(w io.Writer, c *Config) *logger {
	if c == nil {
		c = &Config{}
	}
	if c.InitBufSize < 0 {
		c.InitBufSize = 0
	}
	bufPool := sync.Pool{New: func() interface{} {
		return make([]byte, 0, c.InitBufSize)
	}}
	return &logger{
		w:       w,
		bufPool: bufPool,
		Config:  *c,
	}
}

func (l *logger) Output(lvl Level, calldepth int, reqID, s string) error {
	now := time.Now() // get this early.

	var file string
	var line int
	if l.Flag&(Lshortfile|Llongfile) != 0 {
		calldepth += l.BaseCalldepth
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
	}

	buf := l.bufPool.Get().([]byte)
	buf = buf[:0]
	// add color to header.
	if l.ForceColors {
		buf = append(buf, lvl.Color()...)
	}
	l.formatHeader(lvl, &buf, now, file, line, reqID)
	// clear color.
	if l.ForceColors {
		buf = append(buf, "\x1b[0m"...)
	}
	buf = append(buf, s...)
	// If there is no newline at the end, add it.
	if len(s) == 0 || s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}
	l.wmut.Lock()
	_, err := l.w.Write(buf)
	l.wmut.Unlock()
	l.bufPool.Put(buf)
	return err
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * level
//   * reqID
//   * file and line number (if corresponding flags are provided).
func (l *logger) formatHeader(lvl Level, buf *[]byte, t time.Time, file string, line int, reqID string) {
	*buf = append(*buf, l.Prefix...)
	if l.Flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.Flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.Flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.Flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.Flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	// log level
	*buf = append(*buf, lvl.LogStr()...)

	// request id
	if reqID != "" {
		*buf = append(*buf, '[')
		*buf = append(*buf, reqID...)
		*buf = append(*buf, ']')
	}

	if l.Flag&(Lshortfile|Llongfile) != 0 {
		if l.Flag&Lshortfile != 0 {
			file = shortFile(file)
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func shortFile(f string) string {
	short := f
	for i := len(f) - 1; i > 0; i-- {
		if f[i] == '/' {
			short = f[i+1:]
			break
		}
	}
	return short
}

func (l *logger) CopyConfig() Config {
	return l.Config
}

func (l *logger) Print(v ...interface{}) {
	_ = l.Output(LevelPrint, 2, "", fmt.Sprint(v...))
}

func (l *logger) Printf(format string, v ...interface{}) {
	_ = l.Output(LevelPrint, 2, "", fmt.Sprintf(format, v...))
}

func (l *logger) Println(v ...interface{}) {
	_ = l.Output(LevelPrint, 2, "", fmt.Sprintln(v...))
}

func (l *logger) Debug(v ...interface{}) {
	if l.Level <= LevelDebug {
		_ = l.Output(LevelDebug, 2, "", fmt.Sprint(v...))
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.Level <= LevelDebug {
		_ = l.Output(LevelDebug, 2, "", fmt.Sprintf(format, v...))
	}
}

func (l *logger) Debugln(v ...interface{}) {
	if l.Level <= LevelDebug {
		_ = l.Output(LevelDebug, 2, "", fmt.Sprintln(v...))
	}
}

func (l *logger) Info(v ...interface{}) {
	if l.Level <= LevelInfo {
		_ = l.Output(LevelInfo, 2, "", fmt.Sprint(v...))
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.Level <= LevelInfo {
		_ = l.Output(LevelInfo, 2, "", fmt.Sprintf(format, v...))
	}
}

func (l *logger) Infoln(v ...interface{}) {
	if l.Level <= LevelInfo {
		_ = l.Output(LevelInfo, 2, "", fmt.Sprintln(v...))
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.Level <= LevelWarn {
		_ = l.Output(LevelWarn, 2, "", fmt.Sprint(v...))
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.Level <= LevelWarn {
		_ = l.Output(LevelWarn, 2, "", fmt.Sprintf(format, v...))
	}
}

func (l *logger) Warnln(v ...interface{}) {
	if l.Level <= LevelWarn {
		_ = l.Output(LevelWarn, 2, "", fmt.Sprintln(v...))
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.Level <= LevelError {
		_ = l.Output(LevelError, 2, "", fmt.Sprint(v...))
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if l.Level <= LevelError {
		_ = l.Output(LevelError, 2, "", fmt.Sprintf(format, v...))
	}
}

func (l *logger) Errorln(v ...interface{}) {
	if l.Level <= LevelError {
		_ = l.Output(LevelError, 2, "", fmt.Sprintln(v...))
	}
}

var osExit = os.Exit // for testing

func (l *logger) Fatal(v ...interface{}) {
	if l.Level <= LevelFatal {
		_ = l.Output(LevelFatal, 2, "", fmt.Sprint(v...))
		osExit(1)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	if l.Level <= LevelFatal {
		_ = l.Output(LevelFatal, 2, "", fmt.Sprintf(format, v...))
		osExit(1)
	}
}

func (l *logger) Fatalln(v ...interface{}) {
	if l.Level <= LevelFatal {
		_ = l.Output(LevelFatal, 2, "", fmt.Sprintln(v...))
		osExit(1)
	}
}

func (l *logger) Panic(v ...interface{}) {
	if l.Level <= LevelPanic {
		s := fmt.Sprint(v...)
		_ = l.Output(LevelPanic, 2, "", s)
		panic(s)
	}
}

func (l *logger) Panicf(format string, v ...interface{}) {
	if l.Level <= LevelPanic {
		s := fmt.Sprintf(format, v...)
		_ = l.Output(LevelPanic, 2, "", s)
		panic(s)
	}
}

func (l *logger) Panicln(v ...interface{}) {
	if l.Level <= LevelPanic {
		s := fmt.Sprintln(v...)
		_ = l.Output(LevelPanic, 2, "", s)
		panic(s)
	}
}

var defaultLogger = New(&Config{Level: LevelError, Flag: LstdFlags, BaseCalldepth: 1})

// ResetDefaultLogger replace defaultLogger with `l`.
// NOTE:
//  1. The BaseCalldepth need one more than l.
//  2. This function is not thread safe.
func ResetDefaultLogger(l Logger) {
	if l != nil {
		defaultLogger = l
	}
}

// Print calls Output to print to the default logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) { defaultLogger.Print(v...) }

// Printf calls Output to print to the default logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) { defaultLogger.Printf(format, v...) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) { defaultLogger.Println(v...) }

// Debug calls Output to print to the default logger.
func Debug(v ...interface{}) { defaultLogger.Debug(v...) }

// Debugf calls Output to print to the default logger.
func Debugf(format string, v ...interface{}) { defaultLogger.Debugf(format, v...) }

// Debugln calls Output to print to the default logger.
func Debugln(v ...interface{}) { defaultLogger.Debugln(v...) }

// Info calls Output to print to the default logger.
func Info(v ...interface{}) { defaultLogger.Info(v...) }

// Infof calls Output to print to the default logger.
func Infof(format string, v ...interface{}) { defaultLogger.Infof(format, v...) }

// Infoln calls Output to print to the default logger.
func Infoln(v ...interface{}) { defaultLogger.Infoln(v...) }

// Warn calls Output to print to the default logger.
func Warn(v ...interface{}) { defaultLogger.Warn(v...) }

// Warnf calls Output to print to the default logger.
func Warnf(format string, v ...interface{}) { defaultLogger.Warnf(format, v...) }

// Warnln calls Output to print to the default logger.
func Warnln(v ...interface{}) { defaultLogger.Warnln(v...) }

// Error calls Output to print to the default logger.
func Error(v ...interface{}) { defaultLogger.Error(v...) }

// Errorf calls Output to print to the default logger.
func Errorf(format string, v ...interface{}) { defaultLogger.Errorf(format, v...) }

// Errorln calls Output to print to the default logger.
func Errorln(v ...interface{}) { defaultLogger.Errorln(v...) }

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) { defaultLogger.Fatal(v...) }

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) { defaultLogger.Fatalf(format, v...) }

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) { defaultLogger.Fatalln(v...) }

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) { defaultLogger.Panic(v...) }

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) { defaultLogger.Panicf(format, v...) }

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) { defaultLogger.Panicln(v...) }

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if Llongfile or Lshortfile is set; a value of 1 will print the details
// for the caller of Output.
func Output(lvl Level, calldepth int, reqID, s string) error {
	return defaultLogger.Output(lvl, calldepth, reqID, s)
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
