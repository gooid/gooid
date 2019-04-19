// utility
package util

import "fmt"
import "testing"

func doAssert(con interface{}, v ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(" recover:\n", err)
		}
	}()
	Assert(con, v...)
}

func testAssert(t *testing.T) {
	var err error
	doAssert(true)
	doAssert(nil)
	doAssert(err)
	doAssert(true, "bool", "(true)", 1)
	doAssert(nil, "nil", "(nil)", 2)
	doAssert(err, "error", "(err)", 3)

	doAssert(false)
	err = fmt.Errorf("Errof(%d)", 1)
	doAssert(err)

	doAssert(false, "bool", "(false)", 4)
	doAssert(err, "error", "(err)", 5)
}

func TestAssert(t *testing.T) {
	Mode(Print)
	Output(Fmt)
	Fmt.Print(">>>>>>>>>> Fmt Print\n")
	testAssert(t)
	Output(Log)
	Log.Print(">>>>>>>>>> Log Print\n")
	testAssert(t)

	Mode(Trace)
	Output(Fmt)
	Fmt.Print(">>>>>>>>>> Fmt Trace\n")
	testAssert(t)
	Output(Log)
	Log.Print(">>>>>>>>>> Log Trace\n")
	testAssert(t)

	Mode(StackTrace)
	Output(Fmt)
	Fmt.Print(">>>>>>>>>> Fmt StackTrace\n")
	testAssert(t)
	Output(Log)
	Log.Print(">>>>>>>>>> Log StackTrace\n")
	testAssert(t)
}

func TestAssertPanic(t *testing.T) {
	Mode(Panic)
	Output(Fmt)
	testAssert(t)
	Output(Log)
	testAssert(t)
}

func TestAssertObject(t *testing.T) {
	a := New(Print, Fmt)
	a.Assert(false, "TestAssertObject Print")
	a.Mode(Trace)
	a.Assert(false, "TestAssertObject Trace")
	a.Mode(StackTrace)
	a.Assert(false, "TestAssertObject StackTrace")
}

func _TestAssertF(t *testing.T) {
	Mode(Fatal)
	Output(Fmt)
	testAssert(t)
	Output(Log)
	testAssert(t)
}
