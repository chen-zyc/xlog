package xlog

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

var pid = uint32(os.Getpid())

// ReqIDGen generates request id。
var ReqIDGen = func() string {
	var buf [12]byte
	binary.LittleEndian.PutUint32(buf[:4], pid)
	binary.LittleEndian.PutUint64(buf[4:], uint64(time.Now().UnixNano()))
	// TODO: memprofile
	return base64.URLEncoding.EncodeToString(buf[:])
}

// ReqLogger is the request level Logger, and the request id will be printed in the log.
type ReqLogger interface {
	Logger
	RequestConfig() *ReqConfig
}

// ReqConfig is the config of ReqLogger.
type ReqConfig struct {
	ReqID string // if it is empty, call ReqIDGen to generate it.
	Level Level
}

type reqLogger struct {
	ReqConfig
	Logger
	calldepth int
}

// NewReqLogger creates a ReqLogger.
// If l is nil, the default logger is used.
func NewReqLogger(l Logger, c ReqConfig) ReqLogger {
	if c.ReqID == "" && ReqIDGen != nil {
		c.ReqID = ReqIDGen()
	}
	// Output() + Printx()
	calldepth := 2
	if l == nil {
		l = defaultLogger
		if c.Level == LevelPrint {
			c.Level = l.CopyConfig().Level
		}
		// 因为 defaultLogger 是通过全局函数调用的，会多加一层，但这里是直接调用 defaultLogger，所以需要减掉一层。
		calldepth--
	}

	return &reqLogger{
		ReqConfig: c,
		Logger:    l,
		calldepth: calldepth,
	}
}

func (rl *reqLogger) RequestConfig() *ReqConfig {
	return &rl.ReqConfig
}

func (rl *reqLogger) Print(v ...interface{}) {
	_ = rl.Output(LevelPrint, rl.calldepth, rl.ReqID, fmt.Sprint(v...))
}

func (rl *reqLogger) Printf(format string, v ...interface{}) {
	_ = rl.Output(LevelPrint, rl.calldepth, rl.ReqID, fmt.Sprintf(format, v...))
}

func (rl *reqLogger) Println(v ...interface{}) {
	_ = rl.Output(LevelPrint, rl.calldepth, rl.ReqID, fmt.Sprintln(v...))
}

func (rl *reqLogger) Debug(v ...interface{}) {
	if rl.Level <= LevelDebug {
		_ = rl.Output(LevelDebug, rl.calldepth, rl.ReqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Debugf(format string, v ...interface{}) {
	if rl.Level <= LevelDebug {
		_ = rl.Output(LevelDebug, rl.calldepth, rl.ReqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Debugln(v ...interface{}) {
	if rl.Level <= LevelDebug {
		_ = rl.Output(LevelDebug, rl.calldepth, rl.ReqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Info(v ...interface{}) {
	if rl.Level <= LevelInfo {
		_ = rl.Output(LevelInfo, rl.calldepth, rl.ReqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Infof(format string, v ...interface{}) {
	if rl.Level <= LevelInfo {
		_ = rl.Output(LevelInfo, rl.calldepth, rl.ReqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Infoln(v ...interface{}) {
	if rl.Level <= LevelInfo {
		_ = rl.Output(LevelInfo, rl.calldepth, rl.ReqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Warn(v ...interface{}) {
	if rl.Level <= LevelWarn {
		_ = rl.Output(LevelWarn, rl.calldepth, rl.ReqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Warnf(format string, v ...interface{}) {
	if rl.Level <= LevelWarn {
		_ = rl.Output(LevelWarn, rl.calldepth, rl.ReqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Warnln(v ...interface{}) {
	if rl.Level <= LevelWarn {
		_ = rl.Output(LevelWarn, rl.calldepth, rl.ReqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Error(v ...interface{}) {
	if rl.Level <= LevelError {
		_ = rl.Output(LevelError, rl.calldepth, rl.ReqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Errorf(format string, v ...interface{}) {
	if rl.Level <= LevelError {
		_ = rl.Output(LevelError, rl.calldepth, rl.ReqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Errorln(v ...interface{}) {
	if rl.Level <= LevelError {
		_ = rl.Output(LevelError, rl.calldepth, rl.ReqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Fatal(v ...interface{}) {
	if rl.Level <= LevelFatal {
		_ = rl.Output(LevelFatal, rl.calldepth, rl.ReqID, fmt.Sprint(v...))
		osExit(1)
	}
}

func (rl *reqLogger) Fatalf(format string, v ...interface{}) {
	if rl.Level <= LevelFatal {
		_ = rl.Output(LevelFatal, rl.calldepth, rl.ReqID, fmt.Sprintf(format, v...))
		osExit(1)
	}
}

func (rl *reqLogger) Fatalln(v ...interface{}) {
	if rl.Level <= LevelFatal {
		_ = rl.Output(LevelFatal, rl.calldepth, rl.ReqID, fmt.Sprintln(v...))
		osExit(1)
	}
}

func (rl *reqLogger) Panic(v ...interface{}) {
	if rl.Level <= LevelPanic {
		s := fmt.Sprint(v...)
		_ = rl.Output(LevelPanic, rl.calldepth, rl.ReqID, s)
		panic(s)
	}
}

func (rl *reqLogger) Panicf(format string, v ...interface{}) {
	if rl.Level <= LevelPanic {
		s := fmt.Sprintf(format, v...)
		_ = rl.Output(LevelPanic, rl.calldepth, rl.ReqID, s)
		panic(s)
	}
}

func (rl *reqLogger) Panicln(v ...interface{}) {
	if rl.Level <= LevelPanic {
		s := fmt.Sprintln(v...)
		_ = rl.Output(LevelPanic, rl.calldepth, rl.ReqID, s)
		panic(s)
	}
}
