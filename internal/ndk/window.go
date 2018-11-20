package app

/*
#include <android/window.h>
#include <android/native_window_jni.h>
*/
import "C"

import (
	"unsafe"
)

/*
 * Pixel formats that a window can use.
 */
const (
	WINDOW_FORMAT_RGBA_8888 = C.WINDOW_FORMAT_RGBA_8888
	WINDOW_FORMAT_RGBX_8888 = C.WINDOW_FORMAT_RGBX_8888
	WINDOW_FORMAT_RGB_565   = C.WINDOW_FORMAT_RGB_565
)

const (
	FLAG_ALLOW_LOCK_WHILE_SCREEN_ON = C.AWINDOW_FLAG_ALLOW_LOCK_WHILE_SCREEN_ON
	FLAG_DIM_BEHIND                 = C.AWINDOW_FLAG_DIM_BEHIND
	FLAG_BLUR_BEHIND                = C.AWINDOW_FLAG_BLUR_BEHIND
	FLAG_NOT_FOCUSABLE              = C.AWINDOW_FLAG_NOT_FOCUSABLE
	FLAG_NOT_TOUCHABLE              = C.AWINDOW_FLAG_NOT_TOUCHABLE
	FLAG_NOT_TOUCH_MODAL            = C.AWINDOW_FLAG_NOT_TOUCH_MODAL
	FLAG_TOUCHABLE_WHEN_WAKING      = C.AWINDOW_FLAG_TOUCHABLE_WHEN_WAKING
	FLAG_KEEP_SCREEN_ON             = C.AWINDOW_FLAG_KEEP_SCREEN_ON
	FLAG_LAYOUT_IN_SCREEN           = C.AWINDOW_FLAG_LAYOUT_IN_SCREEN
	FLAG_LAYOUT_NO_LIMITS           = C.AWINDOW_FLAG_LAYOUT_NO_LIMITS
	FLAG_FULLSCREEN                 = C.AWINDOW_FLAG_FULLSCREEN
	FLAG_FORCE_NOT_FULLSCREEN       = C.AWINDOW_FLAG_FORCE_NOT_FULLSCREEN
	FLAG_DITHER                     = C.AWINDOW_FLAG_DITHER
	FLAG_SECURE                     = C.AWINDOW_FLAG_SECURE
	FLAG_SCALED                     = C.AWINDOW_FLAG_SCALED
	FLAG_IGNORE_CHEEK_PRESSES       = C.AWINDOW_FLAG_IGNORE_CHEEK_PRESSES
	FLAG_LAYOUT_INSET_DECOR         = C.AWINDOW_FLAG_LAYOUT_INSET_DECOR
	FLAG_ALT_FOCUSABLE_IM           = C.AWINDOW_FLAG_ALT_FOCUSABLE_IM
	FLAG_WATCH_OUTSIDE_TOUCH        = C.AWINDOW_FLAG_WATCH_OUTSIDE_TOUCH
	FLAG_SHOW_WHEN_LOCKED           = C.AWINDOW_FLAG_SHOW_WHEN_LOCKED
	FLAG_SHOW_WALLPAPER             = C.AWINDOW_FLAG_SHOW_WALLPAPER
	FLAG_TURN_SCREEN_ON             = C.AWINDOW_FLAG_TURN_SCREEN_ON
	FLAG_DISMISS_KEYGUARD           = C.AWINDOW_FLAG_DISMISS_KEYGUARD
)

type WindowBuffer C.ANativeWindow_Buffer

func (b *WindowBuffer) cptr() *C.ANativeWindow_Buffer {
	return (*C.ANativeWindow_Buffer)(b)
}

func (b *WindowBuffer) Width() int {
	return int(b.width)
}

func (b *WindowBuffer) Height() int {
	return int(b.height)
}

func (b *WindowBuffer) Stride() int {
	return int(b.stride)
}

func (b *WindowBuffer) Format() int {
	return int(b.format)
}

func (b *WindowBuffer) Bits() []byte {
	return ((*[1 << 30]byte)(unsafe.Pointer(b.bits)))[:b.stride*b.height]
}

func (b *WindowBuffer) Bit16s() []uint16 {
	return ((*[1 << 28]uint16)(unsafe.Pointer(b.bits)))[:b.stride*b.height]
}

func (b *WindowBuffer) Bit32s() []uint32 {
	return ((*[1 << 28]uint32)(unsafe.Pointer(b.bits)))[:b.stride*b.height]
}

type Window C.ANativeWindow

func (w *Window) cptr() *C.ANativeWindow {
	return (*C.ANativeWindow)(w)
}

func (w *Window) Pointer() *C.ANativeWindow {
	return (*C.ANativeWindow)(w)
}

/**
 * Return the ANativeWindow associated with a Java Surface object,
 * for interacting with it through native code.  This acquires a reference
 * on the ANativeWindow that is returned; be sure to use ANativeWindow_release()
 * when done with it so that it doesn't leak.
 */
func FromSurface(env *C.JNIEnv, surface C.jobject) *Window {
	return (*Window)(C.ANativeWindow_fromSurface(env, surface))
}

/**
 * Acquire a reference on the given ANativeWindow object.  This prevents the object
 * from being deleted until the reference is removed.
 */
func (w *Window) Acquire() {
	C.ANativeWindow_acquire(w.cptr())
}

/**
 * Remove a reference that was previously acquired with ANativeWindow_acquire().
 */
func (w *Window) Release() {
	C.ANativeWindow_release(w.cptr())
}

/*
 * Return the current width in pixels of the window surface.  Returns a
 * negative value on error.
 */
func (w *Window) Width() int {
	return int(C.ANativeWindow_getWidth(w.cptr()))
}

/*
 * Return the current height in pixels of the window surface.  Returns a
 * negative value on error.
 */
func (w *Window) Height() int {
	return int(C.ANativeWindow_getHeight(w.cptr()))
}

/*
 * Return the current pixel format of the window surface.  Returns a
 * negative value on error.
 */
func (w *Window) Format() int {
	return int(C.ANativeWindow_getFormat(w.cptr()))
}

/*
 * Change the format and size of the window buffers.
 *
 * The width and height control the number of pixels in the buffers, not the
 * dimensions of the window on screen.  If these are different than the
 * window's physical size, then it buffer will be scaled to match that size
 * when compositing it to the screen.
 *
 * For all of these parameters, if 0 is supplied then the window's base
 * value will come back in force.
 *
 * width and height must be either both zero or both non-zero.
 *
 */
func (w *Window) SetBuffersGeometry(width, height, format int) int {
	return int(C.ANativeWindow_setBuffersGeometry(w.cptr(), C.int32_t(width), C.int32_t(height), C.int32_t(format)))
}

/**
 * Lock the window's next drawing surface for writing.
 * inOutDirtyBounds is used as an in/out parameter, upon entering the
 * function, it contains the dirty region, that is, the region the caller
 * intends to redraw. When the function returns, inOutDirtyBounds is updated
 * with the actual area the caller needs to redraw -- this region is often
 * extended by ANativeWindow_lock.
 */
func (w *Window) Lock(inDirtyBounds Rect) (*WindowBuffer, Rect, bool) {
	var outBuffer WindowBuffer
	ret := C.ANativeWindow_lock(w.cptr(), outBuffer.cptr(), inDirtyBounds.cptr())
	return &outBuffer, inDirtyBounds, ret == 0
}

/**
 * Unlock the window's drawing surface after previously locking it,
 * posting the new buffer to the display.
 */
func (w *Window) UnlockAndPost() bool {
	return 0 == C.ANativeWindow_unlockAndPost(w.cptr())
}
