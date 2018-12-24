//

package main

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func CursorPosCallback(w *glfw.Window, x float64, y float64) {
	lastX = int(float32(x) / WINDOWSCALE)
	lastY = int(float32(y) / WINDOWSCALE)
}

func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		switch action {
		case glfw.Press, glfw.Repeat:
			mouseLeft = true
		case glfw.Release:
			mouseLeft = false
		}

		draw()
	}
}

func FocusCallback(w *glfw.Window, focused bool) {
	if !focused {
		//destroyed()
	} else {
		//redraw()
	}
}

func main() {
	runtime.LockOSThread()
	log.Printf("main ...")

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// 以下是指定用 EGL 和 OpenGL ES
	glfw.WindowHint(glfw.ContextCreationAPI, glfw.EGLContextAPI)
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)

	w, err := glfw.CreateWindow(width, height, "Record", nil, nil)
	if err != nil {
		panic(err)
	}
	w.MakeContextCurrent()
	curWin = w

	width, height = w.GetSize()

	create()
	initEGL()

	w.SetCursorPosCallback(CursorPosCallback)
	w.SetMouseButtonCallback(MouseButtonCallback)
	w.SetFocusCallback(FocusCallback)

	for !w.ShouldClose() {
		glfw.WaitEvents()

		draw()
	}

	releaseEGL()
	log.Printf("done")
}

var curWin *glfw.Window

const RECORDPATH = "./sdcard/records"

func getDensity() int { return density }
func SwapBuffers()    { curWin.SwapBuffers() }
func Wake()           { glfw.PostEmptyEvent() }
