package main

import (
	"log"
	"os"
	"runtime"
	"time"

	gl "github.com/gooid/gl/es2"
	"github.com/gooid/imgui"
	"github.com/gooid/imgui/util"
)

//

const (
	WINDOWSCALE  = 1
	AUTOHIDETIME = time.Second * 5
	title        = "Camera2"
)

var (
	width   = (1080 / 8 * 3)
	height  = (1920 / 8 * 3)
	density = (480 / 8 * 3)

	im *util.Render
)

// event
var lastTouch time.Time
var mouseLeft = false
var mouseRight = false
var lastX, lastY int

func mouseEvent(x, y int, left bool) {
	lastTouch = time.Now()
	lastX = x
	lastY = y
	mouseLeft = left
}

func initEGL() {
	log.Println("WINSIZE:", width, height)
	width = int(float32(width) / WINDOWSCALE)
	height = int(float32(height) / WINDOWSCALE)

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

	CameraUI().SetWindow(width, height)
	CameraUI().SetCamera(openCamera("", 0))
}

func releaseEGL() {
	if im != nil {
		im.DestroyDeviceObjects()
		CameraUI().Close()
	}
}

func destroyed() {}

func create() {
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

	density = getDensity()

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
	// TODO libpath, cas...
	//CameraUI().Init()
	CameraUI().SetWakeup(func() {
		Wake()
	})
}

func draw() {
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

	imgui.Begin("_bg_", (*bool)(nil), int(imgui.WindowFlags_NoSavedSettings|
		imgui.WindowFlags_NoTitleBar|
		imgui.WindowFlags_NoResize|
		imgui.WindowFlags_NoMove))

	ds := imgui.GetWindowDrawList()
	ds.AddCallback(func(in interface{}) bool {
		CameraUI().DrawFrame()
		return true
	}, nil)
	CameraUI().DrawFrameMark(ds)

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

		CameraUI().DrawUI()

		/*
			imgui.Separator()
			imgui.Separator()
			if imgui.Button("Init") {
				initCamera()
			}
			imgui.SameLine()
			if imgui.Button("Close") {
				closeCamera()
			}

			imgui.Separator()
			if imgui.Button("Start") {
				startPreview()
			}

			imgui.Separator()
			if imgui.Button("TakePhoto") {
				cam.TakePhoto()
			}

			imgui.Separator()
			if imgui.Button("Stop") {
				stopPreview()
			}
		*/

		imgui.End()
	}

	imgui.Render()

	// Rendering
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	im.Render(imgui.GetDrawData())

	SwapBuffers()
}
