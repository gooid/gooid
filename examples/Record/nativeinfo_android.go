package main

import (
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
