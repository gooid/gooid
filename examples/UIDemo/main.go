package main

import (
	"log"

	"github.com/gooid/gl/egl"
	gl "github.com/gooid/gl/es2"
	"github.com/gooid/gooid"
	"github.com/gooid/gooid/input"
)

func main() {
	context := app.Callbacks{
		FocusChanged:       focus,
		WindowDraw:         draw,
		WindowRedrawNeeded: redraw,
		WindowDestroyed:    destroyed,
		Event:              event,
	}
	app.SetMainCB(func(ctx *app.Context) {
		ctx.Run(context)
	})
	for app.Loop() {
	}
	log.Println("done.")
}

func event(act *app.Activity, e *app.InputEvent) {
	if mot := e.Motion(); mot != nil {
		if mot.GetAction() == input.MOTION_EVENT_ACTION_UP {
			log.Println("event:", mot)
			if index == 0 {
				index = 1
			} else {
				index = 0
			}
			vVertices = verticess[index]
			draw(act, nil)
		}
	}
}

var isFocus = false

func focus(_ *app.Activity, f bool) {
	isFocus = f
}

var (
	width, height int
	program       *Shader
	eglctx        *egl.EGLContext

	index     int
	vVertices = verticess[index]

	verticess = [2][]float32{
		{0.0, 0.5, 0.0,
			-0.5, -0.5, 0.0,
			0.5, -0.5, 0.0},
		{0.0, 0.5, 0.0,
			-1.0, -0.5, 0.0,
			1.0, -0.5, 0.0}}
)

func initEGL(win *app.Window) {
	eglctx = egl.CreateEGLContext(&nativeInfo{win: win})
	if eglctx == nil {
		return
	}

	width, height = win.Width(), win.Height()
	log.Println("WINSIZE:", width, height)

	gl.Init()
	log.Println("RENDERER:", gl.GoStr(gl.GetString(gl.RENDERER)))
	log.Println("VENDOR:", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Println("VERSION:", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Println("EXTENSIONS:", gl.GoStr(gl.GetString(gl.EXTENSIONS)))

	log.Printf("%s %s", gl.GoStr(gl.GetString(gl.RENDERER)), gl.GoStr(gl.GetString(gl.VERSION)))

	var err error
	// Create the program object
	program = CreateProgram(vShaderStr, fShaderStr)
	if program == nil || err != nil {
		log.Println("Create:", err)
	}
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.ClearColor(1.0, 1.0, 1.0, 0.0)
}

func releaseEGL() {
	if program != nil {
		program.Delete()
	}
	if eglctx != nil {
		eglctx.Terminate()
		eglctx = nil
	}
}

func redraw(act *app.Activity, win *app.Window) {
	releaseEGL()
	initEGL(win)
}

func destroyed(act *app.Activity, win *app.Window) {
	releaseEGL()
}

func draw(act *app.Activity, win *app.Window) {
	if isFocus && eglctx != nil {
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.ClearColor(1.0, 1.0, 1.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		if program != nil {
			program.Use()
		}

		// Load the vertex data
		a := program.GetAttribLocation("a_position")

		gl.EnableVertexAttribArray(uint32(a))
		gl.VertexAttribPointer(uint32(a), 3, gl.FLOAT, false, 0, gl.Ptr(vVertices))

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		eglctx.SwapBuffers()
	}
}

const vShaderStr = `
attribute vec4 a_position;
void main()
{
   gl_Position = a_position;
}
`

const fShaderStr = `
precision mediump float;
void main()
{
  gl_FragColor = vec4( 1.0, 0.0, 0.0, 1.0 );
}
`

type Shader struct {
	*gl.Shader
}

func CreateProgram(vertexShader, fragmentShader string) *Shader {
	s := &Shader{}

	s.Shader = gl.NewShader([]string{vertexShader + "\x00"},
		[]string{fragmentShader + "\x00"})
	return s
}

func (s *Shader) Use() {
	gl.UseProgram(s.Handle())
}

func (s *Shader) GetAttribLocation(n string) int32 {
	return gl.GetAttribLocation(s.Handle(), gl.Str(n+"\x00"))
}
