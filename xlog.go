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

// Logger 日志接口。
type Logger interface {
	SetOptions(opts ...Option)
	Level() Level

	Output(lvl Level, calldepth int, reqID, s string) error

	// Print 不受 Level 的影响，不管 Level 是什么总是能够打印出日志
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

// Option 在创建 Logger 时可以指定参数
type Option func(v interface{})

// WriterOpt 设置输出。
func WriterOpt(w io.Writer) Option {
	return func(v interface{}) {
		if l, ok := v.(*logger); ok {
			l.rwm.Lock()
			defer l.rwm.Unlock()
			l.out = w
		}
	}
}

// PrefixOpt 设置前缀。
func PrefixOpt(p string) Option {
	return func(v interface{}) {
		if l, ok := v.(*logger); ok {
			l.rwm.Lock()
			defer l.rwm.Unlock()
			l.prefix = p
		}
	}
}

// FlagOpt 设置输出格式。
func FlagOpt(flag int) Option {
	return func(v interface{}) {
		if l, ok := v.(*logger); ok {
			l.rwm.Lock()
			defer l.rwm.Unlock()
			l.flag = flag
		}
	}
}

// BaseCallDepthOpt 设置基础调用深度。
func BaseCallDepthOpt(depth int) Option {
	return func(v interface{}) {
		if l, ok := v.(*logger); ok {
			l.rwm.Lock()
			defer l.rwm.Unlock()
			l.baseCallDepth = depth
		}
	}
}

// LevelOpt 设置等级。
func LevelOpt(lvl Level) Option {
	return func(v interface{}) {
		if l, ok := v.(*logger); ok {
			l.rwm.Lock()
			defer l.rwm.Unlock()
			l.level = lvl
		}
	}
}

// ForceColorsOpt 设置是否开启颜色输出。
func ForceColorsOpt(force bool) Option {
	return func(v interface{}) {
		if l, ok := v.(*logger); ok {
			l.rwm.Lock()
			defer l.rwm.Unlock()
			l.forceColors = force
		}
	}
}

// LoggerOpt 设置 ReqLogger 的 Logger 字段。仅对 ReqLogger 有效。
func LoggerOpt(l Logger) Option {
	return func(v interface{}) {
		if rl, ok := v.(*reqLogger); ok {
			rl.Logger = l
		}
	}
}

// ReqIDOpt 设置 ReqLogger 的 request id。仅对 ReqLogger 有效。
func ReqIDOpt(id string) Option {
	return func(v interface{}) {
		if rl, ok := v.(*reqLogger); ok {
			rl.reqID = id
		}
	}
}

// ReqLevelOpt 设置 ReqLogger 的 level, 仅对 ReqLogger 有效。
func ReqLevelOpt(lvl Level) Option {
	return func(v interface{}) {
		if rl, ok := v.(*reqLogger); ok {
			rl.level = lvl
		}
	}
}

// New 根据 opts 指定的参数返回 Logger 的一个实现。
func New(opts ...Option) Logger {
	l := &logger{}
	l.SetOptions(opts...)
	if l.out == nil {
		l.out = os.Stderr
	}
	return l
}
