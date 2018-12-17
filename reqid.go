package xlog

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

var pid = uint32(os.Getpid())

// ReqIDGen 生成 request id。
var ReqIDGen = func() string {
	var buf [12]byte
	binary.LittleEndian.PutUint32(buf[:4], pid)
	binary.LittleEndian.PutUint64(buf[4:], uint64(time.Now().UnixNano()))
	// TODO: memprofile
	return base64.URLEncoding.EncodeToString(buf[:])
}

// ReqLogger 请求级别的 Logger, 在日志中会打印出 request id.
type ReqLogger interface {
	Logger
	ReqID() string
}

type reqLogger struct {
	Logger
	reqID string
	level Level
}

// NewReqLogger 返回 ReqLogger。ReqLogger 用于请求级别的日志打印。
func NewReqLogger(opts ...Option) ReqLogger {
	rl := &reqLogger{}
	for _, opt := range opts {
		opt(rl)
	}
	if rl.reqID == "" && ReqIDGen != nil {
		rl.reqID = ReqIDGen()
	}
	if rl.Logger == nil {
		rl.Logger = defaultLogger
	}

	return rl
}

func (rl *reqLogger) ReqID() string {
	return rl.reqID
}

func (rl *reqLogger) Level() Level {
	return rl.level
}

func (rl *reqLogger) Print(v ...interface{}) {
	rl.Output(PrintLevel, 2, rl.reqID, fmt.Sprint(v...))
}

func (rl *reqLogger) Printf(format string, v ...interface{}) {
	rl.Output(PrintLevel, 2, rl.reqID, fmt.Sprintf(format, v...))
}

func (rl *reqLogger) Println(v ...interface{}) {
	rl.Output(PrintLevel, 2, rl.reqID, fmt.Sprintln(v...))
}

func (rl *reqLogger) Debug(v ...interface{}) {
	if rl.level <= DebugLevel {
		rl.Output(DebugLevel, 2, rl.reqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Debugf(format string, v ...interface{}) {
	if rl.level <= DebugLevel {
		rl.Output(DebugLevel, 2, rl.reqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Debugln(v ...interface{}) {
	if rl.level <= DebugLevel {
		rl.Output(DebugLevel, 2, rl.reqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Info(v ...interface{}) {
	if rl.level <= InfoLevel {
		rl.Output(InfoLevel, 2, rl.reqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Infof(format string, v ...interface{}) {
	if rl.level <= InfoLevel {
		rl.Output(InfoLevel, 2, rl.reqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Infoln(v ...interface{}) {
	if rl.level <= InfoLevel {
		rl.Output(InfoLevel, 2, rl.reqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Warn(v ...interface{}) {
	if rl.level <= WarnLevel {
		rl.Output(WarnLevel, 2, rl.reqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Warnf(format string, v ...interface{}) {
	if rl.level <= WarnLevel {
		rl.Output(WarnLevel, 2, rl.reqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Warnln(v ...interface{}) {
	if rl.level <= WarnLevel {
		rl.Output(WarnLevel, 2, rl.reqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Error(v ...interface{}) {
	if rl.level <= ErrorLevel {
		rl.Output(ErrorLevel, 2, rl.reqID, fmt.Sprint(v...))
	}
}

func (rl *reqLogger) Errorf(format string, v ...interface{}) {
	if rl.level <= ErrorLevel {
		rl.Output(ErrorLevel, 2, rl.reqID, fmt.Sprintf(format, v...))
	}
}

func (rl *reqLogger) Errorln(v ...interface{}) {
	if rl.level <= ErrorLevel {
		rl.Output(ErrorLevel, 2, rl.reqID, fmt.Sprintln(v...))
	}
}

func (rl *reqLogger) Fatal(v ...interface{}) {
	if rl.level <= FatalLevel {
		rl.Output(FatalLevel, 2, rl.reqID, fmt.Sprint(v...))
		osExit(1)
	}
}

func (rl *reqLogger) Fatalf(format string, v ...interface{}) {
	if rl.level <= FatalLevel {
		rl.Output(FatalLevel, 2, rl.reqID, fmt.Sprintf(format, v...))
		osExit(1)
	}
}

func (rl *reqLogger) Fatalln(v ...interface{}) {
	if rl.level <= FatalLevel {
		rl.Output(FatalLevel, 2, rl.reqID, fmt.Sprintln(v...))
		osExit(1)
	}
}

func (rl *reqLogger) Panic(v ...interface{}) {
	if rl.level <= PanicLevel {
		s := fmt.Sprint(v...)
		rl.Output(PanicLevel, 2, rl.reqID, s)
		panic(s)
	}
}

func (rl *reqLogger) Panicf(format string, v ...interface{}) {
	if rl.level <= PanicLevel {
		s := fmt.Sprintf(format, v...)
		rl.Output(PanicLevel, 2, rl.reqID, s)
		panic(s)
	}
}

func (rl *reqLogger) Panicln(v ...interface{}) {
	if rl.level <= PanicLevel {
		s := fmt.Sprintln(v...)
		rl.Output(PanicLevel, 2, rl.reqID, s)
		panic(s)
	}
}
