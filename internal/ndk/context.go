// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"os"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

// Callbacks is the set of functions called by the activity.
type Callbacks struct {
	Create               func(*Activity, []byte)
	Start                func(*Activity)
	Stop                 func(*Activity)
	Pause                func(*Activity)
	Resume               func(*Activity)
	Destroy              func(*Activity)
	SaveState            func(*Activity) []byte
	FocusChanged         func(*Activity, bool)
	ContentRectChanged   func(*Activity, *Rect)
	ConfigurationChanged func(*Activity)
	LowMemory            func(*Activity)

	// Window
	WindowCreated      func(*Activity, *Window)
	WindowDestroyed    func(*Activity, *Window)
	WindowDraw         func(*Activity, *Window)
	WindowResized      func(*Activity, *Window)
	WindowRedrawNeeded func(*Activity, *Window)

	// Touch is called by the app when a touch event occurs.
	Event func(*Activity, *InputEvent)
	// Sensor
	Sensor func(*Activity, []SensorEvent)
}

type Context struct {
	Callbacks

	act    *Activity
	window *Window // 在 onStart 后不一定调用 onNativeWindowCreated
	input  *InputQueue

	isResume,
	isFocus, // 在 onStart/onStop 之后不一定会调用 onWindowFocusChanged
	willDestory bool

	looper      *Looper
	read, write *os.File
	sensorQueue *SensorEventQueue
	funcChan    chan func()
	wait        *sync.WaitGroup
	isReady     bool
	isDebug     bool

	className  string
	savedState []byte
}

func (ctx *Context) init() {
	ctx.funcChan = make(chan func(), 3)
	ctx.wait = &sync.WaitGroup{}
}

func (ctx *Context) setCB(cb *Callbacks) {
	if cb != nil {
		ctx.Callbacks = *cb
	}
}

func (ctx *Context) reset() {
	ctx.isResume = false
	ctx.willDestory = false
}

func (ctx *Context) Debug(enable bool) {
	ctx.isDebug = enable
}

func (ctx *Context) WillDestory() bool {
	return ctx.willDestory
}

const (
	looper_ID_MAIN   = 1
	looper_ID_INPUT  = 2
	looper_ID_SENSOR = 3
)

// Init looper
func (ctx *Context) initLooper() {
	var err error
	ctx.read, ctx.write, err = os.Pipe()
	assert(err)
	ctx.looper = looperPrepare(LOOPER_PREPARE_ALLOW_NON_CALLBACKS)
	ctx.looper.AddFd(int(ctx.read.Fd()), looper_ID_MAIN, LOOPER_EVENT_INPUT, nil, unsafe.Pointer(uintptr(looper_ID_MAIN)))
}

func (ctx *Context) begin(cb Callbacks) {
	ctx.setCB(&cb)
	ctx.initLooper()
	ctx.funcChan <- func() {}
}

// Run starts the activity.
//
// It must be called directly from from the main function and will
// block until the app exits.
func (ctx *Context) Run(cb Callbacks) {
	info("ctx.Run")
	ctx.begin(cb)
	ctx.Loop()
}

func (ctx *Context) Begin(cb Callbacks) {
	info("ctx.Begin")
	ctx.begin(cb)
}

func (ctx *Context) Release() {
	ctx.looper.RemoveFd(int(ctx.read.Fd()))
	if ctx.sensorQueue != nil {
		SensorManagerInstance().destroy(ctx.sensorQueue)
	}
	//ctx.looper.Release()
	ctx.read.Close()
	ctx.write.Close()
}

func (ctx *Context) runFunc(fun func(), sync bool) {
	if !ctx.isDebug {
		defer func() {
			if r := recover(); r != nil {
				info("runFunc.recover:", r)
			}
		}()
	}

	//Info("...wait...")
	//oldM := util.Mode(util.StackTrace)
	//Assert(false)
	//util.Mode(oldM)

	var cmd = []byte{looper_ID_MAIN}
	ctx.write.Write(cmd)
	if !sync {
		ctx.funcChan <- fun
	} else {
		ctx.wait.Add(1)
		ctx.funcChan <- func() { defer ctx.wait.Done(); fun() }
		ctx.wait.Wait()
	}
	//Info("...wait<<<", forceEsace)
}

var entryDoFunc uintptr

func (ctx *Context) doFunc() {
	if !ctx.isDebug {
		defer func() {
			if r := recover(); r != nil {
				info("execFunc.recover:", r)
			}
		}()
	}

	if entryDoFunc == 0 {
		pc := make([]uintptr, 1)
		n := runtime.Callers(1, pc)
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()
		entryDoFunc = frame.Entry
	}

	(<-ctx.funcChan)()
}

//
func (ctx *Context) Call(fun func(), sync bool) {
	if sync {
		// 如果是在 doFunc 中被调用，则直接执行，否则会导致死锁。
		pc := make([]uintptr, 10000)
		n := runtime.Callers(2, pc)
		frames := runtime.CallersFrames(pc[:n])
		for {
			frame, more := frames.Next()
			if frame.Entry == entryDoFunc {
				fun()
				return
			}
			if !more {
				break
			}
		}
	}

	ctx.runFunc(fun, sync)
}

func (ctx *Context) loopEvents(timeoutMillis int) bool {
	assert(!ctx.willDestory, "!ctx.willDestory")

	for {
		ident, _, event, data := looperPollAll(timeoutMillis)
		//Info("LooperPollAll is", int(ident))
		switch ident {
		case looper_ID_MAIN:
			//Info("LooperPollAll is MAIN")
			var cmd = []byte{0}
			n, err := ctx.read.Read(cmd)
			assert(n == 1 && err == nil)
			assert(ident == int(cmd[0]))
			assert(ident == int(data))
			assert(event == int(data))
			ctx.doFunc()
			if ctx.willDestory {
				return false
			}
			return true

		case looper_ID_INPUT:
			//Info("LooperPollAll is INPUT")
			for ctx.input != nil {
				event := ctx.input.GetEvent()
				if event == nil {
					break
				}
				if !ctx.input.PreDispatchEvent(event) {
					if !ctx.willDestory {
						ctx.processEvent(event)
					}
					ctx.input.FinishEvent(event, 0)
				}
			}
			if ctx.willDestory {
				return false
			}
			return true

		case looper_ID_SENSOR:
			//Info("LooperPollAll is SENSOR")
			if ctx.sensorQueue != nil {
				var events []SensorEvent
				errno := ctx.sensorQueue.hasEvents()
				if errno > 0 {
					events, errno = ctx.sensorQueue.getEvents(32)
				}

				if errno < 0 {
					info("SensorEventQueue GetEvents error :", errno)
				} else {
					if ctx.Sensor != nil && len(events) > 0 {
						//for _, event := range events {
						//}
						ctx.Sensor(ctx.act, events)
					}
				}
			} else {
				info("looper_ID_SENSOR ...")
			}
			return true

		case LOOPER_POLL_WAKE:
			info("LooperPollAll is ALOOPER_POLL_WAKE")
			return true
		case LOOPER_POLL_TIMEOUT:
			//Info("LooperPollAll is LOOPER_POLL_TIMEOUT")
			return true
		case LOOPER_POLL_ERROR:
			info("LooperPollAll is LOOPER_POLL_ERROR")
			return true
		default:
			info("LooperPollAll is ", ident)
			return true
		}
	}
}

func (ctx *Context) processEvent(e *InputEvent) {
	//Info("processEvent:", e)
	if ctx.Event != nil {
		ctx.Event(ctx.act, e)
	}
}

func (ctx *Context) pollEvent(timeoutMillis int) bool {
	if ctx.loopEvents(timeoutMillis) {
		if ctx.isFocus && ctx.isResume {
			if ctx.window != nil &&
				ctx.WindowDraw != nil {
				ctx.WindowDraw(ctx.act, ctx.window)
			}
		}
		return true
	} else {
		//info("pollEvent: return false")
		ctx.doFunc()
		return false
	}
}

func (ctx *Context) Loop() {
	for ctx.pollEvent(-1) {
		// nothing
	}
}

func (ctx *Context) PollEvent() bool {
	return ctx.pollEvent(0)
}

func (ctx *Context) WaitEvent() bool {
	return ctx.pollEvent(-1)
}

func (ctx *Context) Wake() {
	ctx.looper.Wake()
}

func (ctx *Context) Name() string {
	return ctx.className
}

func (s *Sensor) Enable(act *Activity) {
	ctx := act.Context()
	if ctx.sensorQueue == nil {
		ctx.sensorQueue = SensorManagerInstance().createEventQueue(ctx.looper,
			looper_ID_SENSOR, nil, unsafe.Pointer(uintptr(looper_ID_SENSOR)))
	}
	ctx.sensorQueue.enableSensor(s)
}

func (s *Sensor) Disable(act *Activity) {
	act.Context().sensorQueue.disableSensor(s)
}

func (s *Sensor) SetEventRate(act *Activity, t time.Duration) {
	act.Context().sensorQueue.setEventRate(s, t)
}
