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
	buf           []byte
}

// 默认 Logger.
var defaultLogger = New(WriterOpt(os.Stderr), FlagOpt(LstdFlags), BaseCallDepthOpt(1))

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (l *logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
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

func (l *logger) Output(calldepth int, s string) error {
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
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	// 如果没有添加换行符则在最后添加换行符。
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *logger) Print(v ...interface{}) { l.Output(2, fmt.Sprint(v...)) }

func (l *logger) Printf(format string, v ...interface{}) { l.Output(2, fmt.Sprintf(format, v...)) }

func (l *logger) Println(v ...interface{}) { l.Output(2, fmt.Sprintln(v...)) }

func (l *logger) Fatal(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *logger) Fatalln(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func (l *logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(2, s)
	panic(s)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

func (l *logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
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

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) { defaultLogger.Print(v...) }

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) { defaultLogger.Printf(format, v...) }

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) { defaultLogger.Println(v...) }

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
func Output(calldepth int, s string) error { return defaultLogger.Output(calldepth, s) }

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
