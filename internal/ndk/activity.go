package app

/*
#include <android/native_activity.h>
*/
import "C"

import "unsafe"

type Activity C.ANativeActivity

func (a *Activity) cptr() *C.ANativeActivity {
	return (*C.ANativeActivity)(a)
}

func (a *Activity) AssetManager() *AssetManager {
	return (*AssetManager)(a.cptr().assetManager)
}

func (a *Activity) Configuration() *Configuration {
	return fromAssetManager(a.AssetManager())
}

func (a *Activity) Instance() unsafe.Pointer {
	return a.cptr().instance
}

func (a *Activity) JActivity() uintptr {
	return uintptr(a.cptr().clazz)
}

func (a *Activity) InternalDataPath() string {
	return C.GoString(a.cptr().internalDataPath)
}

func (a *Activity) ExternalDataPath() string {
	return C.GoString(a.cptr().externalDataPath)
}

func (a *Activity) SdkVersion() int {
	return int(a.cptr().sdkVersion)
}

func (a *Activity) ObbPath() string {
	return C.GoString(a.cptr().obbPath)
}

/**
 * Finish the given activity.  Its finish() method will be called, causing it
 * to be stopped and destroyed.  Note that this method can be called from
 * *any* thread; it will send a message to the main thread of the process
 * where the Java finish call will take place.
 */
func (a *Activity) Finish() {
	C.ANativeActivity_finish(a.cptr())
}

/**
 * Change the window format of the given activity.  Calls getWindow().setFormat()
 * of the given activity.  Note that this method can be called from
 * *any* thread; it will send a message to the main thread of the process
 * where the Java finish call will take place.
 */
func (a *Activity) SetWindowFormat(format int) {
	C.ANativeActivity_setWindowFormat(a.cptr(), C.int32_t(format))
}

/**
 * Change the window flags of the given activity.  Calls getWindow().setFlags()
 * of the given activity.  Note that this method can be called from
 * *any* thread; it will send a message to the main thread of the process
 * where the Java finish call will take place.  See window.h for flag constants.
 */
func (a *Activity) SetWindowFlags(addFlags, removeFlags int) {
	C.ANativeActivity_setWindowFlags(a.cptr(),
		C.uint32_t(addFlags), C.uint32_t(removeFlags))
}

/**
 * Flags for ANativeActivity_showSoftInput; see the Java InputMethodManager
 * API for documentation.
 */
const (
	ACTIVITY_SHOW_SOFT_INPUT_IMPLICIT = C.ANATIVEACTIVITY_SHOW_SOFT_INPUT_IMPLICIT
	ACTIVITY_SHOW_SOFT_INPUT_FORCED   = C.ANATIVEACTIVITY_SHOW_SOFT_INPUT_FORCED
)

/**
 * Show the IME while in the given activity.  Calls InputMethodManager.showSoftInput()
 * for the given activity.  Note that this method can be called from
 * *any* thread; it will send a message to the main thread of the process
 * where the Java finish call will take place.
 */
func (a *Activity) ShowSoftInput(flags int) {
	C.ANativeActivity_showSoftInput(a.cptr(), C.uint32_t(flags))
}

/**
 * Flags for ANativeActivity_hideSoftInput; see the Java InputMethodManager
 * API for documentation.
 */
const (
	ACTIVITY_HIDE_SOFT_INPUT_IMPLICIT_ONLY = C.ANATIVEACTIVITY_HIDE_SOFT_INPUT_IMPLICIT_ONLY
	ACTIVITY_HIDE_SOFT_INPUT_NOT_ALWAYS    = C.ANATIVEACTIVITY_HIDE_SOFT_INPUT_NOT_ALWAYS
)

/**
 * Hide the IME while in the given activity.  Calls InputMethodManager.hideSoftInput()
 * for the given activity.  Note that this method can be called from
 * *any* thread; it will send a message to the main thread of the process
 * where the Java finish call will take place.
 */
func (a *Activity) HideSoftInput(flags int) {
	C.ANativeActivity_hideSoftInput(a.cptr(), C.uint32_t(flags))
}
