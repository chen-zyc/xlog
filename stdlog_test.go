package xlog

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	Rdate         = `[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]`
	Rtime         = `[0-9][0-9]:[0-9][0-9]:[0-9][0-9]`
	Rmicroseconds = `\.[0-9][0-9][0-9][0-9][0-9][0-9]`
	Rline         = `(90|92):` // must update if the calls to l.Printf / l.Print below move
	Rlongfile     = `.*/[A-Za-z0-9_\-]+\.go:` + Rline
	Rshortfile    = `[A-Za-z0-9_\-]+\.go:` + Rline
)

func TestCallDepth(t *testing.T) {
	b := new(bytes.Buffer)

	// 自定义 logger
	l := NewWithWriter(b, &Config{Flag: Lshortfile})
	l.Print("d1")
	if !strings.Contains(b.String(), "26: d1") {
		t.Fatal("calldepth not match: ", b.String())
	}

	// 默认 logger
	b.Reset()
	l.Print("d2")
	if !strings.Contains(b.String(), "33: d2") {
		t.Fatal("calldepth not match: ", b.String())
	}

	// 组合, 不需要设置 calldepth
	b.Reset()
	type A struct {
		Logger
	}
	a := A{NewWithWriter(b, &Config{Flag: Lshortfile})}
	a.Print("d3")
	if !strings.Contains(b.String(), "44: d3") {
		t.Fatal("calldepth not match: ", b.String())
	}

	// 匿名函数封装，设置 calldepth = 1
	b.Reset()
	func() {
		l := NewWithWriter(b, &Config{Flag: Lshortfile, BaseCalldepth: 1})
		l.Print("d4")
	}()
	if !strings.Contains(b.String(), "54: d4") {
		t.Fatal("calldepth not match: ", b.String())
	}
}

type tester struct {
	flag    int
	prefix  string
	pattern string // regexp that log output must match; we add ^ and expected_text$ always
}

var tests = []tester{
	// individual pieces:
	{0, "", ""},
	{0, "XXX", "XXX"},
	{Ldate, "", Rdate + " "},
	{Ltime, "", Rtime + " "},
	{Ltime | Lmicroseconds, "", Rtime + Rmicroseconds + " "},
	{Lmicroseconds, "", Rtime + Rmicroseconds + " "}, // microsec implies time
	{Llongfile, "", Rlongfile + " "},
	{Lshortfile, "", Rshortfile + " "},
	{Llongfile | Lshortfile, "", Rshortfile + " "}, // shortfile overrides longfile
	// everything at once:
	{Ldate | Ltime | Lmicroseconds | Llongfile, "XXX", "XXX" + Rdate + " " + Rtime + Rmicroseconds + " " + Rlongfile + " "},
	{Ldate | Ltime | Lmicroseconds | Lshortfile, "XXX", "XXX" + Rdate + " " + Rtime + Rmicroseconds + " " + Rshortfile + " "},
}

// Test using Println("hello", 23, "world") or using Printf("hello %d world", 23)
func testPrint(t *testing.T, flag int, prefix string, pattern string, useFormat bool) {
	buf := new(bytes.Buffer)
	conf := defaultLogger.CopyConfig()
	conf.Flag = flag
	conf.Prefix = prefix
	ResetDefaultLogger(NewWithWriter(buf, &conf))
	if useFormat {
		Printf("hello %d world", 23)
	} else {
		Println("hello", 23, "world")
	}
	line := buf.String()
	line = line[0 : len(line)-1]
	pattern = "^" + pattern + "hello 23 world$"
	matched, err4 := regexp.MatchString(pattern, line)
	if err4 != nil {
		t.Fatal("pattern did not compile:", err4)
	}
	if !matched {
		t.Errorf("log output should match %q is %q", pattern, line)
	}
}

func TestAll(t *testing.T) {
	for _, testcase := range tests {
		testPrint(t, testcase.flag, testcase.prefix, testcase.pattern, false)
		testPrint(t, testcase.flag, testcase.prefix, testcase.pattern, true)
	}
}

func TestOutput(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := NewWithWriter(&b, nil)
	l.Println(testString)
	if expect := testString + "\n"; b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestUTCFlag(t *testing.T) {
	var b bytes.Buffer
	l := NewWithWriter(&b, &Config{Prefix: "Test:", Flag: Ldate | Ltime | LUTC})
	// Verify a log message looks right in the right time zone. Quantize to the second only.
	now := time.Now().UTC()
	l.Print("hello")
	want := fmt.Sprintf("Test:%d/%.2d/%.2d %.2d:%.2d:%.2d hello\n",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	got := b.String()
	if got == want {
		return
	}
	// It's possible we crossed a second boundary between getting now and logging,
	// so add a second and try again. This should very nearly always work.
	now = now.Add(time.Second)
	want = fmt.Sprintf("Test:%d/%.2d/%.2d %.2d:%.2d:%.2d hello\n",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	if got == want {
		return
	}
	t.Errorf("got %q; want %q", got, want)
}

func TestEmptyPrintCreatesLine(t *testing.T) {
	var b bytes.Buffer
	l := NewWithWriter(&b, &Config{Prefix: "Header:", Flag: LstdFlags})
	l.Print()
	l.Println("non-empty")
	output := b.String()
	if n := strings.Count(output, "Header"); n != 2 {
		t.Errorf("expected 2 headers, got %d", n)
	}
	if n := strings.Count(output, "\n"); n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func BenchmarkItoa(b *testing.B) {
	dst := make([]byte, 0, 64)
	for i := 0; i < b.N; i++ {
		dst = dst[0:0]
		itoa(&dst, 2015, 4)   // year
		itoa(&dst, 1, 2)      // month
		itoa(&dst, 30, 2)     // day
		itoa(&dst, 12, 2)     // hour
		itoa(&dst, 56, 2)     // minute
		itoa(&dst, 0, 2)      // second
		itoa(&dst, 987654, 6) // microsecond
	}
}

func BenchmarkPrintln(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := NewWithWriter(&buf, &Config{Flag: LstdFlags})
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkPrintlnNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := NewWithWriter(&buf, nil)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}
