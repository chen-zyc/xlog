package xlog

import (
	"bytes"
	"fmt"
	"testing"
)

func ExampleReqLogger() {
	b := new(bytes.Buffer)
	l := New(WriterOpt(b))
	rl := NewReqLogger(LoggerOpt(l), ReqIDOpt("testReqID"), ReqLevelOpt(DebugLevel))
	rl.Debugln("debug message")
	rl.Infoln("info message")
	rl.Errorln("error message")
	fmt.Println(b.String())

	// Output:
	// [DEBU][testReqID]debug message
	// [INFO][testReqID]info message
	// [ERRO][testReqID]error message
}

func TestReqIDGen(t *testing.T) {
	reqID := ReqIDGen()
	if len(reqID) == 0 {
		t.Fatalf("invalid reqid: %s", reqID)
	}
}

func BenchmarkReqLogger(b *testing.B) {
	buf := new(bytes.Buffer)
	l := New(WriterOpt(buf))
	rl := NewReqLogger(LoggerOpt(l), ReqLevelOpt(DebugLevel))
	for i := 0; i < b.N; i++ {
		buf.Reset()
		rl.Debugln("hello")
	}
}

func BenchmarkReqIDGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ReqIDGen()
	}
}
