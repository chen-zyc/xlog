package xlog

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestReqLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewWithWriter(buf, &Config{Level: LevelError, Prefix: "TEST: "})
	rl := NewReqLogger(l, ReqConfig{ReqID: "reqid", Level: LevelDebug})
	rl.Debugf("debug message, err: %v", errors.New("test error"))
	rl.Infoln("abc", "def")
	rl.Error("error message")

	expected := "TEST: [DEBU][reqid]debug message, err: test error\n"
	expected += "TEST: [INFO][reqid]abc def\n"
	expected += "TEST: [ERRO][reqid]error message\n"
	if buf.String() != expected {
		t.Fatalf("wrong output: %s", strings.Replace(buf.String(), "\n", "|", -1))
	}

	buf.Reset()

	l.Debugf("debug message, err: %v", errors.New("test error"))
	l.Infoln("abc", "def")
	l.Error("error message")
	expected = "TEST: [ERRO]error message\n"
	if buf.String() != expected {
		t.Fatalf("wrong output: %s", strings.Replace(buf.String(), "\n", "|", -1))
	}
}

func ExampleReqLogger() {
	b := new(bytes.Buffer)
	l := NewWithWriter(b, nil)
	rl := NewReqLogger(l, ReqConfig{ReqID: "testReqID", Level: LevelDebug})
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
	rl := NewReqLogger(NewWithWriter(buf, nil), ReqConfig{Level: LevelDebug})
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
