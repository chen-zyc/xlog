package xlog

import "io"

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

// Logger 日志接口。
type Logger interface {
	SetOptions(opts ...Option)

	// Output writes the output for a logging event. The string s contains
	// the text to print after the prefix specified by the flags of the
	// Logger. A newline is appended if the last character of s is not
	// already a newline. Calldepth is used to recover the PC and is
	// provided for generality, although at the moment on all pre-defined
	// paths it will be 2.
	Output(calldepth int, s string) error

	// Print calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Print(v ...interface{})

	// Printf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(format string, v ...interface{})

	// Println calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Println(v ...interface{})

	// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
	Fatal(v ...interface{})

	// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
	Fatalf(format string, v ...interface{})

	// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
	Fatalln(v ...interface{})

	// Panic is equivalent to l.Print() followed by a call to panic().
	Panic(v ...interface{})

	// Panicf is equivalent to l.Printf() followed by a call to panic().
	Panicf(format string, v ...interface{})

	// Panicln is equivalent to l.Println() followed by a call to panic().
	Panicln(v ...interface{})
}

// Option 在创建 Logger 时可以指定参数
type Option func(l *logger)

// WriterOpt 设置输出。
func WriterOpt(w io.Writer) Option {
	return func(l *logger) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.out = w
	}
}

// PrefixOpt 设置前缀。
func PrefixOpt(p string) Option {
	return func(l *logger) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.prefix = p
	}
}

// FlagOpt 设置输出格式。
func FlagOpt(flag int) Option {
	return func(l *logger) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.flag = flag
	}
}

// BaseCallDepthOpt 设置基础调用深度。
func BaseCallDepthOpt(depth int) Option {
	return func(l *logger) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.baseCallDepth = depth
	}
}

// New 根据 opts 指定的参数返回 Logger 的一个实现。
func New(opts ...Option) Logger {
	l := &logger{}
	l.SetOptions(opts...)
	return l
}
