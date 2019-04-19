// utility
package util

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

const iCallerSkip = 3

type iMode int

const (
	Fatal = iMode(iota)
	Panic
	Trace
	StackTrace
	Print
)

type iType int

type AssertPrint interface {
	Print(v ...interface{})
}

type asserter struct {
	mode iMode
	out  AssertPrint
}

func New(mode iMode, out AssertPrint) *asserter {
	return &asserter{mode: mode, out: out}
}

func (a *asserter) Mode(mode iMode) (ret iMode) {
	std.Assert(mode >= Fatal && Fatal <= Print, "error mode = ", mode)
	ret, a.mode = a.mode, mode
	return
}

func (a *asserter) Output(out AssertPrint) {
	a.out = out
}

func (a *asserter) assert(con interface{}, v ...interface{}) {
	if con == nil {
		return
	}

	var s string
	switch obj := con.(type) {
	case bool:
		if !obj {
			if len(v) > 0 {
				s = assertFmt(a.mode, v...)
			} else {
				mode := a.mode
				if mode != StackTrace {
					mode = Trace
				}
				s = assertFmt(a.mode, "Assert(false)")
			}
		}

	case error:
		s = assertFmt(a.mode, append(v, con)...)

	default:
		s = assertFmt(Print, "error <", con, "> must is bool/error.")
	}
	if s != "" {
		switch a.mode {
		case Fatal:
			a.out.Print(s)
			os.Exit(1)
		case Panic:
			a.out.Print(s)
			panic(s)
		case Print, Trace, StackTrace:
			a.out.Print(s)
		}
	}
}

func (a *asserter) Assert(con interface{}, v ...interface{}) {
	a.assert(con, v...)
}

func assertFmt(mode iMode, a ...interface{}) string {
	s := fmt.Sprintln(a...)
	if mode == StackTrace {
		stack := string(debug.Stack())

		if i := strings.Index(stack, strThisFile); i >= 0 {
			i += strings.Index(stack[i:], "\n")
			stack = stack[i+1:]
		}

		for {
			stacks := strings.SplitN(stack, "\n", 3)
			if strings.HasPrefix(strings.TrimSpace(stacks[1]), strThisFile) {
				stack = stacks[2]
			} else {
				stack = strings.Join(stacks, "\n")
				break
			}
		}
		if stack != "" {
			s += "--- caller stack ---\n" + stack
		}
	} else if mode == Trace {
		_, file, line, ok := runtime.Caller(iCallerSkip)
		if ok {
			s += fmt.Sprint("\t@Caller : ", file, ":", line, "\n")
		}
	}
	return s
}

type fmtOutput struct{}

func (*fmtOutput) Print(v ...interface{}) {
	fmt.Print(v...)
}

var (
	Log = log.New(os.Stderr, "", log.LstdFlags)
	Fmt = &fmtOutput{}
	std = New(Fatal, Log)
)

// Set assert out mode
func Mode(mode iMode) (ret iMode) {
	return std.Mode(mode)
}

func Output(out AssertPrint) {
	std.Output(out)
}

func Assert(con interface{}, v ...interface{}) {
	std.assert(con, v...)
}

var strThisFile = "assert.go"

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		strThisFile = file
	}
}
