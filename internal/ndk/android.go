// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Go runtime entry point for apps running on android.
// Sets up everything the runtime needs and exposes
// the entry point to JNI.

package app

/*
#cgo LDFLAGS: -llog -landroid

#include <android/log.h>
#include <android/configuration.h>
#include <android/native_activity.h>
#include <stdlib.h>
#include <time.h>
#include <dlfcn.h>

extern int cgoLooperCallback(int p0, int p1, void* p2);
extern void* onCreate(ANativeActivity* p0, void* p1, size_t p2);
extern void onStart(ANativeActivity* p0);
extern void onResume(ANativeActivity* p0);
extern void onPause(ANativeActivity* p0);
extern void onStop(ANativeActivity* p0);
extern void onDestroy(ANativeActivity* p0);
extern void onWindowFocusChanged(ANativeActivity* p0, int p1);
extern void* onSaveInstanceState(ANativeActivity* p0, size_t* p1);
extern void onNativeWindowCreated(ANativeActivity* p0, ANativeWindow* p1);
extern void onNativeWindowResized(ANativeActivity* p0, ANativeWindow* p1);
extern void onNativeWindowRedrawNeeded(ANativeActivity* p0, ANativeWindow* p1);
extern void onNativeWindowDestroyed(ANativeActivity* p0, ANativeWindow* p1);
extern void onInputQueueCreated(ANativeActivity* p0, AInputQueue* p1);
extern void onInputQueueDestroyed(ANativeActivity* p0, AInputQueue* p1);
extern void onContentRectChanged(ANativeActivity* p0, ARect* p1);
extern void onConfigurationChanged(ANativeActivity* p0);
extern void onLowMemory(ANativeActivity* p0);

static jint _JNI_OnLoad(JavaVM* vm, void* reserved) {
	JNIEnv* env;
	if ((*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6) != JNI_OK) {
		return -1;
	}
	return JNI_VERSION_1_6;
}

static const char* callGetStringMethod(JNIEnv *env, jclass jobj, const char* method) {
	jstring jpath;
	jclass clazz = (*env)->GetObjectClass(env, jobj);
	jmethodID m = (*env)->GetMethodID(env, clazz, method, "()Ljava/lang/String;");
	if (m == 0) {
		(*env)->ExceptionClear(env);
		return NULL;
	}
	jpath = (jstring)(*env)->CallObjectMethod(env, jobj, m, NULL);
	return (*env)->GetStringUTFChars(env, jpath, NULL);
}

static const char* _GetPackageName(JNIEnv *env, jobject jobj) {
	return callGetStringMethod(env, jobj, "getPackageName");
}

static const char* _GetLocalClassName(JNIEnv *env, jobject jobj) {
	return callGetStringMethod(env, jobj, "getLocalClassName");
}

static void* _GetMainPC() { return dlsym(RTLD_DEFAULT, "main.main"); }

static void* _SetActivityCallbacks(ANativeActivity* activity) {
	activity->callbacks->onStart = onStart;
	activity->callbacks->onResume = onResume;
	activity->callbacks->onSaveInstanceState = onSaveInstanceState;
	activity->callbacks->onPause = onPause;
	activity->callbacks->onStop = onStop;
	activity->callbacks->onDestroy = onDestroy;
	activity->callbacks->onWindowFocusChanged = onWindowFocusChanged;
	activity->callbacks->onNativeWindowCreated = onNativeWindowCreated;
	activity->callbacks->onNativeWindowResized = onNativeWindowResized;
	activity->callbacks->onNativeWindowRedrawNeeded = onNativeWindowRedrawNeeded;
	activity->callbacks->onNativeWindowDestroyed = onNativeWindowDestroyed;
	activity->callbacks->onInputQueueCreated = onInputQueueCreated;
	activity->callbacks->onInputQueueDestroyed = onInputQueueDestroyed;
	activity->callbacks->onContentRectChanged = (void*)onContentRectChanged;
	activity->callbacks->onConfigurationChanged = onConfigurationChanged;
	activity->callbacks->onLowMemory = onLowMemory;
	return (void*)activity;
}

*/
import "C"

import (
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/gooid/gooid/internal/callfn"
)

var appContext struct {
	mainPC    unsafe.Pointer
	funcChan  chan func()
	actMaps   map[string]*Activity
	actMainCB func(*Context)
}

// SetMainCB
func SetMainCB(fn func(*Context)) {
	if fn == nil {
		fatal("SetMainCB(nil) is incorrect")
	}
	if appContext.actMainCB != nil {
		info("MainCB is ready")
		return
	}
	appContext.funcChan <- func() { appContext.actMainCB = fn }
	appContext.funcChan <- func() {}
}

// Loop
// applation will close, return false
func Loop() bool {
	if fn, ok := <-appContext.funcChan; ok {
		fn()
		return true
	}
	return false
}

func addActivity(act *Activity, n string) {
	appContext.actMaps[n] = act
}

func clrActivity(act *Activity, n string) {
	delete(appContext.actMaps, n)
	if len(appContext.actMaps) == 0 {
		close(appContext.funcChan)
	}
}

//export JNI_OnLoad
func JNI_OnLoad(vm *C.JavaVM, reserved unsafe.Pointer) C.jint {
	return C._JNI_OnLoad(vm, reserved)
}

func callMain() {
	if appContext.mainPC != nil {
		info("main is runing")
		return
	}

	var waitMainCB func()
	appContext.actMaps = make(map[string]*Activity)
	appContext.funcChan, waitMainCB = make(chan func(), 1), func() { (<-appContext.funcChan)() }
	appContext.mainPC = C._GetMainPC()
	if appContext.mainPC == nil {
		fatal("missing main.main")
	}

	for _, name := range []string{"TMPDIR", "PATH", "LD_LIBRARY_PATH", "BOOTCLASSPATH"} {
		n := C.CString(name)
		os.Setenv(name, C.GoString(C.getenv(n)))
		C.free(unsafe.Pointer(n))
	}

	// Set timezone.
	//
	// Note that Android zoneinfo is stored in /system/usr/share/zoneinfo,
	// but it is in some kind of packed TZiff file that we do not support
	// yet. As a stopgap, we build a fixed zone using the tm_zone name.
	var curtime C.time_t
	var curtm C.struct_tm
	C.time(&curtime)
	C.localtime_r(&curtime, &curtm)
	tzOffset := int(curtm.tm_gmtoff)
	tz := C.GoString(curtm.tm_zone)
	time.Local = time.FixedZone(tz, tzOffset)

	go func() {
		runtime.LockOSThread()
		callfn.CallFn(uintptr(appContext.mainPC))
		os.Exit(0)
	}()

	// 此处会挂起直到SetMainCB完成
	waitMainCB()
}

func createActivity(act *Activity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	lname := C.GoString(C._GetLocalClassName(act.env, act.clazz))
	pname := C.GoString(C._GetPackageName(act.env, act.clazz))
	info("ANativeActivity_onCreate:", pname+"/"+lname)

	callMain()

	C._SetActivityCallbacks(act.cptr())

	buf := []byte{}
	if savedStateSize > 0 {
		buf = (*[1 << 30]byte)(unsafe.Pointer(savedState))[:savedStateSize]
	}

	ctx := &Context{}
	act.instance = unsafe.Pointer(ctx)
	ctx.init()

	ctx.className = lname
	ctx.packageName = pname
	addActivity(act, ctx.className)
	appContext.funcChan <- func() {
		go func() {
			runtime.LockOSThread()
			appContext.actMainCB(ctx)
		}()
	}

	// 等待消息队列初始化（ctx.begin）完成
	ctx.doFunc()
	onCreate(act, buf)
}

//export ANativeActivity_onCreate
func ANativeActivity_onCreate(act *Activity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	createActivity(act, savedState, savedStateSize)
}

//export ANativeActivity_onCreateB
func ANativeActivity_onCreateB(act *Activity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	info("ANativeActivity_onCreateB...")

	createActivity(act, savedState, savedStateSize)
}
