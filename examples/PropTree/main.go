package main

import (
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gooid/gl/egl"
	gl "github.com/gooid/gl/es2"
	"github.com/gooid/gooid"
	"github.com/gooid/gooid/input"
	"github.com/gooid/imgui"
	"github.com/gooid/imgui/util"
)

func main() {
	context := app.Callbacks{
		FocusChanged:       focus,
		WindowCreated:      create,
		WindowDraw:         draw,
		WindowRedrawNeeded: redraw,
		WindowDestroyed:    destroyed,
		Event:              event,
		//Pause:
		//Resume:

	}
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

var isFocus = false

func focus(_ *app.Activity, f bool) {
	isFocus = f
}

const WINDOWSCALE = 1
const title = "显示所有属性"

var (
	width, height int
	eglctx        *egl.EGLContext

	deltaLast = time.Now()
	im        *util.Render
)

func initEGL(win *app.Window) {
	eglctx = egl.CreateEGLContext(&nativeInfo{win: win})
	if eglctx == nil {
		return
	}
	width, height = win.Width(), win.Height()
	log.Println("WINSIZE:", width, height)
	width = int(float32(width) / WINDOWSCALE)
	height = int(float32(height) / WINDOWSCALE)

	// gl init
	gl.Init()
	log.Println("RENDERER:", gl.GoStr(gl.GetString(gl.RENDERER)))
	log.Println("VENDOR:", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Println("VERSION:", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Println("EXTENSIONS:", gl.GoStr(gl.GetString(gl.EXTENSIONS)))

	log.Printf("%s %s", gl.GoStr(gl.GetString(gl.RENDERER)), gl.GoStr(gl.GetString(gl.VERSION)))

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
			//fonts.AddFontFromFileTTF(fontName, 14.0, imgui.SwigcptrFontConfig(0), fonts.GetGlyphRangesChineseSimplifiedCommon())
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

	im = util.NewRender("#version 300 es")

	// 设置完字体再调用
	im.CreateDeviceObjects()

	dsize := imgui.NewVec2(float32(width), float32(height))
	//dsize := imgui.NewVec2(float32(win.Width()/SCALE_BUFFER), float32(win.Height()/SCALE_BUFFER))
	defer dsize.Delete()
	io.SetDisplaySize(dsize)

	// Setup time step
	deltaLast = time.Now()
	io.SetDeltaTime(0)
}

func releaseEGL() {
	if im != nil {
		im.DestroyDeviceObjects()
		im = nil
	}
	if eglctx != nil {
		eglctx.Terminate()
		eglctx = nil
	}
}

func create(act *app.Activity, win *app.Window) {
	loadProps()
}

func redraw(act *app.Activity, win *app.Window) {
	act.Context().Call(func() {
		releaseEGL()
		initEGL(win)
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
	if isFocus && eglctx != nil {
		io := imgui.GetIO()

		pos := imgui.NewVec2(float32(lastX), float32(lastY))
		defer pos.Delete()
		io.SetMousePos(pos)
		io.SetMouseDown([]bool{mouseLeft, false, mouseRight, false, false})

		// Setup time step
		now := time.Now()
		io.SetDeltaTime(float32(now.Sub(deltaLast).Seconds()))
		deltaLast = now

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
		//imgui.Text(title)

		drawPropTree()

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

type prop struct {
	k, v string
}

var props []prop

func loadProps() {
	app.PropVisit(func(k, v string) {
		props = append(props, prop{k: k, v: v})
	})
	sort.SliceStable(props, func(i, j int) bool {
		return strings.Compare(props[i].k, props[j].k) < 0
	})
}

func drawPropTree() {
	lastProp := []string{}
	exposeProp := NewStack()
	for _, prop := range props {
		func(k, v string) {
			curProp := strings.Split(k, ".")
			defer func() {
				lastProp = curProp
			}()

			sameCnt := 0
			exposeDepth := exposeProp.Depth()

			for i := 0; i < len(curProp)-1 && i < exposeDepth; i++ {
				if exposeProp.Poll(exposeDepth-i-1) == curProp[i] {
					sameCnt++
				} else {
					break
				}
			}
			// TreePop
			for i := sameCnt; i < exposeDepth; i++ {
				exposeProp.Pop()
				//fmt.Println(" ****Pop:", v)
				imgui.TreePop()
				imgui.Separator()
			}

			// 最后一个域直接显示Text
			if sameCnt == len(curProp)-1 {
				imgui.Text(curProp[len(curProp)-1] + ": " + v)
			} else if sameCnt < len(curProp)-1 {
				if sameCnt >= len(lastProp) || lastProp[sameCnt] != curProp[sameCnt] {

					// 新的 Tree
					for i := sameCnt; i < len(curProp)-1; i++ {
						//fmt.Println("     +Node:", curProp[i])
						if imgui.TreeNode(curProp[i]) {
							//fmt.Println(" ----Push:", curProp[i])
							exposeProp.Push(curProp[i])
						} else {
							break
						}
					}
					if exposeProp.Depth() == len(curProp)-1 {
						imgui.Text(curProp[len(curProp)-1] + ": " + v)
					}
				}
			}
		}(prop.k, prop.v)
	}

	// TreePop
	depth := exposeProp.Depth()
	for i := 0; i < depth; i++ {
		exposeProp.Pop()
		//fmt.Println(" ----Pop:", v)
		imgui.TreePop()
		imgui.Separator()
	}
}
