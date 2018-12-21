package main

/*
#include <stdlib.h>
#include <dlfcn.h>
*/
import "C"

import (
	"errors"
	"unsafe"

	"github.com/gooid/gooid"
)

type nativeInfo struct {
	win *app.Window
}

func (e *nativeInfo) NativeDisplay() unsafe.Pointer {
	return nil
}
func (e *nativeInfo) NativeWindow() unsafe.Pointer {
	return unsafe.Pointer(e.win)
}
func (e *nativeInfo) WindowSize() (w, h int) {
	return e.win.Width(), e.win.Height()
}

func (e *nativeInfo) SetBuffersGeometry(format int) int {
	return e.win.SetBuffersGeometry(0, 0, format)
}

func DLOpen(name string) (uintptr, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	h := uintptr(C.dlopen(cname, C.RTLD_LAZY|C.RTLD_GLOBAL))
	if h == 0 {
		return 0, errors.New(C.GoString(C.dlerror()))
	}
	return h, nil
}

func DLClose(h uintptr) {
	C.dlclose(unsafe.Pointer(h))
}
