// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package alc provides OpenAL's ALC (Audio Library Context) bindings for Go.
package alc

import (
	"fmt"

	"github.com/gooid/audio/al"
)

// Error returns one of these values.
const (
	InvalidDevice  = 0xA001
	InvalidContext = 0xA002
	InvalidEnum    = 0xA003
	InvalidValue   = 0xA004
	OutOfMemory    = 0xA005

	CaptureSamples = 0x312
)

// Device represents an audio device.
type Device = al.Device

// Context represents a context created in the OpenAL layer. A valid current
// context is required to run OpenAL functions.
// The returned context will be available process-wide if it's made the
// current by calling MakeContextCurrent.
type Context = al.Context

// Open opens a new device in the OpenAL layer.
func Open(name string) *Device {
	return al.Open(name)
}

// GetCurrentContext get current context.
func GetCurrentContext() *Context {
	return al.GetCurrentContext()
}

type CaptureDevice = al.CaptureDevice

func CaptureOpen(name string, frequency uint, format int, buffersize int64) *CaptureDevice {
	return al.CaptureOpen(name, frequency, format, buffersize)
}

func Error(e int32) string {
	switch e {
	case 0:
		return ""
	case InvalidDevice:
		return "Invalid device"
	case InvalidContext:
		return "Invalid context"
	case InvalidEnum:
		return "Invalid enum"
	case InvalidValue:
		return "Invalid value"
	case OutOfMemory:
		return "Out of memory"
	}
	return fmt.Sprintf("unknow error(%d)", e)
}
