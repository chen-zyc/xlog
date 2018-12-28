package xlog

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogLevel(t *testing.T) {
	osExit = func(code int) {}
	defer func() {
		osExit = os.Exit
	}()

	lvls := []Level{LevelPrint, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelPanic}
	info := []string{
		LevelPrint.LogStr() + "P",
		LevelDebug.LogStr() + "D",
		LevelInfo.LogStr() + "I",
		LevelWarn.LogStr() + "W",
		LevelError.LogStr() + "E",
		LevelFatal.LogStr() + "F",
		LevelPanic.LogStr() + "A",
	}

	logFunc := []func(l Logger){
		func(l Logger) {
			l.Print("P")
			l.Debug("D")
			l.Info("I")
			l.Warn("W")
			l.Error("E")
			l.Fatal("F")
			func() {
				defer func() {
					if err := recover(); err == nil {
						t.Fatal("should panic")
					}
				}()
				l.Panic("A")
			}()
		},
		func(l Logger) {
			l.Println("P")
			l.Debugln("D")
			l.Infoln("I")
			l.Warnln("W")
			l.Errorln("E")
			l.Fatalln("F")
			func() {
				defer func() {
					if err := recover(); err == nil {
						t.Fatal("should panic")
					}
				}()
				l.Panicln("A")
			}()
		},
		func(l Logger) {
			l.Printf("%s", "P")
			l.Debugf("%s", "D")
			l.Infof("%s", "I")
			l.Warnf("%s", "W")
			l.Errorf("%s", "E")
			l.Fatalf("%s", "F")
			func() {
				defer func() {
					if err := recover(); err == nil {
						t.Fatal("should panic")
					}
				}()
				l.Panicf("%s", "A")
			}()
		},
	}

	for i, lvl := range lvls {
		b := new(bytes.Buffer)
		l := NewWithWriter(b, &Config{Level: lvl})

		for j, f := range logFunc {
			b.Reset()
			f(l)

			want := strings.Join(info[i:], "\n") + "\n"
			if i != 0 { // Print 不受 level 限制
				want = "P\n" + want
			}
			if want != b.String() {
				t.Fatalf("level = %s, func[%d], got %s, want %s",
					lvl.String(), j,
					strings.Replace(b.String(), "\n", " ", -1),
					strings.Replace(want, "\n", " ", -1))
			}
		}
	}
}

func TestParseLevel(t *testing.T) {
	cases := []struct {
		s      string
		lvl    Level
		hasErr bool
	}{
		{s: "print", lvl: LevelPrint},
		{s: "debug", lvl: LevelDebug},
		{s: "info", lvl: LevelInfo},
		{s: "warn", lvl: LevelWarn},
		{s: "warning", lvl: LevelWarn},
		{s: "error", lvl: LevelError},
		{s: "fatal", lvl: LevelFatal},
		{s: "panic", lvl: LevelPanic},
		{s: "", hasErr: true},
		{s: "invalid", hasErr: true},
	}

	for _, c := range cases {
		lvl, err := ParseLevel(c.s)
		if c.hasErr {
			if err == nil {
				t.Fatalf("case: %#v", c)
			}
			continue
		}
		if err != nil {
			t.Fatalf("don't want err, but got: %v", err)
		}
		if lvl != c.lvl {
			t.Fatalf("unexpected level: %v, want: %v", lvl, c.lvl)
		}
	}
}
