package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gooid/gl/egl"
	gl "github.com/gooid/gl/es2"
	"github.com/gooid/gooid"
	"github.com/gooid/gooid/input"
	"github.com/gooid/gooid/storage"
	"github.com/gooid/imgui"
	"github.com/gooid/imgui/util"
)

func main() {
	context := app.Callbacks{
		//WindowDraw:         draw_,
		WindowRedrawNeeded: redraw,
		WindowDestroyed:    destroyed,
		Event:              event,
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
		draw_(act, nil)

		switch mot.GetAction() {
		case input.MOTION_EVENT_ACTION_UP:
			lastX, lastY = 0, 0
			draw_(act, nil)
		}
	}
}

const WINDOWSCALE = 1
const title = "Asset 访问"

var (
	width, height int
	eglctx        *egl.EGLContext

	deltaLast  = time.Now()
	im         *util.Render
	imagData   image.Image
	imageTxtId uint32
	imageScale float32 = 5
)

func initEGL(act *app.Activity, win *app.Window) {
	eglctx = egl.CreateEGLContext(&nativeInfo{win: win})
	if eglctx == nil {
		return
	}
	width, height = win.Width(), win.Height()
	log.Println("WINSIZE:", width, height)
	width = int(float32(width) / WINDOWSCALE)
	height = int(float32(height) / WINDOWSCALE)

	if im == nil {
		// gl init
		gl.Init()
		log.Println("RENDERER:", gl.GoStr(gl.GetString(gl.RENDERER)))
		log.Println("VENDOR:", gl.GoStr(gl.GetString(gl.VENDOR)))
		log.Println("VERSION:", gl.GoStr(gl.GetString(gl.VERSION)))
		log.Println("EXTENSIONS:", gl.GoStr(gl.GetString(gl.EXTENSIONS)))

		log.Printf("%s %s", gl.GoStr(gl.GetString(gl.RENDERER)), gl.GoStr(gl.GetString(gl.VERSION)))

		im = util.NewRender("#version 100")

		imagData = ReadImage(act, "label_icon.png")
	}
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

	// 设置完字体再调用
	im.CreateDeviceObjects()

	if rgba, ok := imagData.(*image.RGBA); ok {
		imageTxtId = genTexture(gl.Ptr(rgba.Pix), rgba.Bounds().Dx(), rgba.Bounds().Dy())
	}

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
	}
	if imageTxtId != 0 {
		gl.DeleteTextures(1, &imageTxtId)
	}
	if eglctx != nil {
		eglctx.Terminate()
		eglctx = nil
	}
}

func redraw(act *app.Activity, win *app.Window) {
	act.Context().Call(func() {
		releaseEGL()
		initEGL(act, win)
	}, false)
	act.Context().Call(func() {
		draw_(act, nil)
	}, false)
	act.Context().Call(func() {
		draw_(act, nil)
	}, false)
}

func destroyed(act *app.Activity, win *app.Window) {
	releaseEGL()
}

func draw_(act *app.Activity, win *app.Window) {
	if eglctx != nil {
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

		if imgui.CollapsingHeader("image", int(imgui.TreeNodeFlags_DefaultOpen)) {
			imgui.SliderFloat("scale", &imageScale, float32(1), float32(20))
			if imageTxtId != 0 {
				imgsize := imgui.NewVec2(float32(imagData.Bounds().Dx())*imageScale, float32(imagData.Bounds().Dy())*imageScale)
				defer imgsize.Delete()
				imgui.Image(uintptr(imageTxtId), imgsize)
			}
		}

		if imgui.CollapsingHeader("fiile list", int(imgui.TreeNodeFlags_DefaultOpen)) {
			AssetInfo(act, "")
			if imgui.TreeNode("a") {
				AssetInfo(act, "a")
				imgui.TreePop()
				imgui.Separator()
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

func AssetInfo(act *app.Activity, dir string) {
	assetMgr := act.AssetManager()
	assetDir := assetMgr.OpenDir(dir)
	defer assetDir.Close()
	for {
		fname := assetDir.GetNextFileName()
		if fname == "" {
			return
		}
		if imgui.TreeNode(fname) {
			fullpath := ""
			if dir == "" {
				fullpath = fname
			} else {
				fullpath = dir + "/" + fname
			}
			asset := assetMgr.Open(fullpath, storage.ASSET_MODE_RANDOM)
			if asset != nil {
				imgui.Text(fmt.Sprintf("size:%d, IsAllocated:%v", asset.Length(), asset.IsAllocated()))

				context := ""
				isTxt := strings.ToLower(filepath.Ext(fullpath)) == ".txt"

				var buf = make([]byte, 16)
				i, _ := asset.Read(buf)
				if isTxt {
					context += string(buf) + "\n"
				} else {
					context = fmt.Sprintln(buf[:i])
				}
				asset.Seek(-16, io.SeekEnd)
				i, _ = asset.Read(buf)
				context += fmt.Sprintln("\t...")
				if isTxt {
					context += string(buf) + "\n"
				} else {
					context += fmt.Sprintln(buf[:i])
				}
				asset.Close()

				imgui.Text(context)
			} else {
				imgui.Text("open fail.")
			}
			imgui.TreePop()
			imgui.Separator()
		}
	}
}

func ReadImage(act *app.Activity, fname string) image.Image {
	// load image
	m := act.AssetManager()
	fa := m.Open(fname, storage.ASSET_MODE_BUFFER)
	buf := make([]byte, fa.Length())
	fa.Read(buf)

	img, _, err := image.Decode(bytes.NewBuffer(buf))
	if err != nil {
		log.Println("ReadImage:", err)
	}
	fa.Close()

	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)
	log.Printf("   %T %T", img, dst)

	return dst
}

func genTexture(p unsafe.Pointer, w, h int) uint32 {
	var textureId uint32
	// Upload texture to graphics system
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h),
		0, gl.RGBA, gl.UNSIGNED_BYTE, p)

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	return textureId
}
