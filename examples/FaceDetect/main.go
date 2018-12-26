package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gooid/gl/egl"
	gl "github.com/gooid/gl/es2"
	"github.com/gooid/gooid"
	"github.com/gooid/gooid/camera"
	"github.com/gooid/gooid/input"
	"github.com/gooid/imgui"
	"github.com/gooid/imgui/util"
)

func main() {
	context := app.Callbacks{
		WindowDraw: draw,
		//WindowCreated:      winCreate,
		WindowRedrawNeeded: redraw,
		WindowDestroyed:    destroyed,
		Event:              event,
		Create:             create,
	}
	app.SetMainCB(func(ctx *app.Context) {
		ctx.Debug(true)
		ctx.Run(context)
	})
	for app.Loop() {
	}
	log.Println("done.")
}

var lastTouch time.Time
var mouseLeft = false
var mouseRight = false
var lastX, lastY int

func event(act *app.Activity, e *app.InputEvent) {
	if mot := e.Motion(); mot != nil {
		lastTouch = time.Now()

		lastX = int(float32(mot.GetX(0)) / WINDOWSCALE)
		lastY = int(float32(mot.GetY(0)) / WINDOWSCALE)
		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			mouseLeft = false
			//log.Println("event:", mot)
		case input.MOTION_EVENT_ACTION_DOWN:
			mouseLeft = true
			//log.Println("event:", mot)
		case input.MOTION_EVENT_ACTION_MOVE:

		default:
			//log.Println("event:", mot)
			return
		}
		draw(act, nil)

		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			lastX, lastY = 0, 0
			draw(act, nil)
		}
	}
}

const WINDOWSCALE = 1
const AUTOHIDETIME = time.Second * 5

const title = "Camera"

var (
	width, height int
	density       = 160
	eglctx        *egl.EGLContext

	im      *util.Render
	cam     *cameraObj
	flashOn bool
)

func initEGL(act *app.Activity, win *app.Window) {
	width, height = win.Width(), win.Height()
	log.Println("WINSIZE:", width, height)
	width = int(float32(width) / WINDOWSCALE)
	height = int(float32(height) / WINDOWSCALE)

	eglctx = egl.CreateEGLContext(&nativeInfo{win: win})
	if eglctx == nil {
		return
	}

	log.Println("RENDERER:", gl.GoStr(gl.GetString(gl.RENDERER)))
	log.Println("VENDOR:", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Println("VERSION:", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Println("EXTENSIONS:", gl.GoStr(gl.GetString(gl.EXTENSIONS)))

	log.Printf("%s %s", gl.GoStr(gl.GetString(gl.RENDERER)), gl.GoStr(gl.GetString(gl.VERSION)))

	// 设置完字体再调用
	im.CreateDeviceObjects()

	io := imgui.GetIO()

	dsize := imgui.NewVec2(float32(width), float32(height))
	defer dsize.Delete()
	io.SetDisplaySize(dsize)

	// Setup time step
	io.SetDeltaTime(float32(time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)))

	if cam == nil {
		cam = cameraInit(camera.FACING_BACK, func(w, h int, img []byte) bool {
			act.Context().Wake()
			return false
		})
	}
	cam.ResetRender()
}

func releaseEGL() {
	if im != nil {
		im.DestroyDeviceObjects()
	}
	if eglctx != nil {
		eglctx.Terminate()
		eglctx = nil
	}
}

func create(act *app.Activity, _ []byte) {
	// gl init
	gl.Init()

	// Setup Dear ImGui binding
	imgui.CreateContext()

	// Setup style
	imgui.StyleColorsDark()
	//imgui.StyleColorsClassic()

	imgui.GetStyle().SetAlpha(0.6)

	io := imgui.GetIO()
	fonts := io.GetFonts()

	log.Println("GOOS", runtime.GOOS)
	if runtime.GOOS == "android" {
		fontName := "/system/fonts/DroidSansFallback.ttf"
		if _, err := os.Stat(fontName); err != nil {
			fontName = "/system/fonts/NotoSansCJK-Regular.ttc"
			if _, err := os.Stat(fontName); err != nil {
				fontName = "/system/fonts/DroidSans.ttf"
				if _, err := os.Stat(fontName); err != nil {
					fontName = ""
				}
			}
		}
		if fontName != "" {
			_ = fonts
			// 加载所有中文Glyph，但内存开销太大
			//fonts.AddFontFromFileTTF(fontName, 24.0, imgui.SwigcptrFontConfig(0), fonts.GetGlyphRangesChineseSimplifiedCommon())
			// 仅仅加载需要显示的中文Glyph
			fonts.AddFontFromFileTTF(fontName, 24.0, imgui.SwigcptrFontConfig(0), util.GetFontGlyphRanges(title))
		}
	}

	if runtime.GOOS == "android" {
		dstr := app.PropGet("hw.lcd.density")
		if dstr == "" {
			dstr = app.PropGet("qemu.sf.lcd_density")
		}

		log.Println(" lcd_density:", dstr)
		if dstr != "" {
			density, _ = strconv.Atoi(dstr)
		}
	}

	// 通过调整 Style 中元素大小，来控制显示大小，但同时要调整字体大小
	if density > 160 {
		iScale := float32(density) / 160 / float32(WINDOWSCALE)
		io.SetFontGlobalScale(iScale)
		style := imgui.GetStyle()
		style.ScaleAllSizes(iScale)
	}

	// 通过缩放 DisplayFramebuffer 来控制显示大小
	scale := imgui.NewVec2((float32)(WINDOWSCALE), (float32)(WINDOWSCALE))
	defer scale.Delete()
	io.SetDisplayFramebufferScale(scale)

	io.SetConfigFlags(io.GetConfigFlags() | int(imgui.ConfigFlags_IsTouchScreen))

	// render 只需初始化一次
	im = util.NewRender("#version 100")

	faceDetect.LoadOpenCV(act)

	lastTouch = time.Now()
}

func redraw(act *app.Activity, win *app.Window) {
	act.Context().Call(func() {
		releaseEGL()
		initEGL(act, win)
	}, false)
	act.Context().Call(func() {
		draw(act, nil)
	}, false)
}

func destroyed(act *app.Activity, win *app.Window) {
	releaseEGL()
}

func draw(act *app.Activity, _ *app.Window) {
	if eglctx != nil {
		io := imgui.GetIO()

		pos := imgui.NewVec2(float32(lastX), float32(lastY))
		defer pos.Delete()
		io.SetMousePos(pos)
		io.SetMouseDown([]bool{mouseLeft, false, mouseRight, false, false})

		// Setup time step
		io.SetDeltaTime(float32(time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)))

		// Margin
		MARGIN := float32(width / 20)

		imgui.NewFrame()

		// 设置为全屏
		pos = imgui.NewVec2(0, 0)
		defer pos.Delete()
		imgui.SetNextWindowPos(pos)
		size := imgui.NewVec2(width, height)
		defer size.Delete()
		imgui.SetNextWindowSize(size)

		imgui.Begin("_Camera_", (*bool)(nil), int(imgui.WindowFlags_NoSavedSettings|
			imgui.WindowFlags_NoTitleBar|
			imgui.WindowFlags_NoResize|
			imgui.WindowFlags_NoMove))

		ds := imgui.GetWindowDrawList()
		ds.AddCallback(func(in interface{}) bool {
			cam.DrawImage()
			return true
		}, nil)

		// 绘制检测到的人脸所在区域
		c := imgui.NewColor(127, 255, 255).U32()
		for _, r := range faceDetect.LastRects() {
			tl := imgui.NewVec2(cam.irender.Pos(r.X, r.Y))
			defer tl.Delete()
			t2 := imgui.NewVec2(cam.irender.Pos(r.X+r.Width, r.Y+r.Height))
			defer t2.Delete()
			ds.AddRect(tl, t2, c, 0, imgui.DrawCornerFlags_All, 6)
		}

		imgui.End()

		if time.Now().Sub(lastTouch) < AUTOHIDETIME {
			curpos := imgui.NewVec2(MARGIN, MARGIN)
			defer curpos.Delete()
			imgui.SetNextWindowPos(curpos)
			isize := imgui.NewVec2(float32(width)-2*MARGIN, float32(height)-2*MARGIN)
			defer isize.Delete()
			imgui.SetNextWindowSize(isize)

			imgui.Begin(title, (*bool)(nil), int(imgui.WindowFlags_NoSavedSettings|
				imgui.WindowFlags_NoTitleBar|
				imgui.WindowFlags_NoResize|
				imgui.WindowFlags_NoMove))

			cam.Draw()

			imgui.End()
		}

		imgui.Render()

		// Rendering
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		im.Render(imgui.GetDrawData())

		eglctx.SwapBuffers()
	}
}
