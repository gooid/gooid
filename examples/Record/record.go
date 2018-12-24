package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	audio "github.com/gooid/audio"
	"github.com/gooid/audio/al"
	"github.com/gooid/audio/alc"
	gl "github.com/gooid/gl/es2"
	"github.com/gooid/imgui"
	"github.com/gooid/imgui/util"
)

var mouseLeft = false
var mouseRight = false
var lastX, lastY int

//

const WINDOWSCALE = 1
const title = "Record "

var (
	width   = (1080 / 8 * 3)
	height  = (1920 / 8 * 3)
	density = (480 / 8 * 3)

	im *util.Render
)

func create() bool {
	// Create window with graphics context
	log.Println("Create...", width, height)

	// gl init
	gl.Init()

	// Setup Dear ImGui binding
	imgui.CreateContext()

	// Setup style
	//imgui.StyleColorsClassic()
	imgui.StyleColorsDark()
	//imgui.StyleColorsLight()

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

		density = getDensity()
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

	io.SetConfigFlags(io.GetConfigFlags() |
		//int(imgui.ConfigFlags_IsTouchScreen) |
		int(imgui.ConfigFlags_NoMouseCursorChange))

	im = util.NewRender("#version 100")

	if _, err := os.Stat(RECORDPATH); os.IsNotExist(err) {
		os.MkdirAll(RECORDPATH, 0666)
	}

	return true
}

func initEGL() {
	log.Println("WINSIZE:", width, height)
	width = int(float32(width) / WINDOWSCALE)
	height = int(float32(height) / WINDOWSCALE)

	if gl.GoStr(gl.GetString(gl.RENDERER)) == "" {
		panic(gl.GetError())
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
	MARGIN := width / 20

	imgui.NewFrame()
	curpos := imgui.NewVec2(MARGIN, MARGIN)
	defer curpos.Delete()
	imgui.SetNextWindowPos(curpos)
	isize := imgui.NewVec2(width-2*MARGIN, height-2*MARGIN)
	defer isize.Delete()
	imgui.SetNextWindowSize(isize)

	imgui.Begin(title, (*bool)(nil), imgui.WindowFlags_NoSavedSettings|
		//imgui.WindowFlags_NoTitleBar|
		imgui.WindowFlags_NoResize|
		imgui.WindowFlags_NoCollapse|
		imgui.WindowFlags_NoMove)

	if stopFunc == nil {
		if imgui.Button("Record") {
			fname := fmt.Sprint("record_", time.Now().UnixNano(), ".wav")
			fw, err := os.Create(filepath.Join(RECORDPATH, fname))
			if err == nil {
				fn, err := Record(fw)
				if err == nil {
					stopFunc = func() {
						fn()
						fw.Close()
					}
				} else {
					log.Println(" CaptureOpen:", err)
					fw.Close()
					os.Remove(filepath.Join(RECORDPATH, fname))
				}
			} else {
				log.Println(" \t", err)
			}
		}
	} else {
		if imgui.Button("Stop") {
			stopFunc()
			stopFunc = nil
		}
	}

	if stopFunc == nil {
		if imgui.CollapsingHeader("record list", imgui.TreeNodeFlags_DefaultOpen) {
			RecordsList(RECORDPATH)
		}
	}

	imgui.End()

	imgui.Render()

	// Rendering
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	im.Render(imgui.GetDrawData())

	SwapBuffers()
}

// 同时表示在 record
var stopFunc func()

func Play(fname string, stopped func()) (*audio.AudioInfo, func(), error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, err
	}

	af, err := audio.NewAudio(f)
	if err != nil {
		return nil, nil, err
	}

	info := af.Info()
	play, err := audio.NewPlayer(af,
		audio.Format(info.Format),
		int64(info.SampleRate))
	if err != nil {
		log.Println(" NewPlayer:", err)
		return nil, nil, err
	}
	err = play.Play()
	if err != nil {
		log.Println(" Play:", err)
		return nil, nil, err
	}

	go func() {
		time.Sleep(time.Millisecond * 300)
		for play.State() != audio.Stopped {
			time.Sleep(time.Millisecond * 300)
		}
		log.Println(" Stop:...")
		play.Destroy()
		play = nil
		af.Close()
		f.Close()
		delete(recordMap, fname)
		stopped()
	}()

	return &info, func() {
		if play != nil {
			play.Stop()
		}
	}, nil
}

func Record(w io.WriteSeeker) (func(), error) {
	stop := make(chan bool)
	wait := sync.WaitGroup{}

	// SAMPLERATE 和 SAMPLESIZE 分别是 44100、2048 这两个值等比缩小，
	// 调用 CaptureOpen 就能成功
	const DENOMINATOR = 1
	const SAMPLERATE = 44100 / DENOMINATOR
	const SAMPLESIZE = 2048 / DENOMINATOR

	const FORMAT = al.FormatStereo16
	const TIMES = 20

	bytes := 1
	switch FORMAT {
	case al.FormatMono8:
		bytes = 1
	case al.FormatMono16:
		bytes = 2
	case al.FormatStereo8:
		bytes = 2
	case al.FormatStereo16:
		bytes = 4
	}

	// 2倍理论大小的BUFFER，以避免频繁分配
	bs := make([]byte, 2*SAMPLERATE*bytes/TIMES)

	audio.WavWriteHeader(w, SAMPLERATE, FORMAT)

	d := alc.CaptureOpen("", SAMPLERATE, FORMAT, SAMPLESIZE)
	if d != nil {
		wait.Add(1)
		go func() {
			defer wait.Done()

			d.Start()
			log.Println("ALC: Start", alc.Error(d.Error()))

			sample := int32(0)
			done := false
			for !done {
				select {
				case <-stop:
					done = true
				case <-time.After(time.Second / TIMES):
				}

				// 得到采样数
				d.GetIntegerv(alc.CaptureSamples, 4, &sample)
				log.Println("   GetIntegerv:", sample)

				// 数据大小 = 采样数 * 每个采样数大小
				size := sample * int32(bytes)
				if int(size) > len(bs) {
					bs = make([]byte, size)
				}

				d.Samples(bs, int64(sample))
				w.Write(bs[:size])
			}
			log.Println("   CaptureStop")
			d.Stop()
			d.Close()
			audio.WavClose(w)
			close(stop)
		}()

		// 返回停止函数
		return func() {
			stop <- true
			wait.Wait()
		}, nil
	} else {
		err := errors.New(alc.Error(al.Error()))
		return nil, err
	}
}

type recordStat struct {
	stop func()
	info string
}

var recordMap = map[string]recordStat{}

func RecordsList(dir string) {
	fs, _ := filepath.Glob(filepath.Join(dir, "*"))
	for _, fname := range fs {
		var i64 int64
		name := filepath.Base(fname)
		fmt.Sscanf(name, "record_%d.wav", &i64)
		t := time.Duration(i64)
		ti := time.Unix(int64(t/time.Nanosecond/time.Second), int64((t/time.Nanosecond)%time.Second))

		if state, ok := recordMap[fname]; ok {
			if imgui.Button(ti.Format("STOP 2006-01-02 15:04:05")) {
				state.stop()
				delete(recordMap, fname)
			}
			size := imgui.NewVec2(-1, -1)
			defer size.Delete()
			imgui.TextWrapped(state.info)
		} else {
			if imgui.Button(ti.Format("PLAY 2006-01-02 15:04:05")) {
				info, stop, err := Play(fname, func() { Wake() })
				if err != nil {
					log.Println("ERROR:" + err.Error())
				} else {
					log.Println("Play  :", fname)
					text := fmt.Sprintf("%+v", info)
					text = strings.Replace(text, " ", "\n", -1)
					recordMap[fname] = recordStat{stop: stop, info: text}
				}
			}
		}
	}
}
