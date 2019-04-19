package main

import (
	"log"
	"strconv"

	"github.com/gooid/gl/egl"
	app "github.com/gooid/gooid"
	"github.com/gooid/gooid/input"
)

func main() {
	context := app.Callbacks{
		WindowDraw:         _draw,
		WindowRedrawNeeded: redraw,
		WindowDestroyed: func(act *app.Activity, win *app.Window) {
			releaseEGL()
			if eglctx != nil {
				eglctx.Terminate()
				eglctx = nil
			}
			destroyed()
		},
		Event: event,
		Create: func(act *app.Activity, _ []byte) {
			preCreate(act)
			create()
			postCreate(act)
		},
	}
	app.SetMainCB(func(ctx *app.Context) {
		ctx.Debug(true)
		ctx.Run(context)
	})
	for app.Loop() {
	}
	log.Println("done.")
}

var aMouseLeft = false

func event(act *app.Activity, e *app.InputEvent) {
	if mot := e.Motion(); mot != nil {
		x := int(float32(mot.GetX(0)) / WINDOWSCALE)
		y := int(float32(mot.GetY(0)) / WINDOWSCALE)
		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			aMouseLeft = false
			//log.Println("event:", mot)
		case input.MOTION_EVENT_ACTION_DOWN:
			aMouseLeft = true
			//log.Println("event:", mot)
		case input.MOTION_EVENT_ACTION_MOVE:

		default:
			//log.Println("event:", mot)
			return
		}
		mouseEvent(x, y, aMouseLeft)
		_draw(act, nil)

		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			mouseEvent(0, 0, aMouseLeft)
			_draw(act, nil)
		}
	}
}

var curAct *app.Activity
var eglctx *egl.EGLContext

func redraw(act *app.Activity, win *app.Window) {
	curAct = act
	act.Context().Call(func() {
		releaseEGL()
		if eglctx != nil {
			eglctx.Terminate()
			eglctx = nil
		}

		width, height = win.Width(), win.Height()
		eglctx = egl.CreateEGLContext(&nativeInfo{win: win})
		if eglctx == nil {
			return
		}
		initEGL()
	}, false)
}

func _draw(act *app.Activity, win *app.Window) {
	if eglctx != nil {
		draw()
	}
}

func getDensity() int {
	dstr := app.PropGet("hw.lcd.density")
	if dstr == "" {
		dstr = app.PropGet("qemu.sf.lcd_density")
	}

	log.Println(" lcd_density:", dstr)
	if dstr != "" {
		density, _ := strconv.Atoi(dstr)
		return density
	}
	return density
}
func SwapBuffers() { eglctx.SwapBuffers() }
func Wake()        { curAct.Context().Wake() }
