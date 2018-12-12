// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Go runtime entry point for apps running on android.
// Sets up everything the runtime needs and exposes
// the entry point to JNI.

// +build android

package app

/*
#cgo LDFLAGS: -llog -landroid

#include <android/log.h>
#include <android/configuration.h>
#include <android/native_activity.h>
#include <stdlib.h>
#include <time.h>
*/
import "C"

import (
	"unsafe"
)

func onCreate(act *Activity, buf []byte) {
	ctx := act.Context()
	ctx.runFunc(func() {}, true)
	info("onCreate:", act, len(buf))

	ctx.act = (*Activity)(act)
	if len(buf) > 0 {
		ctx.savedState = make([]byte, len(buf))
		copy(ctx.savedState, buf)
	}

	if ctx.Create != nil {
		ctx.runFunc(func() {
			ctx.Create(act, ctx.savedState)
		}, true)
	}
}

//export onStart
func onStart(act *Activity) {
	ctx := act.Context()
	info("onStart:", act)
	ctx.reset()
	if ctx.Start != nil {
		ctx.runFunc(func() {
			ctx.Start(ctx.act)
		}, true)
	}
}

//export onResume
func onResume(act *Activity) {
	ctx := act.Context()
	info("onResume:", act)
	if ctx.Resume != nil {
		ctx.runFunc(func() {
			ctx.Resume(ctx.act)
		}, true)
	}
	ctx.isResume = true
}

//export onPause
func onPause(act *Activity) {
	ctx := act.Context()
	info("onPause:", act)
	ctx.isResume = false
	ctx.runFunc(func() {
		if ctx.Pause != nil {
			ctx.Pause(ctx.act)
		}
	}, true)
}

//export onStop
func onStop(act *Activity) {
	ctx := act.Context()
	info("onStop:", act)
	ctx.isResume = false
	ctx.runFunc(func() {
		if ctx.Stop != nil {
			ctx.Stop(ctx.act)
		}
	}, true)
}

//export onDestroy
func onDestroy(act *Activity) {
	ctx := act.Context()
	info("onDestroy:", act)
	ctx.runFunc(func() {
		ctx.willDestory = true
		if ctx.input != nil {
			ctx.input.DetachLooper()
			ctx.input = nil
		}
		if ctx.Destroy != nil {
			ctx.Destroy(ctx.act)
		}
		ctx.Release()
	}, true)

	ctx.funcChan <- func() {
		info("onDestroy:", act, "complete")
	}
	clrActivity(act, ctx.className)
}

//export onWindowFocusChanged
func onWindowFocusChanged(act *Activity, hasFocus C.int) {
	ctx := act.Context()
	info("onWindowFocusChanged:", act, hasFocus)
	focus := hasFocus != 0
	ctx.isFocus = focus
	ctx.runFunc(func() {
		if ctx.FocusChanged != nil {
			ctx.FocusChanged(ctx.act, focus)
		}
	}, true)
}

//export onSaveInstanceState
func onSaveInstanceState(act *Activity, outSize *C.size_t) unsafe.Pointer {
	ctx := act.Context()
	info("onSaveInstanceState:", act)
	if ctx.SaveState != nil {
		ctx.runFunc(func() {
			ctx.savedState = ctx.SaveState(ctx.act)
		}, true)

		if len(ctx.savedState) > 0 {
			size := len(ctx.savedState)
			info("\t\tsize =", size)
			*outSize = C.size_t(size)
			ptr := C.malloc(C.size_t(size))
			copy((*[1 << 30]byte)(unsafe.Pointer(ptr))[:size], ctx.savedState)
			return ptr
		}
	}
	return nil
}

//export onNativeWindowCreated
func onNativeWindowCreated(act *Activity, window *Window) {
	ctx := act.Context()
	info("onNativeWindowCreated:", act, window)
	ctx.window = window
	if ctx.WindowCreated != nil {
		ctx.runFunc(func() {
			ctx.WindowCreated(act, window)
		}, true)
	}
}

//export onNativeWindowResized
func onNativeWindowResized(act *Activity, window *Window) {
	ctx := act.Context()
	info("onNativeWindowResized:", act, window)
	if ctx.WindowResized != nil {
		ctx.runFunc(func() {
			ctx.WindowResized(act, window)
		}, true)
	}
}

//export onNativeWindowRedrawNeeded
func onNativeWindowRedrawNeeded(act *Activity, window *Window) {
	ctx := act.Context()
	info("onNativeWindowRedrawNeeded:", act, window)
	assert(ctx.window == window)
	if ctx.window != window {
		ctx.window = window
	}

	if ctx.WindowRedrawNeeded != nil {
		ctx.runFunc(func() {
			ctx.WindowRedrawNeeded(act, window)
		}, true)
	}
}

//export onNativeWindowDestroyed
func onNativeWindowDestroyed(act *Activity, window *Window) {
	ctx := act.Context()
	info("onNativeWindowDestroyed:", act, window)
	ctx.runFunc(func() {
		//Info("onNativeWindowDestroyed.func")
		if ctx.WindowDestroyed != nil {
			ctx.WindowDestroyed(act, window)
		}
		ctx.window = nil
	}, true)
}

//export onInputQueueCreated
func onInputQueueCreated(act *Activity, queue *InputQueue) {
	ctx := act.Context()
	info("onInputQueueCreated:", act, queue)
	ctx.runFunc(func() {
		ctx.input = (*InputQueue)(queue)
		ctx.input.AttachLooper(ctx.looper, looper_ID_INPUT, nil, unsafe.Pointer(uintptr(looper_ID_INPUT)))
	}, true)
}

//export onInputQueueDestroyed
func onInputQueueDestroyed(act *Activity, queue *InputQueue) {
	ctx := act.Context()
	info("onInputQueueDestroyed:", act, queue)
	ctx.runFunc(func() {
		ctx.input.DetachLooper()
		ctx.input = nil
	}, true)
}

//export onContentRectChanged
func onContentRectChanged(act *Activity, rect *Rect) {
	ctx := act.Context()
	info("onContentRectChanged:", act, rect)
	if ctx.ContentRectChanged != nil {
		ctx.runFunc(func() {
			ctx.ContentRectChanged(act, rect)
		}, true)
	}
}

//export onConfigurationChanged
func onConfigurationChanged(act *Activity) {
	ctx := act.Context()
	info("onConfigurationChanged:", act)
	if ctx.ConfigurationChanged != nil {
		ctx.runFunc(func() {
			ctx.ConfigurationChanged(act)
		}, true)
	}
}

//export onLowMemory
func onLowMemory(act *Activity) {
	ctx := act.Context()
	info("onLowMemory:", act)
	if ctx.LowMemory != nil {
		ctx.runFunc(func() {
			ctx.LowMemory(act)
		}, true)
	}
}
