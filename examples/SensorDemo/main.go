package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/gooid/gl/egl"
	gl "github.com/gooid/gl/es2"
	"github.com/gooid/gooid"
	"github.com/gooid/gooid/input"
	"github.com/gooid/gooid/sensor"
	"github.com/gooid/imgui"
	"github.com/gooid/imgui/util"
)

func main() {
	context := app.Callbacks{
		WindowDraw:         draw,
		WindowRedrawNeeded: redraw,
		WindowDestroyed:    destroyed,
		Event:              event,
		Sensor:             sensorEevent,
		Create:             create,
	}

	//log.Println(ss.SENSOR_TYPE_ACCELEROMETER)
	app.SetMainCB(func(ctx *app.Context) {
		ctx.Debug(true)
		ctx.Run(context)
	})
	for app.Loop() {
	}
	log.Println("done.")
}

var mouseLeft = false
var mouseRight = false
var lastX, lastY int

func event(act *app.Activity, e *app.InputEvent) {
	if mot := e.Motion(); mot != nil {
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
const title = "Sensor "

var (
	width, height int
	eglctx        *egl.EGLContext

	im          *util.Render
	sensorInfos []sensorInfo
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
			desity, _ := strconv.Atoi(dstr)

			if desity > 160 {
				iScale := (float32)(desity / 160 / WINDOWSCALE)
				io.SetFontGlobalScale(iScale)
				style := imgui.GetStyle()
				style.ScaleAllSizes(iScale)

				scale := imgui.NewVec2((float32)(WINDOWSCALE), (float32)(WINDOWSCALE))
				defer scale.Delete()
				io.SetDisplayFramebufferScale(scale)
			}
		}
	}
	io.SetConfigFlags(io.GetConfigFlags() | int(imgui.ConfigFlags_IsTouchScreen))

	// render 只需初始化一次
	im = util.NewRender("#version 100")

	sensorInfos = SensorInfo(act)
	// 排序，在显示的时候保证顺序一致
	sort.SliceStable(sensorInfos, func(i, j int) bool {
		if sensorInfos[i].Type == sensorInfos[j].Type {
			return sensorInfos[i].Name < sensorInfos[j].Name
		}
		return sensorInfos[i].Type < sensorInfos[j].Type
	})
}

func redraw(act *app.Activity, win *app.Window) {
	act.Context().Call(func() {
		releaseEGL()
		initEGL(act, win)
	}, false)
	act.Context().Call(func() {
		draw(act, nil)
	}, false)
	act.Context().Call(func() {
		draw(act, nil)
	}, false)
}

func destroyed(act *app.Activity, win *app.Window) {
	releaseEGL()
}

func draw(act *app.Activity, win *app.Window) {
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
		curpos := imgui.NewVec2(MARGIN, MARGIN)
		defer curpos.Delete()
		imgui.SetNextWindowPos(curpos)
		isize := imgui.NewVec2(float32(width)-2*MARGIN, float32(height)-2*MARGIN)
		defer isize.Delete()
		imgui.SetNextWindowSize(isize)

		imgui.Begin(title, (*bool)(nil), int(imgui.WindowFlags_NoSavedSettings|
			//imgui.WindowFlags_NoTitleBar|
			imgui.WindowFlags_NoResize|
			imgui.WindowFlags_NoCollapse|
			imgui.WindowFlags_NoMove))

		if imgui.CollapsingHeader("sensor list", int(imgui.TreeNodeFlags_DefaultOpen)) {
			lastType := sensor.TYPE_INVALID
			for i, info := range sensorInfos {
				if info.Type != lastType {
					if imgui.TreeNodeEx(info.Type.String(), int(imgui.TreeNodeFlags_DefaultOpen)) {
						for j := i; j < len(sensorInfos) && sensorInfos[j].Type == info.Type; j++ {
							imgui.TextWrapped(fmt.Sprintf("%+v", info))
						}
						imgui.TreePop()
						imgui.Separator()
					}
					lastType = info.Type
				}
			}
		}

		if imgui.CollapsingHeader("sensor events", int(imgui.TreeNodeFlags_DefaultOpen)) {
			lastType := sensor.TYPE_INVALID
			for i, info := range sensorInfos {
				if info.Type != lastType {
					if imgui.TreeNodeEx(info.Type.String(), int(imgui.TreeNodeFlags_DefaultOpen)) {
						for j := i; j < len(sensorInfos) && info.Type == sensorInfos[j].Type; j++ {
							EnableSensor(act, sensorInfos[j].sensor)
						}

						is := []int{}
						for k, _ := range eventMap[info.Type] {
							is = append(is, k)
						}
						sort.Ints(is)
						tms, _ := eventMap[info.Type]
						for _, v := range is {
							imgui.TextWrapped(eventString(tms[v]))
						}
						imgui.TreePop()
						imgui.Separator()
					} else {
						for j := i; j < len(sensorInfos) && info.Type == sensorInfos[j].Type; j++ {
							DisableSensor(act, sensorInfos[j].sensor)
						}
					}
					lastType = info.Type
				}
			}
		} else {
			for _, info := range sensorInfos {
				DisableSensor(act, info.sensor)
			}
		}

		imgui.End()

		imgui.Render()

		// Rendering
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		//log.Println("draw:", drawData.GetCmdListsCount(),	drawData.GetTotalVtxCount(), drawData.GetTotalIdxCount())

		im.Render(imgui.GetDrawData())

		eglctx.SwapBuffers()
	}
}
