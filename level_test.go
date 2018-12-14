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

	lvls := []Level{PrintLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	info := []string{
		PrintLevel.LogStr() + "P",
		DebugLevel.LogStr() + "D",
		InfoLevel.LogStr() + "I",
		WarnLevel.LogStr() + "W",
		ErrorLevel.LogStr() + "E",
		FatalLevel.LogStr() + "F",
		PanicLevel.LogStr() + "A",
	}
	b := new(bytes.Buffer)
	l := New(WriterOpt(b))

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
		l.SetOptions(LevelOpt(lvl))

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
