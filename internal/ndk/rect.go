package app

/*
#include <android/rect.h>
*/
import "C"

//import "unsafe"
type Rect C.ARect

func (rc *Rect) cptr() *C.ARect {
	return (*C.ARect)(rc)
}

func NewRect(l, t, r, b int) Rect {
	var rc Rect
	rc.left = C.int32_t(l)
	rc.top = C.int32_t(t)
	rc.right = C.int32_t(r)
	rc.bottom = C.int32_t(b)
	return rc
}
