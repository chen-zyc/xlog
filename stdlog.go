package xlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// logger 是 Logger 接口的默认实现。
type logger struct {
	mu            sync.Mutex // for out.
	out           io.Writer
	prefix        string // prefix 写在每行前面
	flag          int    // 比如 LstdFlags
	baseCallDepth int    // 基础的深度，实际调用 Caller 的深度等于该值加上 Output 的 calldepth 值。
	level         Level
	buf           []byte
}

// 默认 Logger.
var defaultLogger = New(WriterOpt(os.Stderr), FlagOpt(LstdFlags), BaseCallDepthOpt(1))

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * level
//   * file and line number (if corresponding flags are provided).
func (l *logger) formatHeader(lvl Level, buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, l.prefix...)
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	// log level
	*buf = append(*buf, lvl.LogStr()...)

	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (l *logger) Output(lvl Level, calldepth int, s string) error {
	now := time.Now() // get this early.

	// 获取打印日志的文件位置。
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		calldepth += l.baseCallDepth
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}

	l.buf = l.buf[:0]
	l.formatHeader(lvl, &l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	// 如果没有添加换行符则在最后添加换行符。
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *logger) Print(v ...interface{}) { l.Output(PrintLevel, 2, fmt.Sprint(v...)) }

func (l *logger) Printf(format string, v ...interface{}) {
	l.Output(PrintLevel, 2, fmt.Sprintf(format, v...))
}

func (l *logger) Println(v ...interface{}) { l.Output(PrintLevel, 2, fmt.Sprintln(v...)) }

func (l *logger) Debug(v ...interface{}) {
	if l.level <= DebugLevel {
		l.Output(DebugLevel, 2, fmt.Sprint(v...))
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		l.Output(DebugLevel, 2, fmt.Sprintf(format, v...))
	}
}

func (l *logger) Debugln(v ...interface{}) {
	if l.level <= DebugLevel {
		l.Output(DebugLevel, 2, fmt.Sprintln(v...))
	}
}

func (l *logger) Info(v ...interface{}) {
	if l.level <= InfoLevel {
		l.Output(InfoLevel, 2, fmt.Sprint(v...))
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.Output(InfoLevel, 2, fmt.Sprintf(format, v...))
	}
}

func (l *logger) Infoln(v ...interface{}) {
	if l.level <= InfoLevel {
		l.Output(InfoLevel, 2, fmt.Sprintln(v...))
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.level <= WarnLevel {
		l.Output(WarnLevel, 2, fmt.Sprint(v...))
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		l.Output(WarnLevel, 2, fmt.Sprintf(format, v...))
	}
}

func (l *logger) Warnln(v ...interface{}) {
	if l.level <= WarnLevel {
		l.Output(WarnLevel, 2, fmt.Sprintln(v...))
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.level <= ErrorLevel {
		l.Output(ErrorLevel, 2, fmt.Sprint(v...))
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		l.Output(ErrorLevel, 2, fmt.Sprintf(format, v...))
	}
}

func (l *logger) Errorln(v ...interface{}) {
	if l.level <= ErrorLevel {
		l.Output(ErrorLevel, 2, fmt.Sprintln(v...))
	}
}

var osExit = os.Exit // for testing

func (l *logger) Fatal(v ...interface{}) {
	if l.level <= FatalLevel {
		l.Output(FatalLevel, 2, fmt.Sprint(v...))
		osExit(1)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	if l.level <= FatalLevel {
		l.Output(FatalLevel, 2, fmt.Sprintf(format, v...))
		osExit(1)
	}
}

func (l *logger) Fatalln(v ...interface{}) {
	if l.level <= FatalLevel {
		l.Output(FatalLevel, 2, fmt.Sprintln(v...))
		osExit(1)
	}
}

func (l *logger) Panic(v ...interface{}) {
	if l.level <= PanicLevel {
		s := fmt.Sprint(v...)
		l.Output(PanicLevel, 2, s)
		panic(s)
	}
}

func (l *logger) Panicf(format string, v ...interface{}) {
	if l.level <= PanicLevel {
		s := fmt.Sprintf(format, v...)
		l.Output(PanicLevel, 2, s)
		panic(s)
	}
}

func (l *logger) Panicln(v ...interface{}) {
	if l.level <= PanicLevel {
		s := fmt.Sprintln(v...)
		l.Output(PanicLevel, 2, s)
		panic(s)
	}
}

func (l *logger) SetOptions(opts ...Option) {
	for _, opt := range opts {
		if opt != nil {
			opt(l)
		}
	}
}

// SetOptions 设置默认实现的选项。
func SetOptions(opts ...Option) { defaultLogger.SetOptions(opts...) }

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
func Output(lvl Level, calldepth int, s string) error { return defaultLogger.Output(lvl, calldepth, s) }

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
