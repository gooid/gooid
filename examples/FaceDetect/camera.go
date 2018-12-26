package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gooid/gooid/camera"
	"github.com/gooid/gooid/examples/CameraDemo/render"
	"github.com/gooid/imgui"
)

const (
	NONE = iota
	READY
	RESIZING
)

type cameraUI struct {
	previewSizes [][2]int
	comboText    string

	previewIndex int
	rotation     int
	irender      render.Render
	imgFormat    string
}

type cameraObj struct {
	sync.Mutex
	camera.Camera
	cameraUI

	callback  func([]byte) bool
	w, h      int
	imageData DataUnit
	isUpdate  bool
	camStat   int32 // atomic.???Int32
}

// 为了避免数据拷贝，防止数据在使用过程中被改写
type DataUnit struct {
	lock int32
	data []byte
}

// 成功则返回 true，否则已被 lock（在使用中）
func (d *DataUnit) Lock() bool {
	return atomic.CompareAndSwapInt32(&d.lock, 0, 1)
}

func (d *DataUnit) Unlock() {
	b := atomic.CompareAndSwapInt32(&d.lock, 1, 0)
	if !b {
		panic(errors.New("DataUnit unlock fail"))
	}
}

func (d *DataUnit) Data() []byte {
	return d.data
}

func (d *DataUnit) SetData(in []byte) bool {
	if d.Lock() {
		if len(d.data) < len(in) {
			d.data = make([]byte, len(in))
		}
		copy(d.data, in)
		d.Unlock()
		return true
	}
	return false
}

func (cam *cameraObj) setYuv(buffer []byte) {
	if cam.imageData.SetData(buffer) {
		faceDetect.PrepareYuv(cam.w, cam.h, cam.rotation, cam.imageData)
		cam.isUpdate = true
	} else {
		//log.Println("ignore: buffer", len(buffer))
	}
}

func (cam *cameraObj) Image() (*DataUnit, bool) {
	cam.Lock()
	defer cam.Unlock()
	if len(cam.imageData.Data()) > 0 {
		return &cam.imageData, cam.isUpdate
	}
	return nil, false
}

func (cam *cameraObj) Release() {
	if cam != nil && cam.Camera != 0 {
		cam.Disconnect()
	}
	if cam.irender != nil {
		cam.irender.Release()
		cam.irender = nil
	}
}

func cameraInit(id int, usercb func(w, h int, img []byte) bool) *cameraObj {
	cam := &cameraObj{}

	cb := func(buffer []byte) bool {
		// init done
		if atomic.CompareAndSwapInt32(&cam.camStat, NONE, READY) {
			return true
		}

		cam.Lock()

		wh := cam.previewSizes[cam.previewIndex]
		if cam.camStat == RESIZING && cam.irender.Validate(wh[0], wh[1], buffer) {
			cam.ResetProperty()
			cam.w, cam.h = wh[0], wh[1]
			atomic.CompareAndSwapInt32(&cam.camStat, RESIZING, READY)
		}

		if !usercb(cam.w, cam.h, buffer) {
			cam.setYuv(buffer)
		}
		cam.Unlock()
		runtime.Gosched()
		return true
	}

	cam.Camera = camera.Connect(id, cb)
	if cam.Camera == 0 {
		log.Println("Cammera connect fail")
		return nil
	}
	cam.callback = cb

	cam.w, cam.h = cam.FrameSize()
	log.Println("camera:", cam, cam.w, cam.h, cam.Fps())
	log.Println("\t", cam.PreviewFormat())
	log.Println("\t", cam.SupportedPreviewSizes())

	cam.cameraUI.Init(cam.Camera)

	// 第一个为默认分辨率
	// 因 camera 还未初始化完，需就异步执行
	go func() {
		runtime.LockOSThread()
		for cam.camStat == NONE {
			time.Sleep(time.Millisecond * 100)
		}
		cam.setPreviewSize(0)
	}()
	return cam
}

func (cam *cameraObj) setPreviewSize(i int) {
	w, h := cam.getPreviewSize(i)
	if cam != nil {
		cam.Lock()
		if atomic.CompareAndSwapInt32(&cam.camStat, READY, RESIZING) {
			cam.SetFrameSize(w, h)
			cam.ApplyProperties()
			cam.previewIndex = i
		}
		cam.Unlock()
	}
}

func (cam *cameraUI) Init(nativeObj camera.Camera) {
	str := nativeObj.SupportedPreviewSizes()
	for _, s := range strings.Split(str, ",") {
		w, h := 0, 0
		fmt.Sscanf(string(s), "%dx%d", &w, &h)
		cam.previewSizes = append(cam.previewSizes, [2]int{w, h})
	}

	cam.rotation = render.ROTATION90
	cam.comboText = strings.Replace(str, ",", "\x00", -1) + "\x00"
	cam.imgFormat = nativeObj.PreviewFormat()
}

func (cam *cameraUI) ResetProperty() {
	if cam != nil && cam.irender != nil {
		iw, ih := cam.getPreviewSize(cam.previewIndex)
		cam.irender.SetProperty(width, height, iw, ih, 0, 0, width, height, cam.rotation)
	}
}
func (cam *cameraUI) ResetRender() {
	if cam.irender != nil {
		cam.irender.Release()
	}

	switch cam.imgFormat {
	case "yuv420sp":
		cam.irender = &render.YuvRender{}
	default:
		log.Println("not support format:", cam.imgFormat)
		return
	}

	cam.irender.Init()
	cam.ResetProperty()
}

func (cam *cameraUI) getPreviewSize(i int) (w, h int) {
	if i < len(cam.previewSizes) {
		return cam.previewSizes[i][0], cam.previewSizes[i][1]
	}
	return
}

func (cam *cameraObj) Draw() {
	item := int(cam.previewIndex)
	imgui.Text("pixels")
	imgui.SameLine()
	imgui.Combo("", &item, cam.comboText)
	if item != cam.previewIndex {
		cam.setPreviewSize(item)
	}
	imgui.Separator()
	r := int(cam.rotation)
	if imgui.Combo("degree", &r, "0\x0090\x00180\x00270\x00\x00") {
		cam.rotation = r
		cam.ResetProperty()
	}
}

func (cam *cameraObj) DrawImage() {
	if cam != nil {
		data, _ := cam.Image()
		if data != nil && data.Lock() {
			defer data.Unlock()
			cam.irender.Draw(data.Data())
		}
	}
}
